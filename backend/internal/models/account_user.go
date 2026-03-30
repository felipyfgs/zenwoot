package models

import (
	"time"

	"github.com/uptrace/bun"
)

type AccountUser struct {
	bun.BaseModel `bun:"table:account_users"`
	ID            int64      `bun:"id,pk,autoincrement"      json:"id"`
	AccountID     int64      `bun:"account_id,notnull"        json:"account_id"`
	UserID        int64      `bun:"user_id,notnull"           json:"user_id"`
	Role          int        `bun:"role,default:0"           json:"role"`
	Availability  int        `bun:"availability,default:0"   json:"availability"`
	AutoOffline   bool       `bun:"auto_offline,default:true" json:"auto_offline"`
	ActiveAt      *time.Time `bun:"active_at"                 json:"active_at"`
	CreatedAt     time.Time  `bun:"created_at,notnull"        json:"created_at"`
	UpdatedAt     time.Time  `bun:"updated_at,notnull"        json:"updated_at"`

	User    *User    `bun:"rel:belongs-to,join:user_id=id"    json:"user,omitempty"`
	Account *Account `bun:"rel:belongs-to,join:account_id=id" json:"account,omitempty"`
}
