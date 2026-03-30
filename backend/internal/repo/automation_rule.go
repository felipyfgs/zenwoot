package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type AutomationRuleRepo struct {
	BaseRepo[models.AutomationRule]
}

func NewAutomationRuleRepo(db *bun.DB) *AutomationRuleRepo {
	return &AutomationRuleRepo{BaseRepo: *NewBaseRepo[models.AutomationRule](db)}
}

func (r *AutomationRuleRepo) ListActive(ctx context.Context, accountID int64, event_name string) ([]*models.AutomationRule, error) {
	var items []*models.AutomationRule
	err := r.WithTenant(ctx, accountID).
		Where(`"active" = true AND "event_name" = ?`, event_name).
		Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("automationRuleRepo.ListActive: %w", err)
	}
	return items, nil
}

func (r *AutomationRuleRepo) Create(ctx context.Context, m *models.AutomationRule) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("automationRuleRepo.Create: %w", err)
	}
	return nil
}

func (r *AutomationRuleRepo) Update(ctx context.Context, m *models.AutomationRule) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("automationRuleRepo.Update: %w", err)
	}
	return nil
}

func (r *AutomationRuleRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"automation_rules"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("automationRuleRepo.Delete: %w", err)
	}
	return nil
}

func (r *AutomationRuleRepo) GetByID(ctx context.Context, accountID, id int64) (*models.AutomationRule, error) {
	var m models.AutomationRule
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("automationRuleRepo.GetByID: %w", err)
	}
	return &m, nil
}

func (r *AutomationRuleRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.AutomationRule, error) {
	var items []*models.AutomationRule
	err := r.WithTenant(ctx, accountID).OrderExpr(`"created_at" DESC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("automationRuleRepo.ListByAccount: %w", err)
	}
	return items, nil
}
