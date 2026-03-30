package services

import (
	"context"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type TeamService struct {
	teamRepo *repo.TeamRepo
}

func NewTeamService(teamRepo *repo.TeamRepo) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) List(ctx context.Context, accountID int64) ([]*models.Team, error) {
	return s.teamRepo.ListByAccount(ctx, accountID)
}

func (s *TeamService) GetByID(ctx context.Context, accountID, id int64) (*models.Team, error) {
	return s.teamRepo.GetByID(ctx, accountID, id)
}

func (s *TeamService) Create(ctx context.Context, m *models.Team) (*models.Team, error) {
	if err := s.teamRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *TeamService) Update(ctx context.Context, m *models.Team) (*models.Team, error) {
	if err := s.teamRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *TeamService) Delete(ctx context.Context, accountID, id int64) error {
	return s.teamRepo.Delete(ctx, accountID, id)
}

func (s *TeamService) ListMembers(ctx context.Context, teamID int64) ([]*models.TeamMember, error) {
	return s.teamRepo.ListMembers(ctx, teamID)
}

func (s *TeamService) AddMember(ctx context.Context, teamID, userID int64) error {
	return s.teamRepo.AddMember(ctx, teamID, userID)
}

func (s *TeamService) RemoveMember(ctx context.Context, teamID, userID int64) error {
	return s.teamRepo.RemoveMember(ctx, teamID, userID)
}
