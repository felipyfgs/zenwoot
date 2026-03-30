-- =====================================================
-- Contacts Table (channel-agnostic)
-- =====================================================
CREATE TABLE IF NOT EXISTS "wzContacts" (
    "id" VARCHAR(100) PRIMARY KEY,
    "accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
    "identifier" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) DEFAULT '',
    "pushName" VARCHAR(255) DEFAULT '',
    "avatarUrl" VARCHAR(2048) DEFAULT '',
    "isBlocked" BOOLEAN NOT NULL DEFAULT false,
    "metadata" JSONB NOT NULL DEFAULT '{}',
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idxWzContactsAccount" ON "wzContacts" ("accountId");
CREATE UNIQUE INDEX IF NOT EXISTS "idxWzContactsAccountIdentifier" ON "wzContacts" ("accountId", "identifier");

DROP TRIGGER IF EXISTS "updateWzContactsUpdatedAt" ON "wzContacts";
CREATE TRIGGER "updateWzContactsUpdatedAt"
    BEFORE UPDATE ON "wzContacts"
    FOR EACH ROW
    EXECUTE FUNCTION "updateUpdatedAtColumn"();

COMMENT ON TABLE "wzContacts" IS 'Contacts across all channels';
COMMENT ON COLUMN "wzContacts"."identifier" IS 'Channel-specific identifier (JID, phone, email, etc.)';
