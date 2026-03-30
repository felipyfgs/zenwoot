package models

import (
	"time"

	"github.com/uptrace/bun"
)

type AccessToken struct {
	bun.BaseModel `bun:"table:access_tokens"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	OwnerType     *string   `bun:"owner_type"           json:"owner_type"`
	OwnerID       *int64    `bun:"owner_id"             json:"owner_id"`
	Token         string    `bun:"token,notnull"       json:"token"`
	CreatedAt     time.Time `bun:"created_at,notnull"   json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"   json:"updated_at"`
}
