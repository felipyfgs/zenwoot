CREATE TABLE IF NOT EXISTS "notifications" (
    "id"            BIGSERIAL PRIMARY KEY,
    "user_id"        BIGINT NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "account_id"     BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "notification_type" VARCHAR(50) NOT NULL,
    "title"         VARCHAR(255) NOT NULL,
    "body"          TEXT,
    "primary_color" VARCHAR(10),
    "action_url"    TEXT,
    "read_at"       TIMESTAMPTZ,
    "created_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON "notifications"("user_id");
CREATE INDEX IF NOT EXISTS idx_notifications_account ON "notifications"("account_id");

CREATE TABLE IF NOT EXISTS "notification_settings" (
    "id"           BIGSERIAL PRIMARY KEY,
    "user_id"      BIGINT NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "account_id"   BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "email_enabled" BOOLEAN DEFAULT true,
    "push_enabled"  BOOLEAN DEFAULT true,
    "created_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("user_id", "account_id")
);

CREATE TABLE IF NOT EXISTS "notification_subscriptions" (
    "id"            BIGSERIAL PRIMARY KEY,
    "user_id"       BIGINT NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "account_id"    BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "subscription_type" VARCHAR(50) NOT NULL,
    "subscription_attributes" JSONB DEFAULT '{}',
    "created_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("user_id", "account_id", "subscription_type")
);
