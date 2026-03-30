CREATE TABLE IF NOT EXISTS "users" (
    "id"                 BIGSERIAL PRIMARY KEY,
    "name"               VARCHAR(255) NOT NULL,
    "email"              VARCHAR(255) NOT NULL UNIQUE,
    "password_hash"       VARCHAR(255) NOT NULL,
    "display_name"        VARCHAR(255),
    "avatar_url"          TEXT,
    "locale"             VARCHAR(10)  DEFAULT 'en',
    "type"               VARCHAR(50)  DEFAULT 'User',
    "availability_status" INTEGER      DEFAULT 0,
    "created_at"          TIMESTAMPTZ  DEFAULT NOW() NOT NULL,
    "updated_at"          TIMESTAMPTZ  DEFAULT NOW() NOT NULL,
    "deleted_at"          TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_users_email ON "users"("email");
CREATE TRIGGER set_updated_at BEFORE UPDATE ON "users"
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS "access_tokens" (
    "id"        BIGSERIAL PRIMARY KEY,
    "owner_type" VARCHAR(50),
    "owner_id"   BIGINT,
    "token"     VARCHAR(255) UNIQUE NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_access_tokens_owner ON "access_tokens"("owner_type", "owner_id");
CREATE INDEX IF NOT EXISTS idx_access_tokens_token ON "access_tokens"("token");
