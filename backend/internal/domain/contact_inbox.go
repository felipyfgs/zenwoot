package domain

import "time"

// ContactInbox links a contact to a specific inbox with a channel-specific sourceId.
// This is the Chatwoot pattern: a contact can have multiple identities across inboxes.
type ContactInbox struct {
	ID        string    `json:"id"`
	ContactID string    `json:"contactId"`
	InboxID   string    `json:"inboxId"`
	SourceID  string    `json:"sourceId"` // Channel-specific identifier (JID, phone, email)
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
