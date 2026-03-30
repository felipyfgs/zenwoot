package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Inbox struct {
	bun.BaseModel              `bun:"table:inboxes"`
	ID                         int64          `bun:"id,pk,autoincrement"                     json:"id"`
	AccountID                  int64          `bun:"account_id,notnull"                       json:"account_id"`
	ChannelID                  int64          `bun:"channel_id,notnull"                       json:"channel_id"`
	ChannelType                string         `bun:"channel_type,notnull"                     json:"channel_type"`
	Name                       string         `bun:"name,notnull"                            json:"name"`
	EnableAutoAssignment       bool           `bun:"auto_assign_enabled,default:true"          json:"auto_assign_enabled"`
	GreetingEnabled            bool           `bun:"greeting_enabled,default:false"           json:"greeting_enabled"`
	GreetingMessage            *string        `bun:"greeting_message"                         json:"greeting_message"`
	EmailAddress               *string        `bun:"email_address"                            json:"email_address"`
	WorkingHoursEnabled        bool           `bun:"working_hours_enabled,default:false"      json:"working_hours_enabled"`
	OutOfOfficeMessage         *string        `bun:"out_of_office_message"                      json:"out_of_office_message"`
	Timezone                   string         `bun:"timezone,default:'UTC'"                  json:"timezone"`
	EnableEmailCollect         bool           `bun:"enable_email_collect,default:true"         json:"enable_email_collect"`
	CsatSurveyEnabled          bool           `bun:"csat_enabled,default:false"                json:"csat_enabled"`
	AllowMessagesAfterResolved bool           `bun:"allow_resolved_messages,default:true"      json:"allow_resolved_messages"`
	LockToSingleConversation   bool           `bun:"single_conversation_only,notnull"         json:"single_conversation_only"`
	AutoAssignmentConfig       map[string]any `bun:"auto_assign_config,type:jsonb"            json:"auto_assign_config"`
	CsatConfig                 map[string]any `bun:"csat_config,type:jsonb,notnull"           json:"csat_config"`
	SenderNameType             int            `bun:"sender_name_type,default:0"                json:"sender_name_type"`
	BusinessName               *string        `bun:"business_name"                            json:"business_name"`
	CreatedAt                  time.Time      `bun:"created_at,notnull"                       json:"created_at"`
	UpdatedAt                  time.Time      `bun:"updated_at,notnull"                       json:"updated_at"`
	DeletedAt                  *time.Time     `bun:"deleted_at,soft_delete"                   json:"deleted_at,omitempty"`
}

type InboxMember struct {
	bun.BaseModel `bun:"table:inbox_members"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	InboxID       int64     `bun:"inbox_id,notnull"     json:"inbox_id"`
	UserID        int64     `bun:"user_id,notnull"      json:"user_id"`
	CreatedAt     time.Time `bun:"created_at,notnull"   json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"   json:"updated_at"`

	User *User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}
