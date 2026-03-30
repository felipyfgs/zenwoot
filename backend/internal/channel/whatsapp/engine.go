package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"wzap/internal/broker"
	"wzap/internal/config"
	"wzap/internal/dispatcher"
	"wzap/internal/domain"
	"wzap/internal/model"
	"wzap/internal/processor"
	"wzap/internal/repo"
	"wzap/internal/ws"
)

// InboxEngine connects the WzapClient (HTTP) to the zenwoot domain model.
// It replaces the older embedded engine entirely.
type InboxEngine struct {
	client        *WzapClient
	waChannelRepo *repo.WhatsAppChannelRepository
	inboxRepo     *repo.InboxRepository
	nats          *broker.Nats
	disp          *dispatcher.Dispatcher
	wsHub         *ws.Hub
	msgProcessor  processor.MessageProcessor
	accountID     string
	ctx           context.Context
}

func NewInboxEngine(
	ctx context.Context,
	cfg *config.Config,
	waChannelRepo *repo.WhatsAppChannelRepository,
	inboxRepo *repo.InboxRepository,
	nats *broker.Nats,
	disp *dispatcher.Dispatcher,
	wsHub *ws.Hub,
	msgProcessor processor.MessageProcessor,
	accountID string,
) (*InboxEngine, error) {
	client := NewWzapClient(cfg.WzapBaseURL, cfg.WzapAdminKey)

	e := &InboxEngine{
		client:        client,
		waChannelRepo: waChannelRepo,
		inboxRepo:     inboxRepo,
		nats:          nats,
		disp:          disp,
		wsHub:         wsHub,
		msgProcessor:  msgProcessor,
		accountID:     accountID,
		ctx:           ctx,
	}

	e.wireCallbacks()
	return e, nil
}

// GetClient returns the WzapClient (HTTP client) — used by services.
func (e *InboxEngine) GetClient() *WzapClient {
	return e.client
}

// GetSessionAPIKey returns the API key for a wzap session (channelID).
// For now it returns the admin key; sessions with dedicated keys can extend this.
func (e *InboxEngine) GetSessionAPIKey(_ string) string {
	return e.client.adminKey
}

func (e *InboxEngine) wireCallbacks() {
	ctx := e.ctx

	e.client.OnQRCode(func(channelID string, qr string) {
		if err := e.waChannelRepo.UpdateQRCode(ctx, channelID, qr); err != nil {
			log.Error().Err(err).Str("channelId", channelID).Msg("Failed to update QR code")
		}
		inbox, err := e.inboxRepo.FindByChannelID(ctx, domain.ChannelWhatsApp, channelID)
		if err != nil {
			return
		}
		if qr == "" {
			_ = e.inboxRepo.UpdateStatus(ctx, inbox.ID, domain.InboxStatusDisconnected)
		} else {
			_ = e.inboxRepo.UpdateStatus(ctx, inbox.ID, domain.InboxStatusConnecting)
			e.publishEvent(inbox.ID, model.EventPairSuccess, map[string]interface{}{"qr": qr})
		}
	})

	e.client.OnConnected(func(channelID string, jid string) {
		if err := e.waChannelRepo.UpdateJID(ctx, channelID, jid); err != nil {
			log.Error().Err(err).Str("channelId", channelID).Msg("Failed to update JID")
		}
		_ = e.waChannelRepo.UpdateQRCode(ctx, channelID, "")
		inbox, err := e.inboxRepo.FindByChannelID(ctx, domain.ChannelWhatsApp, channelID)
		if err != nil {
			return
		}
		_ = e.inboxRepo.UpdateStatus(ctx, inbox.ID, domain.InboxStatusActive)
		e.publishEvent(inbox.ID, model.EventConnected, map[string]interface{}{"jid": jid})
		if e.wsHub != nil {
			e.wsHub.Publish(e.accountID, inbox.ID, "inbox.connected", map[string]interface{}{"jid": jid})
		}
		log.Info().Str("channelId", channelID).Str("jid", jid).Msg("Inbox connected")
	})

	e.client.OnDisconnected(func(channelID string) {
		inbox, err := e.inboxRepo.FindByChannelID(ctx, domain.ChannelWhatsApp, channelID)
		if err != nil {
			return
		}
		_ = e.inboxRepo.UpdateStatus(ctx, inbox.ID, domain.InboxStatusDisconnected)
		e.publishEvent(inbox.ID, model.EventDisconnected, map[string]interface{}{})
		if e.wsHub != nil {
			e.wsHub.Publish(e.accountID, inbox.ID, "inbox.disconnected", map[string]interface{}{})
		}
	})

	e.client.OnLoggedOut(func(channelID string) {
		inbox, err := e.inboxRepo.FindByChannelID(ctx, domain.ChannelWhatsApp, channelID)
		if err != nil {
			return
		}
		_ = e.inboxRepo.UpdateStatus(ctx, inbox.ID, domain.InboxStatusInactive)
		e.publishEvent(inbox.ID, model.EventLoggedOut, map[string]interface{}{})
	})

	e.client.OnMessage(func(channelID string, msg *domain.IncomingMessage) {
		inbox, err := e.inboxRepo.FindByChannelID(ctx, domain.ChannelWhatsApp, channelID)
		if err != nil {
			log.Error().Err(err).Str("channelId", channelID).Msg("Inbox not found for incoming message")
			return
		}

		if e.msgProcessor != nil {
			if err := e.msgProcessor.Process(ctx, inbox.ID, msg); err != nil {
				log.Error().Err(err).Str("inboxId", inbox.ID).Msg("Failed to process incoming message")
			}
		}

		e.publishEvent(inbox.ID, model.EventMessage, map[string]interface{}{
			"id":          msg.ExternalID,
			"from":        msg.From,
			"pushName":    msg.PushName,
			"content":     msg.Content,
			"contentType": string(msg.ContentType),
			"fromMe":      msg.IsFromMe,
		})
	})

	e.client.OnStatus(func(channelID string, status *domain.StatusUpdate) {
		inbox, err := e.inboxRepo.FindByChannelID(ctx, domain.ChannelWhatsApp, channelID)
		if err != nil {
			return
		}

		if e.msgProcessor != nil {
			if err := e.msgProcessor.ProcessStatus(ctx, inbox.ID, status); err != nil {
				log.Error().Err(err).Str("inboxId", inbox.ID).Msg("Failed to process message status")
			}
		}

		e.publishEvent(inbox.ID, model.EventReceipt, map[string]interface{}{
			"messageId": status.ExternalID,
			"status":    string(status.Status),
		})
	})
}

// Connect starts QR pairing or reconnects a WhatsApp channel via wzap HTTP API.
func (e *InboxEngine) Connect(ctx context.Context, channelID string) (string, error) {
	// Ensure session exists in wzap; create if not.
	_, err := e.client.GetSession(ctx, channelID)
	if err != nil {
		// Session not found — create it using channelID as the session name.
		if _, createErr := e.client.CreateSession(ctx, channelID); createErr != nil {
			return "", fmt.Errorf("failed to create wzap session: %w", createErr)
		}
	}

	if _, err := e.client.Connect(ctx, channelID, ""); err != nil {
		return "", fmt.Errorf("failed to connect wzap session: %w", err)
	}

	if e.client.IsConnected(channelID) {
		return "CONNECTED", nil
	}
	return "PAIRING", nil
}

func (e *InboxEngine) Disconnect(channelID string) {
	e.client.Disconnect(channelID)
}

func (e *InboxEngine) IsConnected(channelID string) bool {
	return e.client.IsConnected(channelID)
}

func (e *InboxEngine) GetQRCode(ctx context.Context, channelID string) (string, error) {
	return e.waChannelRepo.FindQRByID(ctx, channelID)
}

// ReconnectAll is a no-op for the wzap HTTP provider — wzap manages its own reconnections.
func (e *InboxEngine) ReconnectAll(_ context.Context) error {
	log.Info().Msg("ReconnectAll: wzap manages session persistence — skipping")
	return nil
}

func (e *InboxEngine) publishEvent(inboxID string, eventType model.EventType, payload map[string]interface{}) {
	payload["eventId"] = uuid.NewString()
	payload["inboxId"] = inboxID
	payload["event"] = eventType
	payload["timestamp"] = time.Now().Format(time.RFC3339)

	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Str("inboxId", inboxID).Msg("Failed to marshal event payload")
		return
	}
	if e.nats != nil {
		if err := e.nats.Publish(e.ctx, "wzap.events."+inboxID, data); err != nil {
			log.Error().Err(err).Str("inboxId", inboxID).Msg("Failed to publish NATS event")
		}
	}
	if e.disp != nil {
		go e.disp.Dispatch(inboxID, eventType, data)
	}
}
