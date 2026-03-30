package repo

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type CustomFilterRepo struct {
	db *bun.DB
}

func NewCustomFilterRepo(db *bun.DB) *CustomFilterRepo {
	return &CustomFilterRepo{db: db}
}

func (r *CustomFilterRepo) Create(ctx context.Context, filter *models.CustomFilter) error {
	_, err := r.db.NewInsert().Model(filter).Exec(ctx)
	return err
}

func (r *CustomFilterRepo) Update(ctx context.Context, filter *models.CustomFilter) error {
	_, err := r.db.NewUpdate().Model(filter).Where("id = ?", filter.ID).Exec(ctx)
	return err
}

func (r *CustomFilterRepo) Delete(ctx context.Context, id, accountID, userID int64) error {
	_, err := r.db.NewDelete().
		Model(&models.CustomFilter{}).
		Where("id = ? AND account_id = ? AND user_id = ?", id, accountID, userID).
		Exec(ctx)
	return err
}

func (r *CustomFilterRepo) List(ctx context.Context, accountID, userID int64) ([]models.CustomFilter, error) {
	var filters []models.CustomFilter
	err := r.db.NewSelect().
		Model(&filters).
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Order("created_at DESC").
		Scan(ctx)
	return filters, err
}

func (r *CustomFilterRepo) GetByID(ctx context.Context, id, accountID, userID int64) (*models.CustomFilter, error) {
	var filter models.CustomFilter
	err := r.db.NewSelect().
		Model(&filter).
		Where("id = ? AND account_id = ? AND user_id = ?", id, accountID, userID).
		Scan(ctx)
	return &filter, err
}
