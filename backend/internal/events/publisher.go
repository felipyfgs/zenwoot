package events

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/felipyfgs/zenwoot/backend/internal/logger"
)

type Publisher struct {
	nc *nats.Conn
}

func NewPublisher(nc *nats.Conn) *Publisher {
	return &Publisher{nc: nc}
}

func (p *Publisher) Publish(ctx context.Context, subject string, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error().Err(err).Str("subject", subject).Msg("failed to marshal event")
		return err
	}

	err = p.nc.Publish(subject, data)
	if err != nil {
		logger.Error().Err(err).Str("subject", subject).Msg("failed to publish event")
		return err
	}

	logger.Debug().Str("subject", subject).Msg("event published")
	return nil
}

func (p *Publisher) PublishConversationCreated(ctx context.Context, accountID, conversationID int64) error {
	return p.Publish(ctx, "zenwoot.conversation.created", map[string]any{
		"account_id":      accountID,
		"conversation_id": conversationID,
	})
}

func (p *Publisher) PublishConversationUpdated(ctx context.Context, accountID, conversationID int64) error {
	return p.Publish(ctx, "zenwoot.conversation.updated", map[string]any{
		"account_id":      accountID,
		"conversation_id": conversationID,
	})
}

func (p *Publisher) PublishConversationResolved(ctx context.Context, accountID, conversationID int64) error {
	return p.Publish(ctx, "zenwoot.conversation.resolved", map[string]any{
		"account_id":      accountID,
		"conversation_id": conversationID,
	})
}

func (p *Publisher) PublishMessageCreated(ctx context.Context, accountID, conversationID, messageID int64) error {
	return p.Publish(ctx, "zenwoot.message.created", map[string]any{
		"account_id":      accountID,
		"conversation_id": conversationID,
		"message_id":      messageID,
	})
}

func (p *Publisher) PublishContactCreated(ctx context.Context, accountID, contactID int64) error {
	return p.Publish(ctx, "zenwoot.contact.created", map[string]any{
		"account_id": accountID,
		"contact_id": contactID,
	})
}
