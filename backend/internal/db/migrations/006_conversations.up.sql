-- =====================================================
-- Conversations Table
-- =====================================================
CREATE TABLE IF NOT EXISTS "wzConversations" (
    "id" VARCHAR(100) PRIMARY KEY,
    "accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
    "inboxId" VARCHAR(100) NOT NULL REFERENCES "wzInboxes"("id") ON DELETE CASCADE,
    "contactId" VARCHAR(100) REFERENCES "wzContacts"("id") ON DELETE SET NULL,
    "identifier" VARCHAR(255) NOT NULL,
    "lastMessage" TEXT DEFAULT '',
    "lastMessageAt" TIMESTAMPTZ,
    "unreadCount" INTEGER NOT NULL DEFAULT 0,
    "status" VARCHAR(50) NOT NULL DEFAULT 'open',
    "metadata" JSONB NOT NULL DEFAULT '{}',
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idxWzConversationsAccount" ON "wzConversations" ("accountId");
CREATE INDEX IF NOT EXISTS "idxWzConversationsInbox" ON "wzConversations" ("inboxId");
CREATE INDEX IF NOT EXISTS "idxWzConversationsContact" ON "wzConversations" ("contactId");
CREATE INDEX IF NOT EXISTS "idxWzConversationsStatus" ON "wzConversations" ("status");
CREATE UNIQUE INDEX IF NOT EXISTS "idxWzConversationsInboxIdentifier" ON "wzConversations" ("inboxId", "identifier");

DROP TRIGGER IF EXISTS "updateWzConversationsUpdatedAt" ON "wzConversations";
CREATE TRIGGER "updateWzConversationsUpdatedAt"
    BEFORE UPDATE ON "wzConversations"
    FOR EACH ROW
    EXECUTE FUNCTION "updateUpdatedAtColumn"();

COMMENT ON TABLE "wzConversations" IS 'Chat conversations per inbox';
COMMENT ON COLUMN "wzConversations"."identifier" IS 'Channel-specific chat identifier (chat JID, thread ID, etc.)';
