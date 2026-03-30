package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"wzap/internal/domain"
	"wzap/internal/repo"
)

type TeamService struct {
	teamRepo      *repo.TeamRepository
	teamMemberRepo *repo.TeamMemberRepository
	inboxMemberRepo *repo.InboxMemberRepository
}

func NewTeamService(teamRepo *repo.TeamRepository, teamMemberRepo *repo.TeamMemberRepository, inboxMemberRepo *repo.InboxMemberRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, teamMemberRepo: teamMemberRepo, inboxMemberRepo: inboxMemberRepo}
}

type CreateTeamReq struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	AllowAutoAssign bool   `json:"allowAutoAssign"`
}

func (s *TeamService) Create(ctx context.Context, accountID string, req CreateTeamReq) (*domain.Team, error) {
	now := time.Now()
	team := &domain.Team{
		ID:             uuid.NewString(),
		AccountID:      accountID,
		Name:           req.Name,
		Description:    req.Description,
		AllowAutoAssign: req.AllowAutoAssign,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) Get(ctx context.Context, id string) (*domain.Team, error) {
	team, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("team not found")
	}
	return team, nil
}

func (s *TeamService) ListByAccount(ctx context.Context, accountID string) ([]domain.Team, error) {
	return s.teamRepo.FindByAccountID(ctx, accountID)
}

type UpdateTeamReq struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	AllowAutoAssign *bool  `json:"allowAutoAssign,omitempty"`
}

func (s *TeamService) Update(ctx context.Context, id string, req UpdateTeamReq) (*domain.Team, error) {
	team, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("team not found")
	}

	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Description != "" {
		team.Description = req.Description
	}
	if req.AllowAutoAssign != nil {
		team.AllowAutoAssign = *req.AllowAutoAssign
	}

	if err := s.teamRepo.Update(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) Delete(ctx context.Context, id string) error {
	return s.teamRepo.Delete(ctx, id)
}

// TeamMember operations

func (s *TeamService) AddMember(ctx context.Context, teamID, userID string, role domain.UserRole) error {
	now := time.Now()
	member := &domain.TeamMember{
		ID:        uuid.NewString(),
		TeamID:    teamID,
		UserID:    userID,
		Role:      role,
		CreatedAt: now,
	}
	return s.teamMemberRepo.Create(ctx, member)
}

func (s *TeamService) ListMembers(ctx context.Context, teamID string) ([]domain.TeamMember, error) {
	return s.teamMemberRepo.FindByTeamID(ctx, teamID)
}

func (s *TeamService) RemoveMember(ctx context.Context, teamID, userID string) error {
	return s.teamMemberRepo.Delete(ctx, teamID, userID)
}

// InboxMember operations

func (s *TeamService) AddInboxMember(ctx context.Context, inboxID, userID string) error {
	now := time.Now()
	member := &domain.InboxMember{
		ID:        uuid.NewString(),
		InboxID:   inboxID,
		UserID:    userID,
		CreatedAt: now,
	}
	return s.inboxMemberRepo.Create(ctx, member)
}

func (s *TeamService) ListInboxMembers(ctx context.Context, inboxID string) ([]domain.InboxMember, error) {
	return s.inboxMemberRepo.FindByInboxID(ctx, inboxID)
}

func (s *TeamService) RemoveInboxMember(ctx context.Context, inboxID, userID string) error {
	return s.inboxMemberRepo.Delete(ctx, inboxID, userID)
}

func (s *TeamService) GetUserInboxes(ctx context.Context, userID string) ([]domain.InboxMember, error) {
	return s.inboxMemberRepo.FindByUserID(ctx, userID)
}
