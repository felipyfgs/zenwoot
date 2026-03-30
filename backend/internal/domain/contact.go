package domain

import "time"

type Contact struct {
	ID         string    `json:"id"`
	AccountID  string    `json:"accountId"`
	Identifier string    `json:"identifier"`
	Name       string    `json:"name,omitempty"`
	PushName   string    `json:"pushName,omitempty"`
	AvatarURL  string    `json:"avatarUrl,omitempty"`
	IsBlocked  bool      `json:"isBlocked"`
	Metadata   any       `json:"metadata"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
