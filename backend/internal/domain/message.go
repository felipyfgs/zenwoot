package domain

import "time"

type MessageDirection string

const (
	DirectionIncoming MessageDirection = "incoming"
	DirectionOutgoing MessageDirection = "outgoing"
)

type ContentType string

const (
	ContentText     ContentType = "text"
	ContentImage    ContentType = "image"
	ContentVideo    ContentType = "video"
	ContentAudio    ContentType = "audio"
	ContentDocument ContentType = "document"
	ContentLocation ContentType = "location"
	ContentContact  ContentType = "contact"
	ContentSticker  ContentType = "sticker"
)

type MessageStatus string

const (
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
	StatusFailed    MessageStatus = "failed"
)

type Message struct {
	ID             string           `json:"id"`
	AccountID      string           `json:"accountId"`
	ConversationID string           `json:"conversationId"`
	InboxID        string           `json:"inboxId"`
	ContactID      string           `json:"contactId,omitempty"`
	SenderID       string           `json:"senderId,omitempty"` // User ID for outgoing messages
	ExternalID     string           `json:"externalId,omitempty"`
	ReplyToID      string           `json:"replyToId,omitempty"` // For quoted/reply messages
	Direction      MessageDirection `json:"direction"`
	ContentType    ContentType      `json:"contentType"`
	Content        string           `json:"content,omitempty"`
	MediaURL       string           `json:"mediaUrl,omitempty"`
	MediaType      string           `json:"mediaType,omitempty"`
	Metadata       any              `json:"metadata"`
	Status         MessageStatus    `json:"status"`
	CreatedAt      time.Time        `json:"createdAt"`
}
