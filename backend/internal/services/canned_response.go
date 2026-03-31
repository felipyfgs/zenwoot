package services

import (
	"context"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type CannedResponseService struct {
	cannedRepo *repo.CannedResponseRepo
}

func NewCannedResponseService(cannedRepo *repo.CannedResponseRepo) *CannedResponseService {
	return &CannedResponseService{cannedRepo: cannedRepo}
}

func (s *CannedResponseService) List(ctx context.Context, accountID int64) ([]*models.CannedResponse, error) {
	return s.cannedRepo.ListByAccount(ctx, accountID)
}

func (s *CannedResponseService) GetByID(ctx context.Context, accountID, id int64) (*models.CannedResponse, error) {
	return s.cannedRepo.GetByID(ctx, accountID, id)
}

func (s *CannedResponseService) GetByShortCode(ctx context.Context, accountID int64, shortCode string) (*models.CannedResponse, error) {
	return s.cannedRepo.GetByShortCode(ctx, accountID, shortCode)
}

func (s *CannedResponseService) Create(ctx context.Context, m *models.CannedResponse) (*models.CannedResponse, error) {
	if err := s.cannedRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *CannedResponseService) Update(ctx context.Context, m *models.CannedResponse) (*models.CannedResponse, error) {
	if err := s.cannedRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *CannedResponseService) Delete(ctx context.Context, accountID, id int64) error {
	return s.cannedRepo.Delete(ctx, accountID, id)
}
