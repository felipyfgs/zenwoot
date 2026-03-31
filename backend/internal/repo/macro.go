package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type MacroRepo struct {
	BaseRepo[models.Macro]
}

func NewMacroRepo(db *bun.DB) *MacroRepo {
	return &MacroRepo{BaseRepo: *NewBaseRepo[models.Macro](db)}
}

func (r *MacroRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.Macro, error) {
	var items []*models.Macro
	err := r.WithTenant(ctx, accountID).OrderExpr(`"name" ASC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("macroRepo.ListByAccount: %w", err)
	}
	return items, nil
}

func (r *MacroRepo) Create(ctx context.Context, m *models.Macro) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("macroRepo.Create: %w", err)
	}
	return nil
}

func (r *MacroRepo) Update(ctx context.Context, m *models.Macro) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("macroRepo.Update: %w", err)
	}
	return nil
}

func (r *MacroRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"macros"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("macroRepo.Delete: %w", err)
	}
	return nil
}

func (r *MacroRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Macro, error) {
	var m models.Macro
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("macroRepo.GetByID: %w", err)
	}
	return &m, nil
}
