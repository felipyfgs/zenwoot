CREATE TABLE IF NOT EXISTS "webhooks" (
    "id"            BIGSERIAL PRIMARY KEY,
    "account_id"     BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "inbox_id"       BIGINT REFERENCES "inboxes"("id") ON DELETE CASCADE,
    "url"           TEXT NOT NULL,
    "subscriptions" TEXT[] NOT NULL DEFAULT '{}',
    "hmac_token"     VARCHAR(255),
    "created_at"     TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"     TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "deleted_at"     TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_webhooks_account ON "webhooks"("account_id");

CREATE TABLE IF NOT EXISTS "automation_rules" (
    "id"          BIGSERIAL PRIMARY KEY,
    "account_id"   BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "name"        VARCHAR(255) NOT NULL,
    "description" TEXT,
    "event_name"   VARCHAR(255) NOT NULL,
    "conditions"  JSONB DEFAULT '{}' NOT NULL,
    "actions"     JSONB DEFAULT '{}' NOT NULL,
    "active"      BOOLEAN DEFAULT true NOT NULL,
    "created_at"   TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"   TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "deleted_at"   TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_automation_rules_account ON "automation_rules"("account_id");
CREATE INDEX IF NOT EXISTS idx_automation_rules_event ON "automation_rules"("event_name");
CREATE INDEX IF NOT EXISTS idx_automation_rules_active ON "automation_rules"("active");
CREATE TRIGGER set_updated_at BEFORE UPDATE ON "automation_rules"
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
