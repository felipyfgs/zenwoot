package services

import (
	"context"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type MacroService struct {
	macroRepo *repo.MacroRepo
}

func NewMacroService(macroRepo *repo.MacroRepo) *MacroService {
	return &MacroService{macroRepo: macroRepo}
}

func (s *MacroService) List(ctx context.Context, accountID int64) ([]*models.Macro, error) {
	return s.macroRepo.ListByAccount(ctx, accountID)
}

func (s *MacroService) GetByID(ctx context.Context, accountID, id int64) (*models.Macro, error) {
	return s.macroRepo.GetByID(ctx, accountID, id)
}

func (s *MacroService) Create(ctx context.Context, m *models.Macro) (*models.Macro, error) {
	if err := s.macroRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MacroService) Update(ctx context.Context, m *models.Macro) (*models.Macro, error) {
	if err := s.macroRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MacroService) Delete(ctx context.Context, accountID, id int64) error {
	return s.macroRepo.Delete(ctx, accountID, id)
}
