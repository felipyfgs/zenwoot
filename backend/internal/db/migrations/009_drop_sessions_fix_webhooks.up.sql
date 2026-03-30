-- =====================================================
-- Drop legacy wzSessions table
-- =====================================================
DROP TABLE IF EXISTS "wzSessions" CASCADE;

-- =====================================================
-- wzWebhooks: create with final schema or migrate from old schema
-- =====================================================
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'wzWebhooks') THEN
        CREATE TABLE "wzWebhooks" (
            "id"          VARCHAR(100) PRIMARY KEY,
            "inboxId"     VARCHAR(100) NOT NULL REFERENCES "wzInboxes"("id") ON DELETE CASCADE,
            "url"         TEXT NOT NULL,
            "secret"      VARCHAR(255),
            "events"      JSONB NOT NULL DEFAULT '[]',
            "enabled"     BOOLEAN NOT NULL DEFAULT true,
            "natsEnabled" BOOLEAN NOT NULL DEFAULT false,
            "createdAt"   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
            "updatedAt"   TIMESTAMPTZ NOT NULL DEFAULT NOW()
        );
        CREATE INDEX "idxWzWebhooksInboxId" ON "wzWebhooks" ("inboxId");
    ELSE
        -- Migrate old schema: rename sessionId → inboxId if needed
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'wzWebhooks' AND column_name = 'sessionId') THEN
            ALTER TABLE "wzWebhooks" DROP CONSTRAINT IF EXISTS "wzWebhooks_sessionId_fkey";
            ALTER TABLE "wzWebhooks" RENAME COLUMN "sessionId" TO "inboxId";
            ALTER TABLE "wzWebhooks" ADD CONSTRAINT "wzWebhooks_inboxId_fkey"
                FOREIGN KEY ("inboxId") REFERENCES "wzInboxes"("id") ON DELETE CASCADE;
            DROP INDEX IF EXISTS "idxWzWebhooksSessionId";
            CREATE INDEX IF NOT EXISTS "idxWzWebhooksInboxId" ON "wzWebhooks" ("inboxId");
        END IF;
    END IF;
END $$;
