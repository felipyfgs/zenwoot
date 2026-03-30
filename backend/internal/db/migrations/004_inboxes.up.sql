-- =====================================================
-- Inboxes Table (polymorphic channel reference)
-- =====================================================
CREATE TABLE IF NOT EXISTS "wzInboxes" (
    "id" VARCHAR(100) PRIMARY KEY,
    "accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "channelType" VARCHAR(50) NOT NULL,
    "channelId" VARCHAR(100) NOT NULL,
    "status" VARCHAR(50) NOT NULL DEFAULT 'inactive',
    "settings" JSONB NOT NULL DEFAULT '{}',
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idxWzInboxesAccount" ON "wzInboxes" ("accountId");
CREATE INDEX IF NOT EXISTS "idxWzInboxesChannel" ON "wzInboxes" ("channelType", "channelId");
CREATE INDEX IF NOT EXISTS "idxWzInboxesStatus" ON "wzInboxes" ("status");

DROP TRIGGER IF EXISTS "updateWzInboxesUpdatedAt" ON "wzInboxes";
CREATE TRIGGER "updateWzInboxesUpdatedAt"
    BEFORE UPDATE ON "wzInboxes"
    FOR EACH ROW
    EXECUTE FUNCTION "updateUpdatedAtColumn"();

COMMENT ON TABLE "wzInboxes" IS 'Communication inboxes with polymorphic channel reference';
COMMENT ON COLUMN "wzInboxes"."channelType" IS 'Channel type: whatsapp, telegram, signal, email';
COMMENT ON COLUMN "wzInboxes"."channelId" IS 'FK to channel-specific table (e.g. wzChannelsWhatsapp.id)';
