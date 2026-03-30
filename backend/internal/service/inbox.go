package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	waChannel "wzap/internal/channel/whatsapp"
	"wzap/internal/domain"
	"wzap/internal/dto"
	"wzap/internal/model"
	"wzap/internal/repo"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// InboxService manages the lifecycle of inboxes (WhatsApp channels).
type InboxService struct {
	webhookRepo   *repo.WebhookRepository
	waChannelRepo *repo.WhatsAppChannelRepository
	inboxRepo     *repo.InboxRepository
	engine        *waChannel.InboxEngine
	accountID     string
}

func NewInboxService(
	webhookRepo *repo.WebhookRepository,
	waChannelRepo *repo.WhatsAppChannelRepository,
	inboxRepo *repo.InboxRepository,
	engine *waChannel.InboxEngine,
	accountID string,
) *InboxService {
	return &InboxService{
		webhookRepo:   webhookRepo,
		waChannelRepo: waChannelRepo,
		inboxRepo:     inboxRepo,
		engine:        engine,
		accountID:     accountID,
	}
}

// Create creates a new inbox with WhatsApp channel.
func (s *InboxService) Create(ctx context.Context, req dto.InboxCreateReq) (*dto.InboxCreatedResp, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if !nameRegex.MatchString(req.Name) {
		return nil, fmt.Errorf("name must contain only letters, numbers, hyphens and underscores")
	}

	now := time.Now()
	channelID := uuid.NewString()
	inboxID := uuid.NewString()

	// 1. Create wzChannelsWhatsapp record
	waCh := &repo.WhatsAppChannelRecord{
		ID:        channelID,
		AccountID: s.accountID,
		Provider:  "wzap",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.waChannelRepo.Create(ctx, waCh); err != nil {
		return nil, fmt.Errorf("failed to create whatsapp channel: %w", err)
	}

	// 2. Create wzInboxes record
	inbox := &domain.Inbox{
		ID:          inboxID,
		AccountID:   s.accountID,
		Name:        req.Name,
		ChannelType: domain.ChannelWhatsApp,
		ChannelID:   channelID,
		Status:      domain.InboxStatusInactive,
		Settings:    req.Settings,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.inboxRepo.Create(ctx, inbox); err != nil {
		return nil, fmt.Errorf("failed to create inbox: %w", err)
	}

	resp := &dto.InboxCreatedResp{
		InboxID: inbox.ID,
		Name:    inbox.Name,
		Status:  string(inbox.Status),
	}

	// 3. Create inline webhook if provided
	if req.Webhook != nil && req.Webhook.URL != "" {
		events := make([]string, 0, len(req.Webhook.Events))
		for _, e := range req.Webhook.Events {
			if !model.ValidEventTypes[e] {
				log.Warn().Str("event", string(e)).Str("inbox", inboxID).Msg("Skipping invalid event type")
				continue
			}
			events = append(events, string(e))
		}
		wh := &model.Webhook{
			ID:        uuid.NewString(),
			InboxID:   inboxID,
			URL:       req.Webhook.URL,
			Events:    events,
			Enabled:   true,
			CreatedAt: now,
		}
		if err := s.webhookRepo.Create(ctx, wh); err != nil {
			log.Warn().Err(err).Str("inbox", inboxID).Msg("Failed to create inline webhook")
		} else {
			resp.Webhook = wh
		}
	}

	return resp, nil
}

// List returns all inboxes for the account.
func (s *InboxService) List(ctx context.Context) ([]domain.Inbox, error) {
	return s.inboxRepo.FindByAccountID(ctx, s.accountID)
}

// Get returns a single inbox by ID.
func (s *InboxService) Get(ctx context.Context, id string) (*domain.Inbox, error) {
	return s.inboxRepo.FindByID(ctx, id)
}

// Delete disconnects and removes an inbox.
func (s *InboxService) Delete(ctx context.Context, id string) error {
	inbox, err := s.inboxRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("inbox not found: %w", err)
	}

	s.engine.Disconnect(inbox.ChannelID)
	_ = s.waChannelRepo.Delete(ctx, inbox.ChannelID)
	return s.inboxRepo.Delete(ctx, inbox.ID)
}

// Connect starts WhatsApp connection for an inbox.
func (s *InboxService) Connect(ctx context.Context, id string) (string, error) {
	inbox, err := s.inboxRepo.FindByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("inbox not found: %w", err)
	}
	return s.engine.Connect(ctx, inbox.ChannelID)
}

// Disconnect stops WhatsApp connection for an inbox.
func (s *InboxService) Disconnect(id string) error {
	inbox, err := s.inboxRepo.FindByID(context.Background(), id)
	if err != nil {
		return nil // Already gone
	}
	s.engine.Disconnect(inbox.ChannelID)
	return nil
}

// GetQRCode returns the current QR code for pairing.
func (s *InboxService) GetQRCode(ctx context.Context, id string) (string, error) {
	inbox, err := s.inboxRepo.FindByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("inbox not found: %w", err)
	}
	return s.engine.GetQRCode(ctx, inbox.ChannelID)
}
