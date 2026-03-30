package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type WebhookRepo struct {
	BaseRepo[models.Webhook]
}

func NewWebhookRepo(db *bun.DB) *WebhookRepo {
	return &WebhookRepo{BaseRepo: *NewBaseRepo[models.Webhook](db)}
}

func (r *WebhookRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.Webhook, error) {
	var items []*models.Webhook
	err := r.WithTenant(ctx, accountID).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("webhookRepo.ListByAccount: %w", err)
	}
	return items, nil
}

func (r *WebhookRepo) ListActiveByEvent(ctx context.Context, accountID int64, event string) ([]*models.Webhook, error) {
	var items []*models.Webhook
	err := r.WithTenant(ctx, accountID).
		Where(`? = ANY("subscriptions")`, event).
		Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("webhookRepo.ListActiveByEvent: %w", err)
	}
	return items, nil
}

func (r *WebhookRepo) Create(ctx context.Context, m *models.Webhook) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("webhookRepo.Create: %w", err)
	}
	return nil
}

func (r *WebhookRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"webhooks"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("webhookRepo.Delete: %w", err)
	}
	return nil
}
