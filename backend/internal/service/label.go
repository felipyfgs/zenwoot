package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"wzap/internal/domain"
	"wzap/internal/repo"
)

type LabelService struct {
	labelRepo             *repo.LabelRepository
	conversationLabelRepo *repo.ConversationLabelRepository
}

func NewLabelService(labelRepo *repo.LabelRepository, conversationLabelRepo *repo.ConversationLabelRepository) *LabelService {
	return &LabelService{labelRepo: labelRepo, conversationLabelRepo: conversationLabelRepo}
}

type CreateLabelReq struct {
	Title       string `json:"title"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

func (s *LabelService) Create(ctx context.Context, accountID string, req CreateLabelReq) (*domain.Label, error) {
	color := req.Color
	if color == "" {
		color = "#1F93FF" // default blue
	}

	now := time.Now()
	label := &domain.Label{
		ID:          uuid.NewString(),
		AccountID:   accountID,
		Title:       req.Title,
		Color:       color,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.labelRepo.Create(ctx, label); err != nil {
		return nil, err
	}
	return label, nil
}

func (s *LabelService) Get(ctx context.Context, id string) (*domain.Label, error) {
	label, err := s.labelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("label not found")
	}
	return label, nil
}

func (s *LabelService) ListByAccount(ctx context.Context, accountID string) ([]domain.Label, error) {
	return s.labelRepo.FindByAccountID(ctx, accountID)
}

type UpdateLabelReq struct {
	Title       string `json:"title,omitempty"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

func (s *LabelService) Update(ctx context.Context, id string, req UpdateLabelReq) (*domain.Label, error) {
	label, err := s.labelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, notFoundErrorf("label not found")
	}

	if req.Title != "" {
		label.Title = req.Title
	}
	if req.Color != "" {
		label.Color = req.Color
	}
	if req.Description != "" {
		label.Description = req.Description
	}

	if err := s.labelRepo.Update(ctx, label); err != nil {
		return nil, err
	}
	return label, nil
}

func (s *LabelService) Delete(ctx context.Context, id string) error {
	return s.labelRepo.Delete(ctx, id)
}

// Conversation Label operations

func (s *LabelService) AddToConversation(ctx context.Context, conversationID, labelID string) error {
	now := time.Now()
	cl := &domain.ConversationLabel{
		ID:             uuid.NewString(),
		ConversationID: conversationID,
		LabelID:        labelID,
		CreatedAt:      now,
	}
	return s.conversationLabelRepo.Create(ctx, cl)
}

func (s *LabelService) GetConversationLabels(ctx context.Context, conversationID string) ([]domain.Label, error) {
	return s.conversationLabelRepo.FindByConversationID(ctx, conversationID)
}

func (s *LabelService) RemoveFromConversation(ctx context.Context, conversationID, labelID string) error {
	return s.conversationLabelRepo.Delete(ctx, conversationID, labelID)
}

func (s *LabelService) SetConversationLabels(ctx context.Context, conversationID string, labelIDs []string) error {
	// First remove all existing labels
	if err := s.conversationLabelRepo.DeleteAllByConversation(ctx, conversationID); err != nil {
		return err
	}

	// Then add the new ones
	now := time.Now()
	for _, labelID := range labelIDs {
		cl := &domain.ConversationLabel{
			ID:             uuid.NewString(),
			ConversationID: conversationID,
			LabelID:        labelID,
			CreatedAt:      now,
		}
		if err := s.conversationLabelRepo.Create(ctx, cl); err != nil {
			return err
		}
	}
	return nil
}
