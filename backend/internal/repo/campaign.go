package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type CampaignRepo struct {
	BaseRepo[models.Campaign]
}

func NewCampaignRepo(db *bun.DB) *CampaignRepo {
	return &CampaignRepo{BaseRepo: *NewBaseRepo[models.Campaign](db)}
}

func (r *CampaignRepo) ListByAccount(ctx context.Context, accountID int64) ([]*models.Campaign, error) {
	var items []*models.Campaign
	err := r.WithTenant(ctx, accountID).OrderExpr(`"created_at" DESC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("campaignRepo.ListByAccount: %w", err)
	}
	return items, nil
}

func (r *CampaignRepo) ListByInbox(ctx context.Context, accountID, inboxID int64) ([]*models.Campaign, error) {
	var items []*models.Campaign
	err := r.WithTenant(ctx, accountID).Where(`"inbox_id" = ?`, inboxID).OrderExpr(`"created_at" DESC`).Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("campaignRepo.ListByInbox: %w", err)
	}
	return items, nil
}

func (r *CampaignRepo) Create(ctx context.Context, m *models.Campaign) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("campaignRepo.Create: %w", err)
	}
	return nil
}

func (r *CampaignRepo) Update(ctx context.Context, m *models.Campaign) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("campaignRepo.Update: %w", err)
	}
	return nil
}

func (r *CampaignRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"campaigns"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("campaignRepo.Delete: %w", err)
	}
	return nil
}

func (r *CampaignRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Campaign, error) {
	var m models.Campaign
	err := r.WithTenant(ctx, accountID).Where(`"id" = ?`, id).Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("campaignRepo.GetByID: %w", err)
	}
	return &m, nil
}
