package processor

import (
	"context"

	"wzap/internal/domain"
)

// MessageProcessor persists incoming WhatsApp messages to the domain model.
type MessageProcessor interface {
	Process(ctx context.Context, inboxID string, msg *domain.IncomingMessage) error
	ProcessStatus(ctx context.Context, inboxID string, status *domain.StatusUpdate) error
}
