CREATE TABLE IF NOT EXISTS "channel_web_widgets" (
    "id"                 BIGSERIAL PRIMARY KEY,
    "account_id"          BIGINT      NOT NULL,
    "website_url"         TEXT,
    "website_token"       VARCHAR(255) UNIQUE,
    "widget_color"        VARCHAR(7)  DEFAULT '#1f93ff',
    "welcome_title"       TEXT,
    "welcome_tagline"     TEXT,
    "hmac_token"          VARCHAR(255) UNIQUE,
    "hmac_mandatory"      BOOLEAN     DEFAULT false,
    "pre_chat_form_enabled" BOOLEAN     DEFAULT false,
    "pre_chat_form_options" JSONB       DEFAULT '{}',
    "feature_flags"       INTEGER     DEFAULT 7,
    "reply_time"          INTEGER     DEFAULT 0,
    "allowed_domains"     TEXT        DEFAULT '',
    "continuity_via_email" BOOLEAN     DEFAULT true NOT NULL,
    "created_at"          TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"          TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_channel_web_widgets_token ON "channel_web_widgets"("website_token");

CREATE TABLE IF NOT EXISTS "channel_emails" (
    "id"                     BIGSERIAL PRIMARY KEY,
    "account_id"              BIGINT      NOT NULL,
    "email"                  VARCHAR(255) NOT NULL UNIQUE,
    "forward_to_email"         VARCHAR(255) NOT NULL UNIQUE,
    "imap_enabled"            BOOLEAN     DEFAULT false,
    "imap_address"            VARCHAR(255) DEFAULT '',
    "imap_port"               INTEGER     DEFAULT 0,
    "imap_login"              VARCHAR(255) DEFAULT '',
    "imap_password"           VARCHAR(255) DEFAULT '',
    "imap_enable_ssl"          BOOLEAN     DEFAULT true,
    "smtp_enabled"            BOOLEAN     DEFAULT false,
    "smtp_address"            VARCHAR(255) DEFAULT '',
    "smtp_port"               INTEGER     DEFAULT 0,
    "smtp_login"              VARCHAR(255) DEFAULT '',
    "smtp_password"           VARCHAR(255) DEFAULT '',
    "smtp_domain"             VARCHAR(255) DEFAULT '',
    "smtp_enable_starttls_auto" BOOLEAN     DEFAULT true,
    "smtp_enable_ssl_tls"       BOOLEAN     DEFAULT false,
    "provider"               VARCHAR(50),
    "provider_config"         JSONB       DEFAULT '{}',
    "verified_for_sending"     BOOLEAN     DEFAULT false NOT NULL,
    "created_at"              TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"              TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS "channel_apis" (
    "id"                   BIGSERIAL PRIMARY KEY,
    "account_id"            BIGINT      NOT NULL,
    "webhook_url"           TEXT,
    "identifier"           VARCHAR(255) UNIQUE,
    "hmac_token"            VARCHAR(255) UNIQUE,
    "hmac_mandatory"        BOOLEAN     DEFAULT false,
    "additional_attributes" JSONB       DEFAULT '{}',
    "created_at"            TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"            TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS "channel_telegrams" (
    "id"        BIGSERIAL PRIMARY KEY,
    "account_id" BIGINT      NOT NULL,
    "bot_name"   VARCHAR(255),
    "bot_token"  VARCHAR(255) NOT NULL UNIQUE,
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS "channel_whatsapps" (
    "id"                          BIGSERIAL PRIMARY KEY,
    "account_id"                   BIGINT      NOT NULL,
    "phone_number"                 VARCHAR(50) NOT NULL UNIQUE,
    "provider"                    VARCHAR(50) DEFAULT 'default',
    "provider_config"              JSONB       DEFAULT '{}',
    "message_templates"            JSONB       DEFAULT '{}',
    "message_templates_last_updated" TIMESTAMPTZ,
    "created_at"                   TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"                   TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS "inboxes" (
    "id"                         BIGSERIAL PRIMARY KEY,
    "account_id"                  BIGINT      NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "channel_id"                  BIGINT      NOT NULL,
    "channel_type"                VARCHAR(100) NOT NULL,
    "name"                       VARCHAR(255) NOT NULL,
    "auto_assign_enabled"         BOOLEAN     DEFAULT true,
    "greeting_enabled"            BOOLEAN     DEFAULT false,
    "greeting_message"            TEXT,
    "email_address"              VARCHAR(255),
    "working_hours_enabled"       BOOLEAN     DEFAULT false,
    "out_of_office_message"       TEXT,
    "timezone"                   VARCHAR(64) DEFAULT 'UTC',
    "enable_email_collect"        BOOLEAN     DEFAULT true,
    "csat_enabled"                BOOLEAN     DEFAULT false,
    "allow_resolved_messages"     BOOLEAN     DEFAULT true,
    "single_conversation_only"    BOOLEAN     DEFAULT false NOT NULL,
    "auto_assign_config"          JSONB       DEFAULT '{}',
    "csat_config"                 JSONB       DEFAULT '{}' NOT NULL,
    "sender_name_type"             INTEGER     DEFAULT 0,
    "business_name"               VARCHAR(255),
    "created_at"                  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"                  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "deleted_at"                  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_inboxes_account ON "inboxes"("account_id");
CREATE INDEX IF NOT EXISTS idx_inboxes_channel ON "inboxes"("channel_id", "channel_type");
CREATE TRIGGER set_updated_at BEFORE UPDATE ON "inboxes"
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS "inbox_members" (
    "id"        BIGSERIAL PRIMARY KEY,
    "inbox_id"   BIGINT NOT NULL REFERENCES "inboxes"("id") ON DELETE CASCADE,
    "user_id"    BIGINT NOT NULL REFERENCES "users"("id")   ON DELETE CASCADE,
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("inbox_id", "user_id")
);
CREATE INDEX IF NOT EXISTS idx_inbox_members_inbox ON "inbox_members"("inbox_id");
