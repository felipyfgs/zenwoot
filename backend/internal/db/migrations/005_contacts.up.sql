CREATE TABLE IF NOT EXISTS "companies" (
    "id"            BIGSERIAL PRIMARY KEY,
    "account_id"     BIGINT NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "name"          VARCHAR(255) NOT NULL,
    "domain"        VARCHAR(255),
    "description"   TEXT,
    "contacts_count" INTEGER,
    "created_at"     TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"     TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("account_id", "domain")
);
CREATE INDEX IF NOT EXISTS idx_companies_account ON "companies"("account_id");
CREATE UNIQUE INDEX IF NOT EXISTS idx_companies_domain_partial ON "companies"("account_id", "domain") WHERE "domain" IS NOT NULL;

CREATE TABLE IF NOT EXISTS "contacts" (
    "id"                   BIGSERIAL PRIMARY KEY,
    "account_id"            BIGINT      NOT NULL REFERENCES "accounts"("id") ON DELETE CASCADE,
    "name"                 VARCHAR(255) DEFAULT '',
    "email"                VARCHAR(255),
    "phone_number"          VARCHAR(50),
    "identifier"           VARCHAR(255),
    "avatar_url"            TEXT,
    "extra_attributes"       JSONB       DEFAULT '{}',
    "custom_attributes"      JSONB       DEFAULT '{}',
    "last_activity_at"      TIMESTAMPTZ,
    "contact_type"         INTEGER     DEFAULT 0,
    "first_name"           VARCHAR(255) DEFAULT '',
    "middle_name"          VARCHAR(255) DEFAULT '',
    "last_name"            VARCHAR(255) DEFAULT '',
    "location"             VARCHAR(255) DEFAULT '',
    "country_code"          VARCHAR(10)  DEFAULT '',
    "blocked"              BOOLEAN     DEFAULT false NOT NULL,
    "company_id"            BIGINT REFERENCES "companies"("id") ON DELETE SET NULL,
    "created_at"            TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"            TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "deleted_at"            TIMESTAMPTZ,
    UNIQUE("account_id", "email"),
    UNIQUE("account_id", "identifier")
);
CREATE INDEX IF NOT EXISTS idx_contacts_account ON "contacts"("account_id");
CREATE INDEX IF NOT EXISTS idx_contacts_email ON "contacts"("email", "account_id");
CREATE INDEX IF NOT EXISTS idx_contacts_phone ON "contacts"("phone_number", "account_id");
CREATE INDEX idx_contacts_search ON "contacts" USING gin(
    (COALESCE("name",'') || ' ' || COALESCE("email",'') || ' ' ||
     COALESCE("phone_number",'') || ' ' || COALESCE("identifier",'')) gin_trgm_ops
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON "contacts"
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS "contact_inboxes" (
    "id"           BIGSERIAL PRIMARY KEY,
    "contact_id"    BIGINT REFERENCES "contacts"("id") ON DELETE CASCADE,
    "inbox_id"      BIGINT REFERENCES "inboxes"("id")  ON DELETE CASCADE,
    "source_id"     TEXT        NOT NULL,
    "hmac_verified" BOOLEAN     DEFAULT false,
    "pubsub_token"  VARCHAR(255) UNIQUE,
    "created_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at"    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE("inbox_id", "source_id")
);
CREATE INDEX IF NOT EXISTS idx_contact_inboxes_contact ON "contact_inboxes"("contact_id");
CREATE INDEX IF NOT EXISTS idx_contact_inboxes_source ON "contact_inboxes"("source_id");

CREATE TABLE IF NOT EXISTS "notes" (
    "id"        BIGSERIAL PRIMARY KEY,
    "account_id" BIGINT NOT NULL,
    "contact_id" BIGINT NOT NULL REFERENCES "contacts"("id") ON DELETE CASCADE,
    "user_id"    BIGINT REFERENCES "users"("id") ON DELETE SET NULL,
    "content"   TEXT   NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_notes_contact ON "notes"("contact_id");
