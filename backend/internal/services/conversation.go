package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type ConversationService struct {
	convRepo *repo.ConversationRepo
	nc       *nats.Conn
}

func NewConversationService(convRepo *repo.ConversationRepo, nc *nats.Conn) *ConversationService {
	return &ConversationService{convRepo: convRepo, nc: nc}
}

func (s *ConversationService) GetByID(ctx context.Context, accountID, id int64) (*models.Conversation, error) {
	return s.convRepo.GetByID(ctx, accountID, id)
}

func (s *ConversationService) List(ctx context.Context, accountID int64, f repo.ConversationFilter) ([]*models.Conversation, int, error) {
	return s.convRepo.List(ctx, accountID, f)
}

func (s *ConversationService) Create(ctx context.Context, m *models.Conversation) (*models.Conversation, error) {
	if err := s.convRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	s.publish("zenwoot.conversation.created", m)
	return m, nil
}

func (s *ConversationService) Assign(ctx context.Context, accountID, id int64, assigneeID *int64, teamID *int64) (*models.Conversation, error) {
	conv, err := s.convRepo.GetByID(ctx, accountID, id)
	if err != nil {
		return nil, err
	}
	conv.AssigneeID = assigneeID
	conv.TeamID = teamID
	if err := s.convRepo.Update(ctx, conv); err != nil {
		return nil, err
	}
	s.publish("zenwoot.conversation.updated", conv)
	return conv, nil
}

func (s *ConversationService) Resolve(ctx context.Context, accountID, id int64) (*models.Conversation, error) {
	conv, err := s.convRepo.GetByID(ctx, accountID, id)
	if err != nil {
		return nil, err
	}
	conv.Status = models.ConvStatusResolved
	if err := s.convRepo.Update(ctx, conv); err != nil {
		return nil, err
	}
	s.publish("zenwoot.conversation.resolved", conv)
	return conv, nil
}

func (s *ConversationService) Reopen(ctx context.Context, accountID, id int64) (*models.Conversation, error) {
	conv, err := s.convRepo.GetByID(ctx, accountID, id)
	if err != nil {
		return nil, err
	}
	conv.Status = models.ConvStatusOpen
	conv.WaitingSince = nil
	if err := s.convRepo.Update(ctx, conv); err != nil {
		return nil, err
	}
	s.publish("zenwoot.conversation.updated", conv)
	return conv, nil
}

func (s *ConversationService) Snooze(ctx context.Context, accountID, id int64, until time.Time) (*models.Conversation, error) {
	conv, err := s.convRepo.GetByID(ctx, accountID, id)
	if err != nil {
		return nil, err
	}
	conv.Status = models.ConvStatusSnoozed
	conv.SnoozedUntil = &until
	if err := s.convRepo.Update(ctx, conv); err != nil {
		return nil, err
	}
	s.publish("zenwoot.conversation.updated", conv)
	return conv, nil
}

func (s *ConversationService) publish(subject string, payload any) {
	if s.nc == nil {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_ = s.nc.Publish(subject, data)
}
