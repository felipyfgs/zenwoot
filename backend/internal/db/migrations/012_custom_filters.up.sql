CREATE TABLE IF NOT EXISTS "custom_filters" (
    "id"            BIGSERIAL PRIMARY KEY,
    "account_id"    BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "user_id"       BIGINT NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "name"          VARCHAR(255) NOT NULL,
    "filter_type"   VARCHAR(50) NOT NULL,
    "query"         JSONB DEFAULT '{}' NOT NULL,
    "created_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_custom_filters_account ON "custom_filters"("account_id");
CREATE INDEX IF NOT EXISTS idx_custom_filters_user ON "custom_filters"("user_id");
