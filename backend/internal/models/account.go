package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel       `bun:"table:accounts"`
	ID                  int64          `bun:"id,pk,autoincrement"            json:"id"`
	Name                string         `bun:"name,notnull"                   json:"name"`
	Domain              *string        `bun:"domain"                         json:"domain"`
	SupportEmail        *string        `bun:"support_email"                   json:"support_email"`
	Locale              int            `bun:"locale,default:0"               json:"locale"`
	FeatureFlags        int64          `bun:"feature_flags,notnull,default:0" json:"feature_flags"`
	AutoResolveDuration *int           `bun:"auto_resolve_duration"            json:"auto_resolve_duration"`
	Limits              map[string]any `bun:"limits,type:jsonb"              json:"limits"`
	CustomAttributes    map[string]any `bun:"custom_attributes,type:jsonb"    json:"custom_attributes"`
	Settings            map[string]any `bun:"settings,type:jsonb"            json:"settings"`
	InternalAttributes  map[string]any `bun:"internal_attributes,type:jsonb"    json:"internal_attributes"`
	Status              int            `bun:"status,default:0"               json:"status"`
	CreatedAt           time.Time      `bun:"created_at,notnull"              json:"created_at"`
	UpdatedAt           time.Time      `bun:"updated_at,notnull"              json:"updated_at"`
}
