package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"wzap/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WebhookRepository struct {
	db *pgxpool.Pool
}

func NewWebhookRepository(db *pgxpool.Pool) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) Create(ctx context.Context, w *model.Webhook) error {
	query := `INSERT INTO "wzWebhooks" ("id", "inboxId", "url", "secret", "events", "enabled", "natsEnabled", "createdAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, w.ID, w.InboxID, w.URL, w.Secret, w.Events, w.Enabled, w.NatsEnabled, w.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert webhook: %w", err)
	}
	return nil
}

func (r *WebhookRepository) FindByInboxID(ctx context.Context, inboxID string) ([]model.Webhook, error) {
	query := `SELECT "id", "inboxId", "url", COALESCE("secret", ''), "events", "enabled", "natsEnabled", "createdAt", "updatedAt"
		FROM "wzWebhooks" WHERE "inboxId" = $1 ORDER BY "createdAt" DESC`

	rows, err := r.db.Query(ctx, query, inboxID)
	if err != nil {
		return nil, fmt.Errorf("failed to query webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []model.Webhook
	for rows.Next() {
		var w model.Webhook
		if err := rows.Scan(&w.ID, &w.InboxID, &w.URL, &w.Secret, &w.Events, &w.Enabled, &w.NatsEnabled, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, w)
	}
	return webhooks, rows.Err()
}

// FindByInboxID is the new name (alias for compatibility)
func (r *WebhookRepository) FindBySessionID(ctx context.Context, inboxID string) ([]model.Webhook, error) {
	return r.FindByInboxID(ctx, inboxID)
}

func (r *WebhookRepository) FindActiveByInboxAndEvent(ctx context.Context, inboxID string, eventType string) ([]model.Webhook, error) {
	query := `SELECT "id", "inboxId", "url", COALESCE("secret", ''), "events", "enabled", "natsEnabled", "createdAt", "updatedAt"
		FROM "wzWebhooks"
		WHERE "inboxId" = $1
		  AND "enabled" = true
		  AND ("events" @> $2::jsonb OR "events" @> '["All"]'::jsonb)`

	eventJSON, _ := json.Marshal([]string{eventType})
	rows, err := r.db.Query(ctx, query, inboxID, eventJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to query active webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []model.Webhook
	for rows.Next() {
		var w model.Webhook
		if err := rows.Scan(&w.ID, &w.InboxID, &w.URL, &w.Secret, &w.Events, &w.Enabled, &w.NatsEnabled, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, w)
	}
	return webhooks, rows.Err()
}

// FindActiveBySessionAndEvent is an alias for backward compat
func (r *WebhookRepository) FindActiveBySessionAndEvent(ctx context.Context, inboxID string, eventType string) ([]model.Webhook, error) {
	return r.FindActiveByInboxAndEvent(ctx, inboxID, eventType)
}

func (r *WebhookRepository) Delete(ctx context.Context, inboxID, webhookID string) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM "wzWebhooks" WHERE "id" = $1 AND "inboxId" = $2`,
		webhookID, inboxID)
	if err != nil {
		return fmt.Errorf("failed to delete webhook %s: %w", webhookID, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("webhook %s not found", webhookID)
	}
	return nil
}
