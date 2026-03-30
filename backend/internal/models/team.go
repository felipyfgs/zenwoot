package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Team struct {
	bun.BaseModel   `bun:"table:teams"`
	ID              int64      `bun:"id,pk,autoincrement"          json:"id"`
	AccountID       int64      `bun:"account_id,notnull"            json:"account_id"`
	Name            string     `bun:"name,notnull"                 json:"name"`
	Description     *string    `bun:"description"                  json:"description"`
	AllowAutoAssign bool       `bun:"allow_auto_assign,default:true" json:"allow_auto_assign"`
	CreatedAt       time.Time  `bun:"created_at,notnull"            json:"created_at"`
	UpdatedAt       time.Time  `bun:"updated_at,notnull"            json:"updated_at"`
	DeletedAt       *time.Time `bun:"deleted_at,soft_delete"        json:"deleted_at,omitempty"`
}

type TeamMember struct {
	bun.BaseModel `bun:"table:team_members"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	TeamID        int64     `bun:"team_id,notnull"      json:"team_id"`
	UserID        int64     `bun:"user_id,notnull"      json:"user_id"`
	CreatedAt     time.Time `bun:"created_at,notnull"   json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"   json:"updated_at"`

	User *User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}
