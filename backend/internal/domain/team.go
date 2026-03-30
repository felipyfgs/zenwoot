package domain

import "time"

// Team represents a group of agents within an account.
type Team struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"accountId"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	AllowAutoAssign bool      `json:"allowAutoAssign"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// TeamMember links a user to a team.
type TeamMember struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	UserID    string    `json:"userId"`
	Role      UserRole  `json:"role"` // agent, team_lead
	CreatedAt time.Time `json:"createdAt"`
}

// InboxMember links a user to an inbox (agent assignment).
type InboxMember struct {
	ID             string     `json:"id"`
	InboxID        string     `json:"inboxId"`
	UserID         string     `json:"userId"`
	LastAssignedAt *time.Time `json:"lastAssignedAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
}
