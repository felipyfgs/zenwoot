package model

import "time"

type Webhook struct {
	ID          string    `json:"id"`
	InboxID     string    `json:"inboxId"`
	URL         string    `json:"url"`
	Secret      string    `json:"secret,omitempty"`
	Events      []string  `json:"events"`
	Enabled     bool      `json:"enabled"`
	NatsEnabled bool      `json:"natsEnabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
