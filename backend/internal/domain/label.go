package domain

import "time"

// Label represents a tag for organizing conversations and contacts.
type Label struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"accountId"`
	Title       string    `json:"title"`
	Color       string    `json:"color"` // hex color code
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ConversationLabel links labels to conversations.
type ConversationLabel struct {
	ID             string    `json:"id"`
	ConversationID string   `json:"conversationId"`
	LabelID        string   `json:"labelId"`
	CreatedAt      time.Time `json:"createdAt"`
}
