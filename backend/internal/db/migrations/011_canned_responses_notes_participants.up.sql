-- Canned Responses table
CREATE TABLE IF NOT EXISTS "wzCannedResponses" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
"shortCode" VARCHAR(255) NOT NULL,
"content" TEXT NOT NULL,
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
"updatedAt" TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE("accountId", "shortCode")
);

CREATE INDEX IF NOT EXISTS idx_canned_responses_account ON "wzCannedResponses"("accountId");
CREATE INDEX IF NOT EXISTS idx_canned_responses_short_code ON "wzCannedResponses"("shortCode");

-- Contact Notes table
CREATE TABLE IF NOT EXISTS "wzNotes" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
"contactId" VARCHAR(100) NOT NULL REFERENCES "wzContacts"("id") ON DELETE CASCADE,
"userId" UUID REFERENCES "wzUsers"("id") ON DELETE SET NULL,
"content" TEXT NOT NULL,
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
"updatedAt" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notes_contact ON "wzNotes"("contactId");
CREATE INDEX IF NOT EXISTS idx_notes_account ON "wzNotes"("accountId");

-- Conversation Participants table
CREATE TABLE IF NOT EXISTS "wzConversationParticipants" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
"conversationId" VARCHAR(100) NOT NULL REFERENCES "wzConversations"("id") ON DELETE CASCADE,
"userId" UUID NOT NULL REFERENCES "wzUsers"("id") ON DELETE CASCADE,
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
"updatedAt" TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE("conversationId", "userId")
);

CREATE INDEX IF NOT EXISTS idx_conv_participants_conv ON "wzConversationParticipants"("conversationId");
CREATE INDEX IF NOT EXISTS idx_conv_participants_user ON "wzConversationParticipants"("userId");
