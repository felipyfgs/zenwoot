package models

import (
	"time"

	"github.com/uptrace/bun"
)

type WorkingHours struct {
	bun.BaseModel `bun:"table:working_hours"`
	ID            int64     `bun:"id,pk,autoincrement"      json:"id"`
	InboxID       int64     `bun:"inbox_id,notnull"          json:"inbox_id"`
	DayOfWeek     int       `bun:"day_of_week,notnull"       json:"day_of_week"`
	ClosedAllDay  bool      `bun:"closed_all_day,default:false" json:"closed_all_day"`
	OpenHour      *int      `bun:"open_hour"                 json:"open_hour"`
	OpenMinutes   int       `bun:"open_minutes,default:0"    json:"open_minutes"`
	CloseHour     *int      `bun:"close_hour"               json:"close_hour"`
	CloseMinutes  int       `bun:"close_minutes,default:0"   json:"close_minutes"`
	OpenAllDay    bool      `bun:"open_all_day,default:false" json:"open_all_day"`
	CreatedAt     time.Time `bun:"created_at,notnull"        json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"        json:"updated_at"`
}

type CustomAttributeDefinition struct {
	bun.BaseModel        `bun:"table:custom_attribute_definitions"`
	ID                   int64     `bun:"id,pk,autoincrement"       json:"id"`
	AccountID            *int64    `bun:"account_id"                json:"account_id"`
	AttributeDisplayName *string   `bun:"attribute_display_name"    json:"attribute_display_name"`
	AttributeKey         *string   `bun:"attribute_key"             json:"attribute_key"`
	AttributeDisplayType int       `bun:"attribute_display_type,default:0" json:"attribute_display_type"`
	AttributeModel       int       `bun:"attribute_model,default:0" json:"attribute_model"`
	AttributeValues      []string  `bun:"attribute_values,type:jsonb" json:"attribute_values"`
	AttributeDescription *string   `bun:"attribute_description"    json:"attribute_description"`
	RegexPattern         *string   `bun:"regex_pattern"            json:"regex_pattern"`
	RegexCue             *string   `bun:"regex_cue"                json:"regex_cue"`
	CreatedAt            time.Time `bun:"created_at,notnull"        json:"created_at"`
	UpdatedAt            time.Time `bun:"updated_at,notnull"        json:"updated_at"`
}
