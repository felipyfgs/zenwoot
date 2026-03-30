package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Conversation struct {
	bun.BaseModel        `bun:"table:conversations"`
	ID                   int64          `bun:"id,pk,autoincrement"             json:"id"`
	DisplayID            int            `bun:"display_id,notnull"               json:"display_id"`
	AccountID            int64          `bun:"account_id,notnull"               json:"account_id"`
	InboxID              int64          `bun:"inbox_id,notnull"                 json:"inbox_id"`
	ContactID            *int64         `bun:"contact_id"                       json:"contact_id"`
	ContactInboxID       *int64         `bun:"contact_inbox_id"                  json:"contact_inbox_id"`
	AssigneeID           *int64         `bun:"assignee_id"                      json:"assignee_id"`
	TeamID               *int64         `bun:"team_id"                          json:"team_id"`
	CampaignID           *int64         `bun:"campaign_id"                      json:"campaign_id"`
	SlaPolicyID          *int64         `bun:"sla_policy_id"                     json:"sla_policy_id"`
	Status               int            `bun:"status,notnull,default:0"        json:"status"`
	Priority             *int           `bun:"priority"                        json:"priority"`
	UUID                 uuid.UUID      `bun:"uuid,notnull"                    json:"uuid"`
	Identifier           *string        `bun:"identifier"                      json:"identifier"`
	SnoozedUntil         *time.Time     `bun:"snoozed_until"                    json:"snoozed_until"`
	WaitingSince         *time.Time     `bun:"waiting_since"                    json:"waiting_since"`
	FirstReplyCreatedAt  *time.Time     `bun:"first_reply_at"                   json:"first_reply_at"`
	LastActivityAt       time.Time      `bun:"last_activity_at,notnull"          json:"last_activity_at"`
	ContactLastSeenAt    *time.Time     `bun:"contact_seen_at"                  json:"contact_seen_at"`
	AgentLastSeenAt      *time.Time     `bun:"agent_seen_at"                    json:"agent_seen_at"`
	AssigneeLastSeenAt   *time.Time     `bun:"assignee_seen_at"                 json:"assignee_seen_at"`
	AdditionalAttributes map[string]any `bun:"extra_attributes,type:jsonb"      json:"extra_attributes"`
	CustomAttributes     map[string]any `bun:"custom_attributes,type:jsonb"     json:"custom_attributes"`
	CachedLabelList      *string        `bun:"labels_cache"                     json:"labels_cache"`
	CreatedAt            time.Time      `bun:"created_at,notnull"               json:"created_at"`
	UpdatedAt            time.Time      `bun:"updated_at,notnull"               json:"updated_at"`

	Assignee *User    `bun:"rel:belongs-to,join:assignee_id=id" json:"assignee,omitempty"`
	Contact  *Contact `bun:"rel:belongs-to,join:contact_id=id"  json:"contact,omitempty"`
	Inbox    *Inbox   `bun:"rel:belongs-to,join:inbox_id=id"    json:"inbox,omitempty"`
}

const (
	ConvStatusOpen     = 0
	ConvStatusResolved = 1
	ConvStatusPending  = 2
	ConvStatusSnoozed  = 3
)
