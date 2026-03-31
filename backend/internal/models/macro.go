package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Macro struct {
	bun.BaseModel `bun:"table:macros"`
	ID            int64          `bun:"id,pk,autoincrement"       json:"id"`
	AccountID     int64          `bun:"account_id,notnull"         json:"account_id"`
	Name          *string        `bun:"name"                      json:"name"`
	Actions       map[string]any `bun:"actions,type:jsonb"         json:"actions"`
	ActionTypes   []string       `bun:"action_types,type:jsonb"   json:"action_types"`
	Active        bool           `bun:"active,default:true"        json:"active"`
	CreatedAt     time.Time      `bun:"created_at,notnull"         json:"created_at"`
	UpdatedAt     time.Time      `bun:"updated_at,notnull"         json:"updated_at"`
}
