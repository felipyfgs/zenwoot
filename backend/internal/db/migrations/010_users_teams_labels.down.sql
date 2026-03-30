-- Drop new columns from messages
DROP INDEX IF EXISTS idx_messages_reply_to;
DROP INDEX IF EXISTS idx_messages_sender;
ALTER TABLE "wzMessages" DROP COLUMN IF EXISTS "replyToId";
ALTER TABLE "wzMessages" DROP COLUMN IF EXISTS "senderId";

-- Drop new columns from conversations
DROP INDEX IF EXISTS idx_conversations_status;
DROP INDEX IF EXISTS idx_conversations_team;
DROP INDEX IF EXISTS idx_conversations_assignee;
ALTER TABLE "wzConversations" DROP COLUMN IF EXISTS "lastActivityAt";
ALTER TABLE "wzConversations" DROP COLUMN IF EXISTS "snoozedUntil";
ALTER TABLE "wzConversations" DROP COLUMN IF EXISTS "muted";
ALTER TABLE "wzConversations" DROP COLUMN IF EXISTS "priority";
ALTER TABLE "wzConversations" DROP COLUMN IF EXISTS "teamId";
ALTER TABLE "wzConversations" DROP COLUMN IF EXISTS "assigneeId";

-- Drop tables in reverse order
DROP TABLE IF EXISTS "wzConversationLabels";
DROP TABLE IF EXISTS "wzLabels";
DROP TABLE IF EXISTS "wzInboxMembers";
DROP TABLE IF EXISTS "wzTeamMembers";
DROP TABLE IF EXISTS "wzTeams";
DROP TABLE IF EXISTS "wzAccountUsers";
DROP TABLE IF EXISTS "wzUsers";
