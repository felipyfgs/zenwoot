package domain

import "time"

// CannedResponse represents a preset/template reply agents can quickly insert.
type CannedResponse struct {
	ID        string    `json:"id"`
	AccountID string    `json:"accountId"`
	ShortCode string    `json:"shortCode"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
