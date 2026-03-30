package domain

import "time"

// Note represents a private agent note attached to a contact.
type Note struct {
	ID        string    `json:"id"`
	AccountID string    `json:"accountId"`
	ContactID string    `json:"contactId"`
	UserID    string    `json:"userId,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ConversationParticipant represents an agent who is subscribed/following a conversation.
type ConversationParticipant struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"accountId"`
	ConversationID string    `json:"conversationId"`
	UserID         string    `json:"userId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
