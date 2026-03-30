package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"wzap/internal/domain"
	"wzap/internal/repo"
	"wzap/internal/ws"
)

type NoteService struct {
	noteRepo *repo.NoteRepository
	wsHub    *ws.Hub
	accountID string
}

func NewNoteService(noteRepo *repo.NoteRepository, wsHub *ws.Hub, accountID string) *NoteService {
	return &NoteService{noteRepo: noteRepo, wsHub: wsHub, accountID: accountID}
}

type CreateNoteReq struct {
	Content string `json:"content"`
	UserID  string `json:"userId,omitempty"`
}

func (s *NoteService) Create(ctx context.Context, contactID string, req CreateNoteReq) (*domain.Note, error) {
	now := time.Now()
	note := &domain.Note{
		ID:        uuid.NewString(),
		AccountID: s.accountID,
		ContactID: contactID,
		UserID:    req.UserID,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.noteRepo.Create(ctx, note); err != nil {
		return nil, err
	}

	if s.wsHub != nil {
		s.wsHub.Publish(s.accountID, "", "note.created", map[string]interface{}{
			"noteId":    note.ID,
			"contactId": contactID,
			"content":   req.Content,
		})
	}

	return note, nil
}

func (s *NoteService) ListByContact(ctx context.Context, contactID string) ([]domain.Note, error) {
	return s.noteRepo.FindByContactID(ctx, contactID)
}

func (s *NoteService) Delete(ctx context.Context, id string) error {
	note, err := s.noteRepo.FindByID(ctx, id)
	if err != nil {
		return notFoundErrorf("note not found")
	}

	if err := s.noteRepo.Delete(ctx, id); err != nil {
		return err
	}

	if s.wsHub != nil {
		s.wsHub.Publish(s.accountID, "", "note.deleted", map[string]interface{}{
			"noteId":    id,
			"contactId": note.ContactID,
		})
	}

	return nil
}

// ConversationParticipantService

type ConversationParticipantService struct {
	participantRepo *repo.ConversationParticipantRepository
	wsHub           *ws.Hub
	accountID       string
}

func NewConversationParticipantService(
	participantRepo *repo.ConversationParticipantRepository,
	wsHub *ws.Hub,
	accountID string,
) *ConversationParticipantService {
	return &ConversationParticipantService{
		participantRepo: participantRepo,
		wsHub:           wsHub,
		accountID:       accountID,
	}
}

func (s *ConversationParticipantService) Add(ctx context.Context, conversationID, userID, inboxID string) error {
	now := time.Now()
	p := &domain.ConversationParticipant{
		ID:             uuid.NewString(),
		AccountID:      s.accountID,
		ConversationID: conversationID,
		UserID:         userID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.participantRepo.Create(ctx, p); err != nil {
		return err
	}

	if s.wsHub != nil {
		s.wsHub.Publish(s.accountID, inboxID, "conversation.participant_added", map[string]interface{}{
			"conversationId": conversationID,
			"userId":         userID,
		})
	}
	return nil
}

func (s *ConversationParticipantService) Remove(ctx context.Context, conversationID, userID, inboxID string) error {
	if err := s.participantRepo.Delete(ctx, conversationID, userID); err != nil {
		return err
	}

	if s.wsHub != nil {
		s.wsHub.Publish(s.accountID, inboxID, "conversation.participant_removed", map[string]interface{}{
			"conversationId": conversationID,
			"userId":         userID,
		})
	}
	return nil
}

func (s *ConversationParticipantService) List(ctx context.Context, conversationID string) ([]domain.ConversationParticipant, error) {
	return s.participantRepo.FindByConversationID(ctx, conversationID)
}
