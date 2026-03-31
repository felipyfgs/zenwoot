-- Companies
CREATE TABLE IF NOT EXISTS "companies" (
    "id"               BIGSERIAL PRIMARY KEY,
    "account_id"       BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "name"             VARCHAR(255),
    "description"      TEXT,
    "logo"             VARCHAR(500),
    "website"          VARCHAR(500),
    "identifier"       VARCHAR(255),
    "custom_attributes" JSONB DEFAULT '{}',
    "created_at"       TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"       TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_companies_account ON "companies"("account_id");
CREATE INDEX IF NOT EXISTS idx_companies_name ON "companies"("name");

-- Notes
CREATE TABLE IF NOT EXISTS "notes" (
    "id"          BIGSERIAL PRIMARY KEY,
    "account_id"  BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "contact_id"  BIGINT NOT NULL REFERENCES "contacts"("id") ON DELETE CASCADE,
    "user_id"     BIGINT NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "content"     TEXT NOT NULL,
    "created_at"  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"  TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_notes_account ON "notes"("account_id");
CREATE INDEX IF NOT EXISTS idx_notes_contact ON "notes"("contact_id");

-- Campaigns
CREATE TABLE IF NOT EXISTS "campaigns" (
    "id"              BIGSERIAL PRIMARY KEY,
    "account_id"      BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "inbox_id"        BIGINT NOT NULL REFERENCES "inboxes"("id") ON DELETE CASCADE,
    "name"            VARCHAR(255),
    "description"     TEXT,
    "campaign_type"   INTEGER DEFAULT 0,
    "campaign_status" INTEGER DEFAULT 0,
    "scheduled_at"    TIMESTAMPTZ,
    "sent_at"         TIMESTAMPTZ,
    "message"         TEXT,
    "trigger_rules"   JSONB DEFAULT '{}',
    "created_at"      TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"      TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_campaigns_account ON "campaigns"("account_id");
CREATE INDEX IF NOT EXISTS idx_campaigns_inbox ON "campaigns"("inbox_id");

-- Agent Bots
CREATE TABLE IF NOT EXISTS "agent_bots" (
    "id"            BIGSERIAL PRIMARY KEY,
    "account_id"    BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "name"          VARCHAR(255),
    "description"   TEXT,
    "avatar_url"    VARCHAR(500),
    "bot_type"      INTEGER DEFAULT 0,
    "bot_config"    JSONB DEFAULT '{}',
    "oauth_token"   TEXT,
    "user_id"       BIGINT REFERENCES "users"("id") ON DELETE SET NULL,
    "access_token"  TEXT,
    "created_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_agent_bots_account ON "agent_bots"("account_id");

-- Macros
CREATE TABLE IF NOT EXISTS "macros" (
    "id"            BIGSERIAL PRIMARY KEY,
    "account_id"    BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "name"          VARCHAR(255),
    "actions"       JSONB DEFAULT '{}',
    "action_types"  JSONB DEFAULT '[]',
    "active"        BOOLEAN DEFAULT true,
    "created_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_macros_account ON "macros"("account_id");
