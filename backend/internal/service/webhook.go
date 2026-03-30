package service

import (
	"context"
	"fmt"
	"time"

	"wzap/internal/dto"
	"wzap/internal/model"
	"wzap/internal/repo"

	"github.com/google/uuid"
)

type WebhookService struct {
	repo *repo.WebhookRepository
}

func NewWebhookService(repo *repo.WebhookRepository) *WebhookService {
	return &WebhookService{repo: repo}
}

func (s *WebhookService) Create(ctx context.Context, inboxID string, req dto.CreateWebhookReq) (*model.Webhook, error) {
	if len(req.Events) == 0 {
		return nil, fmt.Errorf("at least one event type is required")
	}
	events := make([]string, 0, len(req.Events))
	for _, e := range req.Events {
		if !model.ValidEventTypes[e] {
			return nil, fmt.Errorf("invalid event type: %s", e)
		}
		events = append(events, string(e))
	}
	webhook := &model.Webhook{
		ID:          uuid.NewString(),
		InboxID:     inboxID,
		URL:         req.URL,
		Secret:      req.Secret,
		Events:      events,
		Enabled:     true,
		NatsEnabled: req.NatsEnabled,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, webhook); err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return webhook, nil
}

func (s *WebhookService) List(ctx context.Context, inboxID string) ([]model.Webhook, error) {
	return s.repo.FindByInboxID(ctx, inboxID)
}

func (s *WebhookService) Delete(ctx context.Context, inboxID, webhookID string) error {
	return s.repo.Delete(ctx, inboxID, webhookID)
}
