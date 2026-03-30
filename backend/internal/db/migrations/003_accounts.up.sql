CREATE TABLE IF NOT EXISTS "accounts" (
    "id"                   BIGSERIAL PRIMARY KEY,
    "name"                 VARCHAR(255) NOT NULL,
    "domain"               VARCHAR(100) UNIQUE,
    "support_email"         VARCHAR(100),
    "locale"               INTEGER      DEFAULT 0,
    "feature_flags"         BIGINT       DEFAULT 0 NOT NULL,
    "auto_resolve_duration"  INTEGER,
    "limits"               JSONB        DEFAULT '{}',
    "custom_attributes"     JSONB        DEFAULT '{}',
    "settings"             JSONB        DEFAULT '{}',
    "internal_meta"         JSONB        DEFAULT '{}' NOT NULL,
    "status"               INTEGER      DEFAULT 0,
    "created_at"            TIMESTAMPTZ  DEFAULT NOW() NOT NULL,
    "updated_at"            TIMESTAMPTZ  DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_accounts_status ON "accounts"("status");
CREATE TRIGGER set_updated_at BEFORE UPDATE ON "accounts"
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS "account_users" (
    "id"           BIGSERIAL PRIMARY KEY,
    "account_id"    BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "user_id"       BIGINT NOT NULL REFERENCES "users"("id")    ON DELETE CASCADE,
    "role"         INTEGER     DEFAULT 0,
    "availability" INTEGER     DEFAULT 0,
    "auto_offline"  BOOLEAN     DEFAULT true,
    "active_at"     TIMESTAMPTZ,
    "created_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("account_id", "user_id")
);
CREATE INDEX IF NOT EXISTS idx_account_users_account ON "account_users"("account_id");
CREATE INDEX IF NOT EXISTS idx_account_users_user ON "account_users"("user_id");
