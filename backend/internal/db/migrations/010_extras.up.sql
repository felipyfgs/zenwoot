CREATE TABLE IF NOT EXISTS "canned_responses" (
    "id"        BIGSERIAL PRIMARY KEY,
    "account_id" BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "short_code" VARCHAR(255),
    "content"   TEXT,
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_canned_responses_account ON "canned_responses"("account_id");
CREATE INDEX IF NOT EXISTS idx_canned_responses_short ON "canned_responses"("short_code", "account_id");

CREATE TABLE IF NOT EXISTS "working_hours" (
    "id"           BIGSERIAL PRIMARY KEY,
    "inbox_id"      BIGINT NOT NULL REFERENCES "inboxes"("id") ON DELETE CASCADE,
    "day_of_week"    INTEGER NOT NULL,
    "closed_all_day" BOOLEAN DEFAULT false,
    "open_hour"     INTEGER,
    "open_minutes"  INTEGER DEFAULT 0,
    "close_hour"    INTEGER,
    "close_minutes" INTEGER DEFAULT 0,
    "open_all_day"   BOOLEAN DEFAULT false,
    "created_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_working_hours_inbox ON "working_hours"("inbox_id");

CREATE TABLE IF NOT EXISTS "custom_attribute_definitions" (
    "id"                   BIGSERIAL PRIMARY KEY,
    "account_id"            BIGINT REFERENCES "accounts"("id") ON DELETE CASCADE,
    "attribute_display_name" VARCHAR(255),
    "attribute_key"         VARCHAR(255),
    "attribute_display_type" INTEGER DEFAULT 0,
    "attribute_model"       INTEGER DEFAULT 0,
    "attribute_values"      JSONB DEFAULT '[]',
    "attribute_description" TEXT,
    "regex_pattern"         VARCHAR(255),
    "regex_cue"             VARCHAR(255),
    "created_at"            TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"            TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("attribute_key", "attribute_model", "account_id")
);
CREATE INDEX IF NOT EXISTS idx_custom_attr_account ON "custom_attribute_definitions"("account_id");
