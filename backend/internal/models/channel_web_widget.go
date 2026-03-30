package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ChannelWebWidget struct {
	bun.BaseModel      `bun:"table:channel_web_widgets"`
	ID                 int64          `bun:"id,pk,autoincrement"              json:"id"`
	AccountID          int64          `bun:"account_id,notnull"                json:"account_id"`
	WebsiteURL         *string        `bun:"website_url"                       json:"website_url"`
	WebsiteToken       *string        `bun:"website_token"                     json:"website_token"`
	WidgetColor        string         `bun:"widget_color,default:'#1f93ff'"    json:"widget_color"`
	WelcomeTitle       *string        `bun:"welcome_title"                     json:"welcome_title"`
	WelcomeTagline     *string        `bun:"welcome_tagline"                   json:"welcome_tagline"`
	HmacToken          *string        `bun:"hmac_token"                        json:"-"`
	HmacMandatory      bool           `bun:"hmac_mandatory,default:false"      json:"hmac_mandatory"`
	PreChatFormEnabled bool           `bun:"pre_chat_form_enabled,default:false" json:"pre_chat_form_enabled"`
	PreChatFormOptions map[string]any `bun:"pre_chat_form_options,type:jsonb"    json:"pre_chat_form_options"`
	FeatureFlags       int            `bun:"feature_flags,default:7"           json:"feature_flags"`
	ReplyTime          int            `bun:"reply_time,default:0"              json:"reply_time"`
	AllowedDomains     string         `bun:"allowed_domains,default:''"        json:"allowed_domains"`
	ContinuityViaEmail bool           `bun:"continuity_via_email,notnull"       json:"continuity_via_email"`
	CreatedAt          time.Time      `bun:"created_at,notnull"                json:"created_at"`
	UpdatedAt          time.Time      `bun:"updated_at,notnull"                json:"updated_at"`
}

type ChannelEmail struct {
	bun.BaseModel          `bun:"table:channel_emails"`
	ID                     int64          `bun:"id,pk,autoincrement"                  json:"id"`
	AccountID              int64          `bun:"account_id,notnull"                    json:"account_id"`
	Email                  string         `bun:"email,notnull"                        json:"email"`
	ForwardToEmail         string         `bun:"forward_to_email,notnull"               json:"forward_to_email"`
	ImapEnabled            bool           `bun:"imap_enabled,default:false"            json:"imap_enabled"`
	ImapAddress            string         `bun:"imap_address,default:''"              json:"imap_address"`
	ImapPort               int            `bun:"imap_port,default:0"                   json:"imap_port"`
	ImapLogin              string         `bun:"imap_login,default:''"               json:"imap_login"`
	ImapPassword           string         `bun:"imap_password,default:''"            json:"-"`
	ImapEnableSsl          bool           `bun:"imap_enable_ssl,default:true"          json:"imap_enable_ssl"`
	SmtpEnabled            bool           `bun:"smtp_enabled,default:false"           json:"smtp_enabled"`
	SmtpAddress            string         `bun:"smtp_address,default:''"             json:"smtp_address"`
	SmtpPort               int            `bun:"smtp_port,default:0"                  json:"smtp_port"`
	SmtpLogin              string         `bun:"smtp_login,default:''"              json:"smtp_login"`
	SmtpPassword           string         `bun:"smtp_password,default:''"           json:"-"`
	SmtpDomain             string         `bun:"smtp_domain,default:''"             json:"smtp_domain"`
	SmtpEnableStarttlsAuto bool           `bun:"smtp_enable_starttls_auto,default:true" json:"smtp_enable_starttls_auto"`
	SmtpEnableSslTls       bool           `bun:"smtp_enable_ssl_tls,default:false"      json:"smtp_enable_ssl_tls"`
	Provider               *string        `bun:"provider"                            json:"provider"`
	ProviderConfig         map[string]any `bun:"provider_config,type:jsonb"           json:"provider_config"`
	VerifiedForSending     bool           `bun:"verified_for_sending,notnull"          json:"verified_for_sending"`
	CreatedAt              time.Time      `bun:"created_at,notnull"                   json:"created_at"`
	UpdatedAt              time.Time      `bun:"updated_at,notnull"                   json:"updated_at"`
}

type ChannelApi struct {
	bun.BaseModel        `bun:"table:channel_apis"`
	ID                   int64          `bun:"id,pk,autoincrement"            json:"id"`
	AccountID            int64          `bun:"account_id,notnull"              json:"account_id"`
	WebhookURL           *string        `bun:"webhook_url"                     json:"webhook_url"`
	Identifier           *string        `bun:"identifier"                     json:"identifier"`
	HmacToken            *string        `bun:"hmac_token"                      json:"-"`
	HmacMandatory        bool           `bun:"hmac_mandatory,default:false"    json:"hmac_mandatory"`
	AdditionalAttributes map[string]any `bun:"additional_attributes,type:jsonb" json:"additional_attributes"`
	CreatedAt            time.Time      `bun:"created_at,notnull"              json:"created_at"`
	UpdatedAt            time.Time      `bun:"updated_at,notnull"              json:"updated_at"`
}
