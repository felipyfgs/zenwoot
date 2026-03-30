DROP TRIGGER IF EXISTS "updateWzContactInboxesUpdatedAt" ON "wzContactInboxes";
DROP TABLE IF EXISTS "wzContactInboxes";
DROP FUNCTION IF EXISTS ensure_default_account();

-- Restore NOT NULL constraints (would fail if data has nulls, so this is one-way in practice)
-- ALTER TABLE "wzMessages" ALTER COLUMN "accountId" SET NOT NULL;
-- ALTER TABLE "wzConversations" ALTER COLUMN "accountId" SET NOT NULL;
-- ALTER TABLE "wzContacts" ALTER COLUMN "accountId" SET NOT NULL;
-- ALTER TABLE "wzInboxes" ALTER COLUMN "accountId" SET NOT NULL;
-- ALTER TABLE "wzChannelsWhatsapp" ALTER COLUMN "accountId" SET NOT NULL;
