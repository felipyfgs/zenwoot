package service

import (
	"context"

	"wzap/internal/domain"
	"wzap/internal/dto"
	"wzap/internal/repo"
)

type ConversationService struct {
	convRepo *repo.ConversationRepository
	msgRepo  *repo.MessageRepository
}

func NewConversationService(convRepo *repo.ConversationRepository, msgRepo *repo.MessageRepository) *ConversationService {
	return &ConversationService{convRepo: convRepo, msgRepo: msgRepo}
}

func (s *ConversationService) ListByInbox(ctx context.Context, inboxID string, req dto.PaginationReq) ([]domain.Conversation, error) {
	limit, offset := req.LimitOffset()
	return s.convRepo.FindByInboxID(ctx, inboxID, limit, offset)
}

func (s *ConversationService) Get(ctx context.Context, id string) (*domain.Conversation, error) {
	return s.convRepo.FindByID(ctx, id)
}

func (s *ConversationService) ToggleStatus(ctx context.Context, id string) (*domain.Conversation, error) {
	conv, err := s.convRepo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("conversation not found")
	}

	next := domain.ConversationOpen
	if conv.Status == domain.ConversationOpen {
		next = domain.ConversationResolved
	}

	if err := s.convRepo.UpdateStatus(ctx, id, next); err != nil {
		return nil, err
	}

	conv.Status = next
	return conv, nil
}

func (s *ConversationService) MarkRead(ctx context.Context, id string) error {
	return s.convRepo.MarkRead(ctx, id)
}

func (s *ConversationService) ListMessages(ctx context.Context, conversationID string, req dto.PaginationReq) ([]domain.Message, error) {
	limit, offset := req.LimitOffset()
	return s.msgRepo.FindByConversation(ctx, conversationID, limit, offset)
}
