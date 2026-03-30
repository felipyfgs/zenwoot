-- Users table
CREATE TABLE IF NOT EXISTS "wzUsers" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"email" VARCHAR(255) NOT NULL UNIQUE,
"name" VARCHAR(255) NOT NULL,
"displayName" VARCHAR(255),
"avatarUrl" TEXT,
"role" VARCHAR(50) NOT NULL DEFAULT 'agent',
"status" VARCHAR(50) NOT NULL DEFAULT 'active',
"passwordHash" VARCHAR(255),
"provider" VARCHAR(50),
"uid" VARCHAR(255),
"settings" JSONB DEFAULT '{}',
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
"updatedAt" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON "wzUsers"("email");
CREATE INDEX IF NOT EXISTS idx_users_role ON "wzUsers"("role");

-- AccountUsers table (links users to accounts)
CREATE TABLE IF NOT EXISTS "wzAccountUsers" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
"userId" UUID NOT NULL REFERENCES "wzUsers"("id") ON DELETE CASCADE,
"role" VARCHAR(50) NOT NULL DEFAULT 'agent',
"active" BOOLEAN NOT NULL DEFAULT true,
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
"updatedAt" TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE("accountId", "userId")
);

CREATE INDEX IF NOT EXISTS idx_account_users_account ON "wzAccountUsers"("accountId");
CREATE INDEX IF NOT EXISTS idx_account_users_user ON "wzAccountUsers"("userId");

-- Teams table
CREATE TABLE IF NOT EXISTS "wzTeams" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
"name" VARCHAR(255) NOT NULL,
"description" TEXT,
"allowAutoAssign" BOOLEAN NOT NULL DEFAULT false,
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
"updatedAt" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_teams_account ON "wzTeams"("accountId");

-- TeamMembers table
CREATE TABLE IF NOT EXISTS "wzTeamMembers" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"teamId" UUID NOT NULL REFERENCES "wzTeams"("id") ON DELETE CASCADE,
"userId" UUID NOT NULL REFERENCES "wzUsers"("id") ON DELETE CASCADE,
"role" VARCHAR(50) NOT NULL DEFAULT 'agent',
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE("teamId", "userId")
);

CREATE INDEX IF NOT EXISTS idx_team_members_team ON "wzTeamMembers"("teamId");
CREATE INDEX IF NOT EXISTS idx_team_members_user ON "wzTeamMembers"("userId");

-- InboxMembers table (agents assigned to inboxes)
CREATE TABLE IF NOT EXISTS "wzInboxMembers" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"inboxId" VARCHAR(100) NOT NULL REFERENCES "wzInboxes"("id") ON DELETE CASCADE,
"userId" UUID NOT NULL REFERENCES "wzUsers"("id") ON DELETE CASCADE,
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE("inboxId", "userId")
);

CREATE INDEX IF NOT EXISTS idx_inbox_members_inbox ON "wzInboxMembers"("inboxId");
CREATE INDEX IF NOT EXISTS idx_inbox_members_user ON "wzInboxMembers"("userId");

-- Labels table
CREATE TABLE IF NOT EXISTS "wzLabels" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"accountId" VARCHAR(100) NOT NULL REFERENCES "wzAccounts"("id") ON DELETE CASCADE,
"title" VARCHAR(255) NOT NULL,
"color" VARCHAR(7) NOT NULL DEFAULT '#1F93FF',
"description" TEXT,
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
"updatedAt" TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE("accountId", "title")
);

CREATE INDEX IF NOT EXISTS idx_labels_account ON "wzLabels"("accountId");

-- ConversationLabels table (many-to-many)
CREATE TABLE IF NOT EXISTS "wzConversationLabels" (
"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
"conversationId" VARCHAR(100) NOT NULL REFERENCES "wzConversations"("id") ON DELETE CASCADE,
"labelId" UUID NOT NULL REFERENCES "wzLabels"("id") ON DELETE CASCADE,
"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
UNIQUE("conversationId", "labelId")
);

CREATE INDEX IF NOT EXISTS idx_conversation_labels_conv ON "wzConversationLabels"("conversationId");
CREATE INDEX IF NOT EXISTS idx_conversation_labels_label ON "wzConversationLabels"("labelId");

-- Add new columns to conversations
ALTER TABLE "wzConversations" ADD COLUMN IF NOT EXISTS "assigneeId" UUID REFERENCES "wzUsers"("id") ON DELETE SET NULL;
ALTER TABLE "wzConversations" ADD COLUMN IF NOT EXISTS "teamId" UUID REFERENCES "wzTeams"("id") ON DELETE SET NULL;
ALTER TABLE "wzConversations" ADD COLUMN IF NOT EXISTS "priority" VARCHAR(50) DEFAULT 'low';
ALTER TABLE "wzConversations" ADD COLUMN IF NOT EXISTS "muted" BOOLEAN DEFAULT false;
ALTER TABLE "wzConversations" ADD COLUMN IF NOT EXISTS "snoozedUntil" TIMESTAMP;
ALTER TABLE "wzConversations" ADD COLUMN IF NOT EXISTS "lastActivityAt" TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_conversations_assignee ON "wzConversations"("assigneeId");
CREATE INDEX IF NOT EXISTS idx_conversations_team ON "wzConversations"("teamId");
CREATE INDEX IF NOT EXISTS idx_conversations_status ON "wzConversations"("status");

-- Add new columns to messages
ALTER TABLE "wzMessages" ADD COLUMN IF NOT EXISTS "senderId" UUID REFERENCES "wzUsers"("id") ON DELETE SET NULL;
ALTER TABLE "wzMessages" ADD COLUMN IF NOT EXISTS "replyToId" VARCHAR(100) REFERENCES "wzMessages"("id") ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_messages_sender ON "wzMessages"("senderId");
CREATE INDEX IF NOT EXISTS idx_messages_reply_to ON "wzMessages"("replyToId");

-- Round-robin tracking for inbox members
ALTER TABLE "wzInboxMembers" ADD COLUMN IF NOT EXISTS "lastAssignedAt" TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_inbox_members_last_assigned ON "wzInboxMembers"("lastAssignedAt");
