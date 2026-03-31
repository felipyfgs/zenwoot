package models

import (
	"time"

	"github.com/uptrace/bun"
)

type AgentBot struct {
	bun.BaseModel `bun:"table:agent_bots"`
	ID            int64          `bun:"id,pk,autoincrement"        json:"id"`
	AccountID     int64          `bun:"account_id,notnull"          json:"account_id"`
	Name          *string        `bun:"name"                       json:"name"`
	Description   *string        `bun:"description"                 json:"description"`
	AvatarURL     *string        `bun:"avatar_url"                 json:"avatar_url"`
	BotType       int            `bun:"bot_type,default:0"         json:"bot_type"`
	BotConfig     map[string]any `bun:"bot_config,type:jsonb"      json:"bot_config"`
	OauthToken    *string        `bun:"oauth_token"                json:"-"`
	UserID        *int64         `bun:"user_id"                    json:"user_id"`
	AccessToken   *string        `bun:"access_token"               json:"-"`
	CreatedAt     time.Time      `bun:"created_at,notnull"          json:"created_at"`
	UpdatedAt     time.Time      `bun:"updated_at,notnull"          json:"updated_at"`
}
