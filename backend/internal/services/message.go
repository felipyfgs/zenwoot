package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type MessageService struct {
	msgRepo  *repo.MessageRepo
	convRepo *repo.ConversationRepo
	db       *bun.DB
	nc       *nats.Conn
}

func NewMessageService(msgRepo *repo.MessageRepo, convRepo *repo.ConversationRepo, db *bun.DB, nc *nats.Conn) *MessageService {
	return &MessageService{msgRepo: msgRepo, convRepo: convRepo, db: db, nc: nc}
}

func (s *MessageService) List(ctx context.Context, accountID, conversationID int64, before *int64, limit int) ([]*models.Message, error) {
	return s.msgRepo.ListByConversation(ctx, accountID, conversationID, before, limit)
}

func (s *MessageService) GetByID(ctx context.Context, accountID, id int64) (*models.Message, error) {
	return s.msgRepo.GetByID(ctx, accountID, id)
}

func (s *MessageService) Update(ctx context.Context, m *models.Message) (*models.Message, error) {
	if err := s.msgRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MessageService) Delete(ctx context.Context, accountID, id int64) error {
	return s.msgRepo.Delete(ctx, accountID, id)
}

type CreateMessageInput struct {
	ConversationID int64
	AccountID      int64
	InboxID        int64
	SenderType     string
	SenderID       int64
	Content        string
	MessageType    int
	Private        bool
	ContentAttrs   map[string]any
}

func (s *MessageService) Create(ctx context.Context, in CreateMessageInput) (*models.Message, error) {
	msg := &models.Message{
		ConversationID:    in.ConversationID,
		AccountID:         in.AccountID,
		InboxID:           in.InboxID,
		SenderType:        &in.SenderType,
		SenderID:          &in.SenderID,
		Content:           &in.Content,
		MessageType:       in.MessageType,
		Private:           in.Private,
		ContentAttributes: in.ContentAttrs,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("messageService.Create begin tx: %w", err)
	}
	defer tx.Rollback()

	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	now := time.Now()
	_, err = tx.NewUpdate().TableExpr(`"conversations"`).
		Set(`"last_activity_at" = ?`, now).
		Where(`"id" = ? AND "account_id" = ?`, in.ConversationID, in.AccountID).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("messageService.Create update last_activity_at: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("messageService.Create commit: %w", err)
	}

	s.publish("zenwoot.message.created", msg)
	return msg, nil
}

func (s *MessageService) publish(subject string, payload any) {
	if s.nc == nil {
		return
	}
	data, _ := json.Marshal(payload)
	_ = s.nc.Publish(subject, data)
}
