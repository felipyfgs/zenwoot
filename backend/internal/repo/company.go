package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type CompanyRepo struct {
	BaseRepo[models.Company]
}

func NewCompanyRepo(db *bun.DB) *CompanyRepo {
	return &CompanyRepo{BaseRepo: *NewBaseRepo[models.Company](db)}
}

func (r *CompanyRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.Company, error) {
	var items []*models.Company
	err := r.WithTenant(ctx, accountID).OrderExpr(`"name" ASC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("companyRepo.ListByAccount: %w", err)
	}
	return items, nil
}

func (r *CompanyRepo) Search(ctx context.Context, accountID int64, q string) ([]*models.Company, error) {
	var items []*models.Company
	err := r.WithTenant(ctx, accountID).
		Where(`"name" ILIKE ?`, "%"+q+"%").
		Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("companyRepo.Search: %w", err)
	}
	return items, nil
}

func (r *CompanyRepo) Create(ctx context.Context, m *models.Company) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("companyRepo.Create: %w", err)
	}
	return nil
}

func (r *CompanyRepo) Update(ctx context.Context, m *models.Company) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("companyRepo.Update: %w", err)
	}
	return nil
}

func (r *CompanyRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"companies"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("companyRepo.Delete: %w", err)
	}
	return nil
}

func (r *CompanyRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Company, error) {
	var m models.Company
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("companyRepo.GetByID: %w", err)
	}
	return &m, nil
}
