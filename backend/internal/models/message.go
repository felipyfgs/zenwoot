package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Message struct {
	bun.BaseModel           `bun:"table:messages"`
	ID                      int64          `bun:"id,pk,autoincrement"             json:"id"`
	AccountID               int64          `bun:"account_id,notnull"               json:"account_id"`
	InboxID                 int64          `bun:"inbox_id,notnull"                 json:"inbox_id"`
	ConversationID          int64          `bun:"conversation_id,notnull"          json:"conversation_id"`
	SenderType              *string        `bun:"sender_type"                      json:"sender_type"`
	SenderID                *int64         `bun:"sender_id"                        json:"sender_id"`
	Content                 *string        `bun:"content"                         json:"content"`
	ContentType             int            `bun:"content_type,notnull,default:0"   json:"content_type"`
	MessageType             int            `bun:"message_type,notnull"             json:"message_type"`
	Status                  int            `bun:"status,default:0"                json:"status"`
	SourceID                *string        `bun:"source_id"                        json:"source_id"`
	Private                 bool           `bun:"private,notnull,default:false"   json:"private"`
	ContentAttributes       map[string]any `bun:"content_meta,type:jsonb"           json:"content_meta"`
	AdditionalAttributes    map[string]any `bun:"extra_attributes,type:jsonb"       json:"extra_attributes"`
	ExternalSourceIds       map[string]any `bun:"external_ids,type:jsonb"            json:"external_ids"`
	ProcessedMessageContent *string        `bun:"processed_content"                 json:"processed_content"`
	Sentiment               map[string]any `bun:"sentiment,type:jsonb"            json:"sentiment"`
	CreatedAt               time.Time      `bun:"created_at,notnull"               json:"created_at"`
	UpdatedAt               time.Time      `bun:"updated_at,notnull"               json:"updated_at"`
	DeletedAt               *time.Time     `bun:"deleted_at,soft_delete"           json:"deleted_at,omitempty"`

	Attachments []*Attachment `bun:"rel:has-many,join:id=message_id" json:"attachments,omitempty"`
}

const (
	MsgTypeIncoming = 0
	MsgTypeOutgoing = 1
	MsgTypeActivity = 2
	MsgTypeTemplate = 3
)

type Attachment struct {
	bun.BaseModel   `bun:"table:attachments"`
	ID              int64          `bun:"id,pk,autoincrement" json:"id"`
	MessageID       int64          `bun:"message_id,notnull"   json:"message_id"`
	AccountID       int64          `bun:"account_id,notnull"   json:"account_id"`
	FileType        int            `bun:"file_type,default:0"  json:"file_type"`
	ExternalURL     *string        `bun:"external_url"         json:"external_url"`
	CoordinatesLat  float64        `bun:"coordinates_lat"      json:"coordinates_lat"`
	CoordinatesLong float64        `bun:"coordinates_long"     json:"coordinates_long"`
	FallbackTitle   *string        `bun:"fallback_title"       json:"fallback_title"`
	Extension       *string        `bun:"extension"           json:"extension"`
	Meta            map[string]any `bun:"meta,type:jsonb"     json:"meta"`
	CreatedAt       time.Time      `bun:"created_at,notnull"   json:"created_at"`
	UpdatedAt       time.Time      `bun:"updated_at,notnull"   json:"updated_at"`
}
