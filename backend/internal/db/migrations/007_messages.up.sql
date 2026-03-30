CREATE TABLE IF NOT EXISTS "messages" (
    "id"                      BIGSERIAL PRIMARY KEY,
    "account_id"               BIGINT  NOT NULL REFERENCES "accounts"("id")      ON DELETE CASCADE,
    "inbox_id"                 BIGINT  NOT NULL REFERENCES "inboxes"("id")       ON DELETE CASCADE,
    "conversation_id"          BIGINT  NOT NULL REFERENCES "conversations"("id") ON DELETE CASCADE,
    "sender_type"              VARCHAR(50),
    "sender_id"                BIGINT,
    "content"                 TEXT,
    "content_type"             INTEGER DEFAULT 0 NOT NULL,
    "message_type"             INTEGER NOT NULL,
    "status"                  INTEGER DEFAULT 0,
    "source_id"                TEXT,
    "private"                 BOOLEAN DEFAULT false NOT NULL,
    "content_meta"            JSON    DEFAULT '{}',
    "extra_attributes"        JSONB   DEFAULT '{}',
    "external_ids"            JSONB   DEFAULT '{}',
    "processed_content"       TEXT,
    "sentiment"               JSONB   DEFAULT '{}',
    "created_at"               TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"               TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "deleted_at"               TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_messages_conv ON "messages"("conversation_id");
CREATE INDEX IF NOT EXISTS idx_messages_account ON "messages"("account_id");
CREATE INDEX IF NOT EXISTS idx_messages_sender ON "messages"("sender_type", "sender_id");
CREATE INDEX IF NOT EXISTS idx_messages_source ON "messages"("source_id");
CREATE INDEX IF NOT EXISTS idx_messages_inbox ON "messages"("account_id", "inbox_id");
CREATE INDEX IF NOT EXISTS idx_messages_created ON "messages"("created_at");
CREATE INDEX idx_messages_content_gin ON "messages" USING gin("content" gin_trgm_ops);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON "messages"
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS "attachments" (
    "id"              BIGSERIAL PRIMARY KEY,
    "message_id"       BIGINT  NOT NULL REFERENCES "messages"("id") ON DELETE CASCADE,
    "account_id"       BIGINT  NOT NULL,
    "file_type"        INTEGER DEFAULT 0,
    "external_url"     TEXT,
    "coordinates_lat"  FLOAT   DEFAULT 0.0,
    "coordinates_long" FLOAT   DEFAULT 0.0,
    "fallback_title"   VARCHAR(255),
    "extension"       VARCHAR(50),
    "meta"            JSONB   DEFAULT '{}',
    "created_at"       TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"       TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_attachments_message ON "attachments"("message_id");
CREATE INDEX IF NOT EXISTS idx_attachments_account ON "attachments"("account_id");
