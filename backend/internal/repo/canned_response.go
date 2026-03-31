package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type CannedResponseRepo struct {
	BaseRepo[models.CannedResponse]
}

func NewCannedResponseRepo(db *bun.DB) *CannedResponseRepo {
	return &CannedResponseRepo{BaseRepo: *NewBaseRepo[models.CannedResponse](db)}
}

func (r *CannedResponseRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.CannedResponse, error) {
	var items []*models.CannedResponse
	err := r.WithTenant(ctx, accountID).OrderExpr(`"short_code" ASC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("cannedResponseRepo.ListByAccount: %w", err)
	}
	return items, nil
}

func (r *CannedResponseRepo) Create(ctx context.Context, m *models.CannedResponse) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("cannedResponseRepo.Create: %w", err)
	}
	return nil
}

func (r *CannedResponseRepo) Update(ctx context.Context, m *models.CannedResponse) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("cannedResponseRepo.Update: %w", err)
	}
	return nil
}

func (r *CannedResponseRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"canned_responses"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("cannedResponseRepo.Delete: %w", err)
	}
	return nil
}

func (r *CannedResponseRepo) GetByID(ctx context.Context, accountID, id int64) (*models.CannedResponse, error) {
	var m models.CannedResponse
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("cannedResponseRepo.GetByID: %w", err)
	}
	return &m, nil
}

func (r *CannedResponseRepo) GetByShortCode(ctx context.Context, accountID int64, shortCode string) (*models.CannedResponse, error) {
	var m models.CannedResponse
	err := r.WithTenant(ctx, accountID).Where(`"short_code" = ?`, shortCode).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("cannedResponseRepo.GetByShortCode: %w", err)
	}
	return &m, nil
}
