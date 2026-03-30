-- This migration is not reversible in practice because wzSessions was dropped
-- and data would be lost. This is provided for schema rollback only.

ALTER TABLE "wzWebhooks" DROP CONSTRAINT IF EXISTS "wzWebhooks_inboxId_fkey";
ALTER TABLE "wzWebhooks" RENAME COLUMN "inboxId" TO "sessionId";

-- Recreate wzSessions (empty)
CREATE TABLE IF NOT EXISTS "wzSessions" (
    "id" VARCHAR(100) PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL UNIQUE,
    "apiKey" VARCHAR(255) NOT NULL UNIQUE,
    "jid" VARCHAR(255) DEFAULT '',
    "qrCode" TEXT DEFAULT '',
    "connected" INTEGER DEFAULT 0,
    "status" VARCHAR(50) NOT NULL DEFAULT 'disconnected',
    "proxy" JSONB NOT NULL DEFAULT '{}',
    "settings" JSONB NOT NULL DEFAULT '{}',
    "metadata" JSONB,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
