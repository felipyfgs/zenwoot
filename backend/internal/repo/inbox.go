package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type InboxRepo struct {
	BaseRepo[models.Inbox]
}

func NewInboxRepo(db *bun.DB) *InboxRepo {
	return &InboxRepo{BaseRepo: *NewBaseRepo[models.Inbox](db)}
}

func (r *InboxRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Inbox, error) {
	var m models.Inbox
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("inboxRepo.GetByID: %w", err)
	}
	return &m, nil
}

func (r *InboxRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.Inbox, error) {
	var items []*models.Inbox
	err := r.WithTenant(ctx, accountID).OrderExpr(`"name" ASC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("inboxRepo.ListByAccount: %w", err)
	}
	return items, nil
}

func (r *InboxRepo) Create(ctx context.Context, m *models.Inbox) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("inboxRepo.Create: %w", err)
	}
	return nil
}

func (r *InboxRepo) Update(ctx context.Context, m *models.Inbox) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("inboxRepo.Update: %w", err)
	}
	return nil
}

func (r *InboxRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"inboxes"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("inboxRepo.Delete: %w", err)
	}
	return nil
}
