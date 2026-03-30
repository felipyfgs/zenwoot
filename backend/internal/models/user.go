package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel      `bun:"table:users"`
	ID                 int64      `bun:"id,pk,autoincrement"          json:"id"`
	Name               string     `bun:"name,notnull"                 json:"name"`
	Email              string     `bun:"email,notnull"                json:"email"`
	PasswordHash       string     `bun:"password_hash,notnull"         json:"-"`
	DisplayName        *string    `bun:"display_name"                 json:"display_name"`
	AvatarURL          *string    `bun:"avatar_url"                    json:"avatar_url"`
	Locale             string     `bun:"locale,default:'en'"          json:"locale"`
	Type               string     `bun:"type,default:'User'"          json:"type"`
	AvailabilityStatus int        `bun:"availability,default:0"       json:"availability"`
	CreatedAt          time.Time  `bun:"created_at,notnull"            json:"created_at"`
	UpdatedAt          time.Time  `bun:"updated_at,notnull"            json:"updated_at"`
	DeletedAt          *time.Time `bun:"deleted_at,soft_delete"        json:"deleted_at,omitempty"`
}
