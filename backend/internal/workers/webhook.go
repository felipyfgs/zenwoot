package workers

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/logger"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type WebhookWorker struct {
	nc         *nats.Conn
	db         *bun.DB
	webhookSvc *services.WebhookService
	subs       []*nats.Subscription
}

func NewWebhookWorker(nc *nats.Conn, db *bun.DB, webhookSvc *services.WebhookService) *WebhookWorker {
	return &WebhookWorker{nc: nc, db: db, webhookSvc: webhookSvc}
}

func (w *WebhookWorker) Start() error {
	subjects := []string{
		"zenwoot.conversation.created",
		"zenwoot.conversation.updated",
		"zenwoot.conversation.resolved",
		"zenwoot.message.created",
		"zenwoot.contact.created",
	}
	for _, subj := range subjects {
		sub, err := w.nc.Subscribe(subj, w.handle(subj))
		if err != nil {
			return err
		}
		w.subs = append(w.subs, sub)
	}
	logger.Info().Msg("WebhookWorker started")
	return nil
}

func (w *WebhookWorker) Stop() {
	for _, sub := range w.subs {
		_ = sub.Unsubscribe()
	}
}

func (w *WebhookWorker) handle(event string) nats.MsgHandler {
	return func(msg *nats.Msg) {
		if w.webhookSvc == nil {
			logger.Error().Msg("webhookSvc is nil in WebhookWorker")
			return
		}

		var payload map[string]any
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			logger.Warn().Err(err).Msg("failed to unmarshal webhook event")
			return
		}
		accountIDf, ok := payload["account_id"].(float64)
		if !ok {
			logger.Warn().Msg("account_id not found in payload or invalid type")
			return
		}
		accountID := int64(accountIDf)
		if accountID == 0 {
			return
		}
		w.webhookSvc.Dispatch(context.Background(), accountID, event, payload)
	}
}
