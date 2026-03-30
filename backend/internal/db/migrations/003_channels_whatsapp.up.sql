-- =====================================================
-- WhatsApp Channel Table (channel-specific config)
-- =====================================================
CREATE TABLE IF NOT EXISTS "wzChannelsWhatsapp" (
    "id" VARCHAR(100) PRIMARY KEY,
    "accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
    "phoneNumber" VARCHAR(50),
    "jid" VARCHAR(255) DEFAULT '',
    "provider" VARCHAR(50) NOT NULL DEFAULT 'whatsmeow',
    "providerConfig" JSONB NOT NULL DEFAULT '{}',
    "qrCode" TEXT DEFAULT '',
    "connected" BOOLEAN NOT NULL DEFAULT false,
    "connectedAt" TIMESTAMPTZ,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "idxWzChannelsWhatsappAccount" ON "wzChannelsWhatsapp" ("accountId");

CREATE UNIQUE INDEX IF NOT EXISTS "idxWzChannelsWhatsappJid"
    ON "wzChannelsWhatsapp" ("jid")
    WHERE "jid" IS NOT NULL AND "jid" != '';

DROP TRIGGER IF EXISTS "updateWzChannelsWhatsappUpdatedAt" ON "wzChannelsWhatsapp";
CREATE TRIGGER "updateWzChannelsWhatsappUpdatedAt"
    BEFORE UPDATE ON "wzChannelsWhatsapp"
    FOR EACH ROW
    EXECUTE FUNCTION "updateUpdatedAtColumn"();

COMMENT ON TABLE "wzChannelsWhatsapp" IS 'WhatsApp channel configurations';
COMMENT ON COLUMN "wzChannelsWhatsapp"."provider" IS 'WhatsApp provider: whatsmeow, baileys, cloud';
