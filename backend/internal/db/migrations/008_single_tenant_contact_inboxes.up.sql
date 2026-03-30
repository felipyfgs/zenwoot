-- =====================================================
-- Make accountId nullable for single-tenant mode
-- =====================================================

-- Drop NOT NULL constraints on accountId columns
ALTER TABLE "wzChannelsWhatsapp" ALTER COLUMN "accountId" DROP NOT NULL;
ALTER TABLE "wzInboxes" ALTER COLUMN "accountId" DROP NOT NULL;
ALTER TABLE "wzContacts" ALTER COLUMN "accountId" DROP NOT NULL;
ALTER TABLE "wzConversations" ALTER COLUMN "accountId" DROP NOT NULL;
ALTER TABLE "wzMessages" ALTER COLUMN "accountId" DROP NOT NULL;

-- Make FK constraints DEFERRABLE so we can set accountId after row creation
ALTER TABLE "wzChannelsWhatsapp" DROP CONSTRAINT IF EXISTS "wzChannelsWhatsapp_accountId_fkey";
ALTER TABLE "wzChannelsWhatsapp" ADD CONSTRAINT "wzChannelsWhatsapp_accountId_fkey"
    FOREIGN KEY ("accountId") REFERENCES "wzAccounts"("id") ON DELETE SET NULL DEFERRABLE;

ALTER TABLE "wzInboxes" DROP CONSTRAINT IF EXISTS "wzInboxes_accountId_fkey";
ALTER TABLE "wzInboxes" ADD CONSTRAINT "wzInboxes_accountId_fkey"
    FOREIGN KEY ("accountId") REFERENCES "wzAccounts"("id") ON DELETE SET NULL DEFERRABLE;

ALTER TABLE "wzContacts" DROP CONSTRAINT IF EXISTS "wzContacts_accountId_fkey";
ALTER TABLE "wzContacts" ADD CONSTRAINT "wzContacts_accountId_fkey"
    FOREIGN KEY ("accountId") REFERENCES "wzAccounts"("id") ON DELETE SET NULL DEFERRABLE;

ALTER TABLE "wzConversations" DROP CONSTRAINT IF EXISTS "wzConversations_accountId_fkey";
ALTER TABLE "wzConversations" ADD CONSTRAINT "wzConversations_accountId_fkey"
    FOREIGN KEY ("accountId") REFERENCES "wzAccounts"("id") ON DELETE SET NULL DEFERRABLE;

ALTER TABLE "wzMessages" DROP CONSTRAINT IF EXISTS "wzMessages_accountId_fkey";
ALTER TABLE "wzMessages" ADD CONSTRAINT "wzMessages_accountId_fkey"
    FOREIGN KEY ("accountId") REFERENCES "wzAccounts"("id") ON DELETE SET NULL DEFERRABLE;

-- =====================================================
-- Contact Inboxes Table (links contacts to inboxes)
-- =====================================================
CREATE TABLE IF NOT EXISTS "wzContactInboxes" (
    "id" VARCHAR(100) PRIMARY KEY,
    "contactId" VARCHAR(100) NOT NULL REFERENCES "wzContacts"("id") ON DELETE CASCADE,
    "inboxId" VARCHAR(100) NOT NULL REFERENCES "wzInboxes"("id") ON DELETE CASCADE,
    "sourceId" VARCHAR(255) NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idxWzContactInboxesContact" ON "wzContactInboxes" ("contactId");
CREATE INDEX IF NOT EXISTS "idxWzContactInboxesInbox" ON "wzContactInboxes" ("inboxId");
CREATE UNIQUE INDEX IF NOT EXISTS "idxWzContactInboxesInboxSource" ON "wzContactInboxes" ("inboxId", "sourceId");

DROP TRIGGER IF EXISTS "updateWzContactInboxesUpdatedAt" ON "wzContactInboxes";
CREATE TRIGGER "updateWzContactInboxesUpdatedAt"
    BEFORE UPDATE ON "wzContactInboxes"
    FOR EACH ROW
    EXECUTE FUNCTION "updateUpdatedAtColumn"();

COMMENT ON TABLE "wzContactInboxes" IS 'Links contacts to specific inboxes with channel-specific sourceId';
COMMENT ON COLUMN "wzContactInboxes"."sourceId" IS 'Channel-specific identifier (e.g. WhatsApp JID, phone number)';

-- =====================================================
-- Default account auto-creation trigger
-- =====================================================
CREATE OR REPLACE FUNCTION ensure_default_account()
RETURNS VOID AS $$
BEGIN
    INSERT INTO "wzAccounts" ("id", "name", "apiKey", "settings", "createdAt", "updatedAt")
    VALUES ('default', 'Default Account', 'sk_default', '{}', NOW(), NOW())
    ON CONFLICT ("id") DO NOTHING;
END;
$$ language 'plpgsql';

-- Run once on migration
SELECT ensure_default_account();

COMMENT ON FUNCTION ensure_default_account() IS 'Ensures a default account exists for single-tenant mode';
