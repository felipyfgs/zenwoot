package services

import (
	"context"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type AgentBotService struct {
	agentBotRepo *repo.AgentBotRepo
}

func NewAgentBotService(agentBotRepo *repo.AgentBotRepo) *AgentBotService {
	return &AgentBotService{agentBotRepo: agentBotRepo}
}

func (s *AgentBotService) List(ctx context.Context, accountID int64) ([]*models.AgentBot, error) {
	return s.agentBotRepo.ListByAccount(ctx, accountID)
}

func (s *AgentBotService) GetByID(ctx context.Context, accountID, id int64) (*models.AgentBot, error) {
	return s.agentBotRepo.GetByID(ctx, accountID, id)
}

func (s *AgentBotService) Create(ctx context.Context, m *models.AgentBot) (*models.AgentBot, error) {
	if err := s.agentBotRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *AgentBotService) Update(ctx context.Context, m *models.AgentBot) (*models.AgentBot, error) {
	if err := s.agentBotRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *AgentBotService) Delete(ctx context.Context, accountID, id int64) error {
	return s.agentBotRepo.Delete(ctx, accountID, id)
}
