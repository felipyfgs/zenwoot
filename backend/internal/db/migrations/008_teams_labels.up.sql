CREATE TABLE IF NOT EXISTS "labels" (
    "id"            BIGSERIAL PRIMARY KEY,
    "account_id"     BIGINT REFERENCES "accounts"("id") ON DELETE CASCADE,
    "title"         VARCHAR(255),
    "description"   TEXT,
    "color"         VARCHAR(7) DEFAULT '#1f93ff' NOT NULL,
    "show_on_sidebar" BOOLEAN,
    "created_at"     TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"     TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("title", "account_id")
);
CREATE INDEX IF NOT EXISTS idx_labels_account ON "labels"("account_id");

CREATE TABLE IF NOT EXISTS "conversation_labels" (
    "conversation_id" BIGINT NOT NULL REFERENCES "conversations"("id") ON DELETE CASCADE,
    "label_id"        BIGINT NOT NULL REFERENCES "labels"("id")        ON DELETE CASCADE,
    PRIMARY KEY ("conversation_id", "label_id")
);
CREATE INDEX IF NOT EXISTS idx_conv_labels_label ON "conversation_labels"("label_id");

CREATE TABLE IF NOT EXISTS "teams" (
    "id"              BIGSERIAL PRIMARY KEY,
    "account_id"       BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "name"            VARCHAR(255) NOT NULL,
    "description"     TEXT,
    "allow_auto_assign" BOOLEAN DEFAULT true,
    "created_at"       TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"       TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "deleted_at"       TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_teams_account ON "teams"("account_id");
CREATE TRIGGER set_updated_at BEFORE UPDATE ON "teams"
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS "team_members" (
    "id"        BIGSERIAL PRIMARY KEY,
    "team_id"    BIGINT NOT NULL REFERENCES "teams"("id") ON DELETE CASCADE,
    "user_id"    BIGINT NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("team_id", "user_id")
);
CREATE INDEX IF NOT EXISTS idx_team_members_team ON "team_members"("team_id");
