package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type ConversationRepo struct {
	BaseRepo[models.Conversation]
}

func NewConversationRepo(db *bun.DB) *ConversationRepo {
	return &ConversationRepo{BaseRepo: *NewBaseRepo[models.Conversation](db)}
}

type ConversationFilter struct {
	Status     *int
	AssigneeID *int64
	InboxID    *int64
	TeamID     *int64
	Page       int
	PageSize   int
}

func (r *ConversationRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Conversation, error) {
	var m models.Conversation
	err := r.WithTenant(ctx, accountID).
		Where(`"id" = ?`, id).
		Relation("Assignee").
		Relation("Contact").
		Relation("Inbox").
		Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("conversationRepo.GetByID: %w", err)
	}
	return &m, nil
}

func (r *ConversationRepo) List(ctx context.Context, accountID int64, f ConversationFilter) ([]*models.Conversation, int, error) {
	var items []*models.Conversation
	q := r.WithTenant(ctx, accountID)

	if f.Status != nil {
		q = q.Where(`"status" = ?`, *f.Status)
	}
	if f.AssigneeID != nil {
		q = q.Where(`"assignee_id" = ?`, *f.AssigneeID)
	}
	if f.InboxID != nil {
		q = q.Where(`"inbox_id" = ?`, *f.InboxID)
	}
	if f.TeamID != nil {
		q = q.Where(`"team_id" = ?`, *f.TeamID)
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 25
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("conversationRepo.List count: %w", err)
	}
	err = q.OrderExpr(`"last_activity_at" DESC`).
		Limit(f.PageSize).
		Offset((f.Page-1)*f.PageSize).
		Scan(ctx, &items)
	if err != nil {
		return nil, 0, fmt.Errorf("conversationRepo.List scan: %w", err)
	}
	return items, total, nil
}

func (r *ConversationRepo) Create(ctx context.Context, m *models.Conversation) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("conversationRepo.Create: %w", err)
	}
	return nil
}

func (r *ConversationRepo) Update(ctx context.Context, m *models.Conversation) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("conversationRepo.Update: %w", err)
	}
	return nil
}
