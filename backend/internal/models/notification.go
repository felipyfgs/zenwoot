package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Notification struct {
	bun.BaseModel    `bun:"table:notifications"`
	ID               int64      `bun:"id,pk,autoincrement"         json:"id"`
	UserID           int64      `bun:"user_id,notnull"             json:"user_id"`
	AccountID        int64      `bun:"account_id,notnull"          json:"account_id"`
	NotificationType string     `bun:"notification_type,notnull"  json:"notification_type"`
	Title            string     `bun:"title,notnull"               json:"title"`
	Body             *string    `bun:"body"                        json:"body"`
	PrimaryColor     *string    `bun:"primary_color"               json:"primary_color"`
	ActionURL        *string    `bun:"action_url"                 json:"action_url"`
	ReadAt           *time.Time `bun:"read_at"                    json:"read_at"`
	CreatedAt        time.Time  `bun:"created_at,notnull"         json:"created_at"`
}

type NotificationSetting struct {
	bun.BaseModel `bun:"table:notification_settings"`
	ID            int64     `bun:"id,pk,autoincrement"      json:"id"`
	UserID        int64     `bun:"user_id,notnull"          json:"user_id"`
	AccountID     int64     `bun:"account_id,notnull"       json:"account_id"`
	EmailEnabled  bool      `bun:"email_enabled,default:true" json:"email_enabled"`
	PushEnabled   bool      `bun:"push_enabled,default:true"  json:"push_enabled"`
	CreatedAt     time.Time `bun:"created_at,notnull"      json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"      json:"updated_at"`
}

type NotificationSubscription struct {
	bun.BaseModel          `bun:"table:notification_subscriptions"`
	ID                     int64          `bun:"id,pk,autoincrement"           json:"id"`
	UserID                 int64          `bun:"user_id,notnull"               json:"user_id"`
	AccountID              int64          `bun:"account_id,notnull"            json:"account_id"`
	SubscriptionType       string         `bun:"subscription_type,notnull"     json:"subscription_type"`
	SubscriptionAttributes map[string]any `bun:"subscription_attributes,type:jsonb" json:"subscription_attributes"`
	CreatedAt              time.Time      `bun:"created_at,notnull"             json:"created_at"`
}
