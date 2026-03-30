-- =====================================================
-- Accounts Table (multi-tenant root entity)
-- =====================================================
CREATE TABLE IF NOT EXISTS "wzAccounts" (
    "id" VARCHAR(100) PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "domain" VARCHAR(100),
    "apiKey" VARCHAR(255) NOT NULL UNIQUE,
    "settings" JSONB NOT NULL DEFAULT '{}',
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "idxWzAccountsDomain"
    ON "wzAccounts" ("domain")
    WHERE "domain" IS NOT NULL AND "domain" != '';

CREATE UNIQUE INDEX IF NOT EXISTS "idxWzAccountsApiKey" ON "wzAccounts" ("apiKey");

DROP TRIGGER IF EXISTS "updateWzAccountsUpdatedAt" ON "wzAccounts";
CREATE TRIGGER "updateWzAccountsUpdatedAt"
    BEFORE UPDATE ON "wzAccounts"
    FOR EACH ROW
    EXECUTE FUNCTION "updateUpdatedAtColumn"();

COMMENT ON TABLE "wzAccounts" IS 'Multi-tenant accounts';
