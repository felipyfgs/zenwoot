package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Campaign struct {
	bun.BaseModel  `bun:"table:campaigns"`
	ID             int64          `bun:"id,pk,autoincrement"       json:"id"`
	AccountID      int64          `bun:"account_id,notnull"         json:"account_id"`
	InboxID        int64          `bun:"inbox_id,notnull"           json:"inbox_id"`
	Name           *string        `bun:"name"                      json:"name"`
	Description    *string        `bun:"description"                json:"description"`
	CampaignType   int            `bun:"campaign_type,default:0"    json:"campaign_type"`
	CampaignStatus int            `bun:"campaign_status,default:0"  json:"campaign_status"`
	ScheduledAt    *time.Time     `bun:"scheduled_at"               json:"scheduled_at"`
	SentAt         *time.Time     `bun:"sent_at"                    json:"sent_at"`
	Message        *string        `bun:"message"                    json:"message"`
	TriggerRules   map[string]any `bun:"trigger_rules,type:jsonb"   json:"trigger_rules"`
	CreatedAt      time.Time      `bun:"created_at,notnull"         json:"created_at"`
	UpdatedAt      time.Time      `bun:"updated_at,notnull"         json:"updated_at"`
}
