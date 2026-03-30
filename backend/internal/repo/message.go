package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type MessageRepo struct {
	BaseRepo[models.Message]
}

func NewMessageRepo(db *bun.DB) *MessageRepo {
	return &MessageRepo{BaseRepo: *NewBaseRepo[models.Message](db)}
}

func (r *MessageRepo) ListByConversation(ctx context.Context, accountID, conversationID int64, before *int64, limit int) ([]*models.Message, error) {
	var items []*models.Message
	q := r.WithTenant(ctx, accountID).
		Where(`"conversation_id" = ?`, conversationID).
		Relation("Attachments")

	if before != nil {
		q = q.Where(`"id" < ?`, *before)
	}
	if limit <= 0 {
		limit = 25
	}
	err := q.OrderExpr(`"created_at" DESC`).Limit(limit).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("messageRepo.ListByConversation: %w", err)
	}
	return items, nil
}

func (r *MessageRepo) Create(ctx context.Context, m *models.Message) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("messageRepo.Create: %w", err)
	}
	return nil
}

func (r *MessageRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Message, error) {
	var m models.Message
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("messageRepo.GetByID: %w", err)
	}
	return &m, nil
}

func (r *MessageRepo) Update(ctx context.Context, m *models.Message) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("messageRepo.Update: %w", err)
	}
	return nil
}

func (r *MessageRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"messages"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("messageRepo.Delete: %w", err)
	}
	return nil
}
