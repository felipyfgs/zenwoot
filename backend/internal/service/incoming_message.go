package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"


	"wzap/internal/domain"
	"wzap/internal/processor"
	"wzap/internal/repo"
	"wzap/internal/ws"
)

// IncomingMessageProcessor persists incoming WhatsApp messages to the domain model.
// This is the Chatwoot pattern: Contact → ContactInbox → Conversation → Message.
type IncomingMessageProcessor struct {
	contactRepo      *repo.ContactRepository
	contactInboxRepo *repo.ContactInboxRepository
	convRepo         *repo.ConversationRepository
	msgRepo          *repo.MessageRepository
	inboxRepo        *repo.InboxRepository
	wsHub            *ws.Hub
	accountID        string // default account ID (single-tenant)
}

var _ processor.MessageProcessor = (*IncomingMessageProcessor)(nil)

func NewIncomingMessageProcessor(
	contactRepo *repo.ContactRepository,
	contactInboxRepo *repo.ContactInboxRepository,
	convRepo *repo.ConversationRepository,
	msgRepo *repo.MessageRepository,
	inboxRepo *repo.InboxRepository,
	wsHub *ws.Hub,
	accountID string,
) *IncomingMessageProcessor {
	return &IncomingMessageProcessor{
		contactRepo:      contactRepo,
		contactInboxRepo: contactInboxRepo,
		convRepo:         convRepo,
		msgRepo:          msgRepo,
		inboxRepo:        inboxRepo,
		wsHub:            wsHub,
		accountID:        accountID,
	}
}

func (p *IncomingMessageProcessor) Process(ctx context.Context, inboxID string, msg *domain.IncomingMessage) error {
	// Skip messages from me (echoes) - they're handled by outgoing message tracking
	if msg.IsFromMe {
		return nil
	}

	now := time.Now()

	// 1. Upsert contact (by phone/JID identifier)
	contact := &domain.Contact{
		ID:         uuid.NewString(),
		AccountID:  p.accountID,
		Identifier: msg.From,
		Name:       msg.PushName,
		PushName:   msg.PushName,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := p.contactRepo.Upsert(ctx, contact); err != nil {
		return fmt.Errorf("failed to upsert contact: %w", err)
	}

	// 2. Upsert contact_inbox (contact + inbox + sourceId)
	contactInbox := &domain.ContactInbox{
		ID:        uuid.NewString(),
		ContactID: contact.ID,
		InboxID:   inboxID,
		SourceID:  msg.From,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := p.contactInboxRepo.Upsert(ctx, contactInbox); err != nil {
		return fmt.Errorf("failed to upsert contact_inbox: %w", err)
	}

	// 3. Upsert conversation (inbox + identifier + contact)
	// The conversation identifier is the contact's JID (individual chat) or group JID
	conversationIdentifier := msg.From
	if msg.IsFromMe {
		// For echoes, the conversation is still with the recipient
		conversationIdentifier = msg.From
	}

	conv := &domain.Conversation{
		ID:          uuid.NewString(),
		AccountID:   p.accountID,
		InboxID:     inboxID,
		ContactID:   contact.ID,
		Identifier:  conversationIdentifier,
		LastMessage: msg.Content,
		LastMessageAt: &now,
		UnreadCount: 1,
		Status:      domain.ConversationOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := p.convRepo.Upsert(ctx, conv); err != nil {
		return fmt.Errorf("failed to upsert conversation: %w", err)
	}

	// 4. Create message
	message := &domain.Message{
		ID:           uuid.NewString(),
		AccountID:    p.accountID,
		ConversationID: conv.ID,
		InboxID:      inboxID,
		ContactID:    contact.ID,
		ExternalID:   msg.ExternalID,
		Direction:    domain.DirectionIncoming,
		ContentType:  msg.ContentType,
		Content:      msg.Content,
		MediaURL:     msg.MediaURL,
		MediaType:    msg.MediaType,
		Status:       domain.StatusDelivered,
		CreatedAt:    now,
	}
	if err := p.msgRepo.Create(ctx, message); err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// 5. Update conversation lastMessage (already done in Upsert with ON CONFLICT)
	// But we need to increment unreadCount for existing conversations
	// The Upsert already handles this via ON CONFLICT DO UPDATE

	// 6. Publish to WS Hub (real-time)
	if p.wsHub != nil {
		p.wsHub.Publish(p.accountID, inboxID, "message.created", message)
	}

	log.Info().
		Str("inboxId", inboxID).
		Str("contactId", contact.ID).
		Str("conversationId", conv.ID).
		Str("messageId", message.ID).
		Str("externalId", msg.ExternalID).
		Msg("Incoming message processed and persisted")

	return nil
}

// ProcessStatus updates message status based on delivery/read receipts.
func (p *IncomingMessageProcessor) ProcessStatus(ctx context.Context, inboxID string, status *domain.StatusUpdate) error {
	// Find message by external ID
	msg, err := p.msgRepo.FindByExternalID(ctx, status.ExternalID)
	if err != nil {
		// Message not found - might be an outgoing message we didn't persist yet
		log.Debug().Str("externalId", status.ExternalID).Msg("Message not found for status update")
		return nil
	}

	if err := p.msgRepo.UpdateStatus(ctx, msg.ID, status.Status); err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	// Publish to WS Hub
	if p.wsHub != nil {
		p.wsHub.Publish(p.accountID, inboxID, "message.status_updated", map[string]interface{}{
			"messageId": msg.ID,
			"status":    string(status.Status),
		})
	}

	log.Info().
		Str("inboxId", inboxID).
		Str("messageId", msg.ID).
		Str("status", string(status.Status)).
		Msg("Message status updated")

	return nil
}
