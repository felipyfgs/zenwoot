package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Note struct {
	bun.BaseModel `bun:"table:notes"`
	ID            int64     `bun:"id,pk,autoincrement"  json:"id"`
	AccountID     int64     `bun:"account_id,notnull"    json:"account_id"`
	ContactID     int64     `bun:"contact_id,notnull"    json:"contact_id"`
	UserID        int64     `bun:"user_id,notnull"       json:"user_id"`
	Content       string    `bun:"content,notnull"       json:"content"`
	CreatedAt     time.Time `bun:"created_at,notnull"    json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"    json:"updated_at"`
}
