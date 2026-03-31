package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Company struct {
	bun.BaseModel    `bun:"table:companies"`
	ID               int64          `bun:"id,pk,autoincrement"        json:"id"`
	AccountID        int64          `bun:"account_id,notnull"          json:"account_id"`
	Name             *string        `bun:"name"                       json:"name"`
	Description      *string        `bun:"description"                json:"description"`
	Logo             *string        `bun:"logo"                       json:"logo"`
	Website          *string        `bun:"website"                    json:"website"`
	Identifier       *string        `bun:"identifier"                 json:"identifier"`
	CustomAttributes map[string]any `bun:"custom_attributes,type:jsonb" json:"custom_attributes"`
	CreatedAt        time.Time      `bun:"created_at,notnull"          json:"created_at"`
	UpdatedAt        time.Time      `bun:"updated_at,notnull"          json:"updated_at"`
}
