package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Webhook struct {
	bun.BaseModel `bun:"table:webhooks"`
	ID            int64      `bun:"id,pk,autoincrement"   json:"id"`
	AccountID     int64      `bun:"account_id,notnull"     json:"account_id"`
	InboxID       *int64     `bun:"inbox_id"               json:"inbox_id"`
	URL           string     `bun:"url,notnull"           json:"url"`
	Subscriptions []string   `bun:"subscriptions,array"   json:"subscriptions"`
	HmacToken     *string    `bun:"hmac_token"             json:"-"`
	CreatedAt     time.Time  `bun:"created_at,notnull"     json:"created_at"`
	UpdatedAt     time.Time  `bun:"updated_at,notnull"     json:"updated_at"`
	DeletedAt     *time.Time `bun:"deleted_at,soft_delete" json:"deleted_at,omitempty"`
}

type AutomationRule struct {
	bun.BaseModel `bun:"table:automation_rules"`
	ID            int64          `bun:"id,pk,autoincrement"          json:"id"`
	AccountID     int64          `bun:"account_id,notnull"            json:"account_id"`
	Name          string         `bun:"name,notnull"                 json:"name"`
	Description   *string        `bun:"description"                  json:"description"`
	EventName     string         `bun:"event_name,notnull"            json:"event_name"`
	Conditions    map[string]any `bun:"conditions,type:jsonb,notnull" json:"conditions"`
	Actions       map[string]any `bun:"actions,type:jsonb,notnull"    json:"actions"`
	Active        bool           `bun:"active,notnull,default:true"   json:"active"`
	CreatedAt     time.Time      `bun:"created_at,notnull"            json:"created_at"`
	UpdatedAt     time.Time      `bun:"updated_at,notnull"            json:"updated_at"`
	DeletedAt     *time.Time     `bun:"deleted_at,soft_delete"        json:"deleted_at,omitempty"`
}

type CannedResponse struct {
	bun.BaseModel `bun:"table:canned_responses"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	AccountID     int64     `bun:"account_id,notnull"   json:"account_id"`
	ShortCode     *string   `bun:"short_code"           json:"short_code"`
	Content       *string   `bun:"content"             json:"content"`
	CreatedAt     time.Time `bun:"created_at,notnull"   json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"   json:"updated_at"`
}
