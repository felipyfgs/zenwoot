package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Label struct {
	bun.BaseModel `bun:"table:labels"`
	ID            int64     `bun:"id,pk,autoincrement"            json:"id"`
	AccountID     int64     `bun:"account_id"                      json:"account_id"`
	Title         string    `bun:"title"                          json:"title"`
	Description   *string   `bun:"description"                    json:"description"`
	Color         string    `bun:"color,notnull,default:'#1f93ff'" json:"color"`
	ShowOnSidebar *bool     `bun:"show_on_sidebar"                  json:"show_on_sidebar"`
	CreatedAt     time.Time `bun:"created_at,notnull"              json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"              json:"updated_at"`
}

type ConversationLabel struct {
	bun.BaseModel  `bun:"table:conversation_labels"`
	ConversationID int64 `bun:"conversation_id,pk" json:"conversation_id"`
	LabelID        int64 `bun:"label_id,pk"        json:"label_id"`
}
