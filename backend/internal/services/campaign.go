package services

import (
	"context"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type CampaignService struct {
	campaignRepo *repo.CampaignRepo
}

func NewCampaignService(campaignRepo *repo.CampaignRepo) *CampaignService {
	return &CampaignService{campaignRepo: campaignRepo}
}

func (s *CampaignService) List(ctx context.Context, accountID int64) ([]*models.Campaign, error) {
	return s.campaignRepo.ListByAccount(ctx, accountID)
}

func (s *CampaignService) ListByInbox(ctx context.Context, accountID, inboxID int64) ([]*models.Campaign, error) {
	return s.campaignRepo.ListByInbox(ctx, accountID, inboxID)
}

func (s *CampaignService) GetByID(ctx context.Context, accountID, id int64) (*models.Campaign, error) {
	return s.campaignRepo.GetByID(ctx, accountID, id)
}

func (s *CampaignService) Create(ctx context.Context, m *models.Campaign) (*models.Campaign, error) {
	if err := s.campaignRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *CampaignService) Update(ctx context.Context, m *models.Campaign) (*models.Campaign, error) {
	if err := s.campaignRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *CampaignService) Delete(ctx context.Context, accountID, id int64) error {
	return s.campaignRepo.Delete(ctx, accountID, id)
}
