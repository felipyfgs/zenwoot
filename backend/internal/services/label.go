package services

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type LabelService struct {
	labelRepo *repo.LabelRepo
	db        *bun.DB
}

func NewLabelService(labelRepo *repo.LabelRepo, db *bun.DB) *LabelService {
	return &LabelService{labelRepo: labelRepo, db: db}
}

func (s *LabelService) List(ctx context.Context, accountID int64) ([]*models.Label, error) {
	return s.labelRepo.ListByAccount(ctx, accountID)
}

func (s *LabelService) GetByID(ctx context.Context, accountID, id int64) (*models.Label, error) {
	return s.labelRepo.GetByID(ctx, accountID, id)
}

func (s *LabelService) Create(ctx context.Context, m *models.Label) (*models.Label, error) {
	if err := s.labelRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *LabelService) Update(ctx context.Context, m *models.Label) (*models.Label, error) {
	if err := s.labelRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *LabelService) Delete(ctx context.Context, accountID, id int64) error {
	return s.labelRepo.Delete(ctx, accountID, id)
}

func (s *LabelService) AddToConversation(ctx context.Context, accountID, conversationID, labelID int64) error {
	cl := &models.ConversationLabel{ConversationID: conversationID, LabelID: labelID}
	_, err := s.db.NewInsert().Model(cl).On("CONFLICT DO NOTHING").Exec(ctx)
	return err
}

func (s *LabelService) RemoveFromConversation(ctx context.Context, conversationID, labelID int64) error {
	_, err := s.db.NewDelete().TableExpr(`"conversation_labels"`).
		Where(`"conversation_id" = ? AND "label_id" = ?`, conversationID, labelID).
		Exec(ctx)
	return err
}
