package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type ContactRepo struct {
	BaseRepo[models.Contact]
}

func NewContactRepo(db *bun.DB) *ContactRepo {
	return &ContactRepo{BaseRepo: *NewBaseRepo[models.Contact](db)}
}

func (r *ContactRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Contact, error) {
	var m models.Contact
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("contactRepo.GetByID: %w", err)
	}
	return &m, nil
}

func (r *ContactRepo) Search(ctx context.Context, accountID int64, q string, page, pageSize int) ([]*models.Contact, int, error) {
	var items []*models.Contact
	query := r.WithTenant(ctx, accountID).
		Where(`(COALESCE("name",'') || ' ' || COALESCE("email",'') || ' ' ||
                COALESCE("phone_number",'') || ' ' || COALESCE("identifier",'')) ILIKE ?`,
			"%"+q+"%")
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("contactRepo.Search count: %w", err)
	}
	err = query.OrderExpr(`"name" ASC`).Limit(pageSize).Offset((page-1)*pageSize).Scan(ctx, &items)
	if err != nil {
		return nil, 0, fmt.Errorf("contactRepo.Search scan: %w", err)
	}
	return items, total, nil
}

func (r *ContactRepo) Create(ctx context.Context, m *models.Contact) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("contactRepo.Create: %w", err)
	}
	return nil
}

func (r *ContactRepo) Update(ctx context.Context, m *models.Contact) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("contactRepo.Update: %w", err)
	}
	return nil
}
