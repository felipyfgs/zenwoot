-- =====================================================
-- Messages Table
-- =====================================================
CREATE TABLE IF NOT EXISTS "wzMessages" (
    "id" VARCHAR(100) PRIMARY KEY,
    "accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
    "conversationId" VARCHAR(100) NOT NULL REFERENCES "wzConversations"("id") ON DELETE CASCADE,
    "inboxId" VARCHAR(100) NOT NULL REFERENCES "wzInboxes"("id") ON DELETE CASCADE,
    "contactId" VARCHAR(100) REFERENCES "wzContacts"("id") ON DELETE SET NULL,
    "externalId" VARCHAR(255) DEFAULT '',
    "direction" VARCHAR(20) NOT NULL,
    "contentType" VARCHAR(50) NOT NULL DEFAULT 'text',
    "content" TEXT DEFAULT '',
    "mediaUrl" VARCHAR(2048) DEFAULT '',
    "mediaType" VARCHAR(50) DEFAULT '',
    "metadata" JSONB NOT NULL DEFAULT '{}',
    "status" VARCHAR(50) NOT NULL DEFAULT 'sent',
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idxWzMessagesAccount" ON "wzMessages" ("accountId");
CREATE INDEX IF NOT EXISTS "idxWzMessagesConversation" ON "wzMessages" ("conversationId", "createdAt");
CREATE INDEX IF NOT EXISTS "idxWzMessagesInbox" ON "wzMessages" ("inboxId");
CREATE INDEX IF NOT EXISTS "idxWzMessagesExternal" ON "wzMessages" ("externalId")
    WHERE "externalId" IS NOT NULL AND "externalId" != '';
CREATE INDEX IF NOT EXISTS "idxWzMessagesCreatedAt" ON "wzMessages" ("createdAt");

COMMENT ON TABLE "wzMessages" IS 'Chat messages across all channels';
COMMENT ON COLUMN "wzMessages"."direction" IS 'incoming or outgoing';
COMMENT ON COLUMN "wzMessages"."contentType" IS 'text, image, video, audio, document, location, contact, sticker';
COMMENT ON COLUMN "wzMessages"."externalId" IS 'Message ID from the external channel';
