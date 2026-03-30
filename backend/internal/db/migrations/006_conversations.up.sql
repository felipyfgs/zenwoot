CREATE TABLE IF NOT EXISTS "conversations" (
    "id"                   BIGSERIAL PRIMARY KEY,
    "display_id"            INTEGER     NOT NULL,
    "account_id"            BIGINT      NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "inbox_id"              BIGINT      NOT NULL REFERENCES "inboxes"("id")  ON DELETE CASCADE,
    "contact_id"            BIGINT      REFERENCES "contacts"("id") ON DELETE SET NULL,
    "contact_inbox_id"       BIGINT      REFERENCES "contact_inboxes"("id"),
    "assignee_id"           BIGINT      REFERENCES "users"("id") ON DELETE SET NULL,
    "team_id"               BIGINT,
    "campaign_id"           BIGINT,
    "sla_policy_id"          BIGINT,
    "status"               INTEGER     DEFAULT 0 NOT NULL,
    "priority"             INTEGER,
    "uuid"                 UUID        DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    "identifier"           VARCHAR(255),
    "snoozed_until"         TIMESTAMPTZ,
    "waiting_since"         TIMESTAMPTZ,
    "first_reply_at"         TIMESTAMPTZ,
    "last_activity_at"      TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "contact_seen_at"       TIMESTAMPTZ,
    "agent_seen_at"         TIMESTAMPTZ,
    "assignee_seen_at"      TIMESTAMPTZ,
    "extra_attributes"      JSONB       DEFAULT '{}',
    "custom_attributes"     JSONB       DEFAULT '{}',
    "labels_cache"          TEXT,
    "created_at"            TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"            TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("account_id", "display_id")
);
CREATE INDEX IF NOT EXISTS idx_conversations_account ON "conversations"("account_id");
CREATE INDEX IF NOT EXISTS idx_conversations_filter ON "conversations"("account_id", "inbox_id", "status", "assignee_id");
CREATE INDEX IF NOT EXISTS idx_conversations_status ON "conversations"("status", "account_id");
CREATE INDEX IF NOT EXISTS idx_conversations_uuid ON "conversations"("uuid");
CREATE INDEX IF NOT EXISTS idx_conversations_contact ON "conversations"("contact_id");
CREATE INDEX IF NOT EXISTS idx_conversations_team ON "conversations"("team_id");
CREATE INDEX IF NOT EXISTS idx_conversations_waiting ON "conversations"("waiting_since");
CREATE INDEX IF NOT EXISTS idx_conversations_priority ON "conversations"("priority");
CREATE TRIGGER set_conversation_display_id
    BEFORE INSERT ON "conversations"
    FOR EACH ROW EXECUTE FUNCTION generate_conversation_display_id();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON "conversations"
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS "conversation_participants" (
    "id"             BIGSERIAL PRIMARY KEY,
    "account_id"      BIGINT NOT NULL,
    "user_id"         BIGINT NOT NULL REFERENCES "users"("id")          ON DELETE CASCADE,
    "conversation_id" BIGINT NOT NULL REFERENCES "conversations"("id")  ON DELETE CASCADE,
    "created_at"      TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"      TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("user_id", "conversation_id")
);
CREATE INDEX IF NOT EXISTS idx_conv_participants_conv ON "conversation_participants"("conversation_id");

CREATE TABLE IF NOT EXISTS "mentions" (
    "id"             BIGSERIAL PRIMARY KEY,
    "account_id"      BIGINT NOT NULL,
    "user_id"         BIGINT NOT NULL REFERENCES "users"("id")          ON DELETE CASCADE,
    "conversation_id" BIGINT NOT NULL REFERENCES "conversations"("id")  ON DELETE CASCADE,
    "mentioned_at"    TIMESTAMPTZ NOT NULL,
    "created_at"      TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"      TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("user_id", "conversation_id")
);
CREATE INDEX IF NOT EXISTS idx_mentions_conv ON "mentions"("conversation_id");
CREATE INDEX IF NOT EXISTS idx_mentions_user ON "mentions"("user_id");
