package services

import (
	"context"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type AutomationService struct {
	autoRepo *repo.AutomationRuleRepo
}

func NewAutomationService(autoRepo *repo.AutomationRuleRepo) *AutomationService {
	return &AutomationService{autoRepo: autoRepo}
}

func (s *AutomationService) List(ctx context.Context, accountID int64) ([]*models.AutomationRule, error) {
	return s.autoRepo.ListByAccount(ctx, accountID)
}

func (s *AutomationService) GetByID(ctx context.Context, accountID, id int64) (*models.AutomationRule, error) {
	return s.autoRepo.GetByID(ctx, accountID, id)
}

func (s *AutomationService) Create(ctx context.Context, m *models.AutomationRule) (*models.AutomationRule, error) {
	if err := s.autoRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *AutomationService) Update(ctx context.Context, m *models.AutomationRule) (*models.AutomationRule, error) {
	if err := s.autoRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *AutomationService) Delete(ctx context.Context, accountID, id int64) error {
	return s.autoRepo.Delete(ctx, accountID, id)
}
