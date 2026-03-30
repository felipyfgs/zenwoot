package models

import (
	"time"

	"github.com/uptrace/bun"
)

type CustomFilter struct {
	bun.BaseModel `bun:"table:custom_filters"`
	ID            int64          `bun:"id,pk,autoincrement"     json:"id"`
	AccountID     int64          `bun:"account_id,notnull"      json:"account_id"`
	UserID        int64          `bun:"user_id,notnull"        json:"user_id"`
	Name          string         `bun:"name,notnull"           json:"name"`
	FilterType    string         `bun:"filter_type,notnull"    json:"filter_type"`
	Query         map[string]any `bun:"query,type:jsonb"       json:"query"`
	CreatedAt     time.Time      `bun:"created_at,notnull"     json:"created_at"`
	UpdatedAt     time.Time      `bun:"updated_at,notnull"     json:"updated_at"`
}
