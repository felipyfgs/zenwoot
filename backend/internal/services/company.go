package services

import (
	"context"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type CompanyService struct {
	companyRepo *repo.CompanyRepo
}

func NewCompanyService(companyRepo *repo.CompanyRepo) *CompanyService {
	return &CompanyService{companyRepo: companyRepo}
}

func (s *CompanyService) List(ctx context.Context, accountID int64) ([]*models.Company, error) {
	return s.companyRepo.ListByAccount(ctx, accountID)
}

func (s *CompanyService) GetByID(ctx context.Context, accountID, id int64) (*models.Company, error) {
	return s.companyRepo.GetByID(ctx, accountID, id)
}

func (s *CompanyService) Search(ctx context.Context, accountID int64, q string) ([]*models.Company, error) {
	return s.companyRepo.Search(ctx, accountID, q)
}

func (s *CompanyService) Create(ctx context.Context, m *models.Company) (*models.Company, error) {
	if err := s.companyRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *CompanyService) Update(ctx context.Context, m *models.Company) (*models.Company, error) {
	if err := s.companyRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *CompanyService) Delete(ctx context.Context, accountID, id int64) error {
	return s.companyRepo.Delete(ctx, accountID, id)
}
