package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Contact struct {
	bun.BaseModel        `bun:"table:contacts"`
	ID                   int64          `bun:"id,pk,autoincrement"             json:"id"`
	AccountID            int64          `bun:"account_id,notnull"               json:"account_id"`
	Name                 string         `bun:"name,default:''"                json:"name"`
	Email                *string        `bun:"email"                           json:"email"`
	PhoneNumber          *string        `bun:"phone_number"                     json:"phone_number"`
	Identifier           *string        `bun:"identifier"                      json:"identifier"`
	AvatarURL            *string        `bun:"avatar_url"                       json:"avatar_url"`
	AdditionalAttributes map[string]any `bun:"additional_attributes,type:jsonb" json:"additional_attributes"`
	CustomAttributes     map[string]any `bun:"custom_attributes,type:jsonb"     json:"custom_attributes"`
	LastActivityAt       *time.Time     `bun:"last_activity_at"                  json:"last_activity_at"`
	ContactType          int            `bun:"contact_type,default:0"           json:"contact_type"`
	FirstName            string         `bun:"first_name,default:''"          json:"first_name"`
	MiddleName           string         `bun:"middle_name,default:''"         json:"middle_name"`
	LastName             string         `bun:"last_name,default:''"           json:"last_name"`
	Location             string         `bun:"location,default:''"            json:"location"`
	CountryCode          string         `bun:"country_code,default:''"         json:"country_code"`
	Blocked              bool           `bun:"blocked,notnull,default:false"   json:"blocked"`
	CompanyID            *int64         `bun:"company_id"                       json:"company_id"`
	CreatedAt            time.Time      `bun:"created_at,notnull"               json:"created_at"`
	UpdatedAt            time.Time      `bun:"updated_at,notnull"               json:"updated_at"`
	DeletedAt            *time.Time     `bun:"deleted_at,soft_delete"           json:"deleted_at,omitempty"`
}

type ContactInbox struct {
	bun.BaseModel `bun:"table:contact_inboxes"`
	ID            int64     `bun:"id,pk,autoincrement"        json:"id"`
	ContactID     int64     `bun:"contact_id"                  json:"contact_id"`
	InboxID       int64     `bun:"inbox_id"                    json:"inbox_id"`
	SourceID      string    `bun:"source_id,notnull"           json:"source_id"`
	HmacVerified  bool      `bun:"hmac_verified,default:false" json:"hmac_verified"`
	PubsubToken   *string   `bun:"pubsub_token"                json:"pubsub_token"`
	CreatedAt     time.Time `bun:"created_at,notnull"          json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"          json:"updated_at"`

	Contact *Contact `bun:"rel:belongs-to,join:contact_id=id" json:"contact,omitempty"`
	Inbox   *Inbox   `bun:"rel:belongs-to,join:inbox_id=id"   json:"inbox,omitempty"`
}
