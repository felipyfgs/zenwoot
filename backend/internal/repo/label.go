package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type LabelRepo struct {
	BaseRepo[models.Label]
}

func NewLabelRepo(db *bun.DB) *LabelRepo {
	return &LabelRepo{BaseRepo: *NewBaseRepo[models.Label](db)}
}

func (r *LabelRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.Label, error) {
	var items []*models.Label
	err := r.WithTenant(ctx, accountID).OrderExpr(`"title" ASC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("labelRepo.ListByAccount: %w", err)
	}
	return items, nil
}

func (r *LabelRepo) Create(ctx context.Context, m *models.Label) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("labelRepo.Create: %w", err)
	}
	return nil
}

func (r *LabelRepo) Update(ctx context.Context, m *models.Label) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("labelRepo.Update: %w", err)
	}
	return nil
}

func (r *LabelRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"labels"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("labelRepo.Delete: %w", err)
	}
	return nil
}

func (r *LabelRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Label, error) {
	var item models.Label
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &item)
	if err != nil {
		return nil, fmt.Errorf("labelRepo.GetByID: %w", err)
	}
	return &item, nil
}
