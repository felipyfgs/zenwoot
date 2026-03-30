package domain

import "time"

type ConversationStatus string

const (
	ConversationOpen     ConversationStatus = "open"
	ConversationResolved ConversationStatus = "resolved"
	ConversationPending  ConversationStatus = "pending"
	ConversationSnoozed  ConversationStatus = "snoozed"
)

type ConversationPriority string

const (
	PriorityLow    ConversationPriority = "low"
	PriorityMedium ConversationPriority = "medium"
	PriorityHigh   ConversationPriority = "high"
	PriorityUrgent ConversationPriority = "urgent"
)

type Conversation struct {
	ID             string               `json:"id"`
	AccountID      string               `json:"accountId"`
	InboxID        string               `json:"inboxId"`
	ContactID      string               `json:"contactId,omitempty"`
	AssigneeID     string               `json:"assigneeId,omitempty"`
	TeamID         string               `json:"teamId,omitempty"`
	Identifier     string               `json:"identifier"`
	LastMessage    string               `json:"lastMessage,omitempty"`
	LastMessageAt  *time.Time           `json:"lastMessageAt,omitempty"`
	LastActivityAt *time.Time           `json:"lastActivityAt,omitempty"`
	UnreadCount    int                  `json:"unreadCount"`
	Status         ConversationStatus   `json:"status"`
	Priority       ConversationPriority `json:"priority"`
	Muted          bool                 `json:"muted"`
	SnoozedUntil   *time.Time           `json:"snoozedUntil,omitempty"`
	Metadata       any                  `json:"metadata"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
}
