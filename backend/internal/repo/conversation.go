package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type ConversationRepository struct {
	db *pgxpool.Pool
}

func NewConversationRepository(db *pgxpool.Pool) *ConversationRepository {
	return &ConversationRepository{db: db}
}

const conversationCols = `"id","accountId","inboxId",COALESCE("contactId",''),
	COALESCE("assigneeId",''),COALESCE("teamId",''),
	"identifier",COALESCE("lastMessage",''),"lastMessageAt","lastActivityAt",
	"unreadCount","status","priority","muted","snoozedUntil","metadata","createdAt","updatedAt"`

func scanConversation(row interface{ Scan(...any) error }) (*domain.Conversation, error) {
	var c domain.Conversation
	err := row.Scan(
		&c.ID, &c.AccountID, &c.InboxID, &c.ContactID,
		&c.AssigneeID, &c.TeamID,
		&c.Identifier, &c.LastMessage, &c.LastMessageAt, &c.LastActivityAt,
		&c.UnreadCount, &c.Status, &c.Priority, &c.Muted, &c.SnoozedUntil, &c.Metadata, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConversationRepository) Upsert(ctx context.Context, c *domain.Conversation) error {
	q := `INSERT INTO "wzConversations" ("id","accountId","inboxId","contactId","identifier","lastMessage","lastMessageAt","unreadCount","status","metadata","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		  ON CONFLICT ("inboxId","identifier") DO UPDATE SET
		    "lastMessage"=EXCLUDED."lastMessage",
		    "lastMessageAt"=EXCLUDED."lastMessageAt",
		    "unreadCount"="wzConversations"."unreadCount"+1,
		    "updatedAt"=NOW()
		  RETURNING "id"`
	err := r.db.QueryRow(ctx, q,
		c.ID, c.AccountID, c.InboxID, nilIfEmpty(c.ContactID),
		c.Identifier, c.LastMessage, c.LastMessageAt,
		c.UnreadCount, c.Status, c.Metadata, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID)
	if err != nil {
		return fmt.Errorf("failed to upsert conversation: %w", err)
	}
	return nil
}

func (r *ConversationRepository) FindByID(ctx context.Context, id string) (*domain.Conversation, error) {
	q := `SELECT ` + conversationCols + ` FROM "wzConversations" WHERE "id"=$1`
	return scanConversation(r.db.QueryRow(ctx, q, id))
}

func (r *ConversationRepository) FindByInboxAndIdentifier(ctx context.Context, inboxID string, identifier string) (*domain.Conversation, error) {
	q := `SELECT ` + conversationCols + ` FROM "wzConversations" WHERE "inboxId"=$1 AND "identifier"=$2 LIMIT 1`
	c, err := scanConversation(r.db.QueryRow(ctx, q, inboxID, identifier))
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}
	return c, nil
}

func (r *ConversationRepository) FindByInboxID(ctx context.Context, inboxID string, limit, offset int) ([]domain.Conversation, error) {
	q := `SELECT ` + conversationCols + ` FROM "wzConversations"
		  WHERE "inboxId"=$1 ORDER BY COALESCE("lastMessageAt","createdAt") DESC
		  LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, q, inboxID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations: %w", err)
	}
	defer rows.Close()

	var list []domain.Conversation
	for rows.Next() {
		c, err := scanConversation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *c)
	}
	return list, rows.Err()
}

func (r *ConversationRepository) FindByContactID(ctx context.Context, contactID string, limit, offset int) ([]domain.Conversation, error) {
	q := `SELECT ` + conversationCols + ` FROM "wzConversations"
		  WHERE "contactId"=$1 ORDER BY COALESCE("lastMessageAt","createdAt") DESC
		  LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, q, contactID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations by contact: %w", err)
	}
	defer rows.Close()

	var list []domain.Conversation
	for rows.Next() {
		c, err := scanConversation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *c)
	}
	return list, rows.Err()
}

func (r *ConversationRepository) UpdateStatus(ctx context.Context, id string, status domain.ConversationStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzConversations" SET "status"=$1,"updatedAt"=NOW() WHERE "id"=$2`, status, id)
	return err
}

func (r *ConversationRepository) MarkRead(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzConversations" SET "unreadCount"=0,"updatedAt"=NOW() WHERE "id"=$1`, id)
	return err
}

func (r *ConversationRepository) UpdateLastMessage(ctx context.Context, id string, text string, at time.Time) error {
	_, err := r.db.Exec(ctx,
		`UPDATE "wzConversations" SET "lastMessage"=$1,"lastMessageAt"=$2,"lastActivityAt"=$3,"updatedAt"=NOW() WHERE "id"=$4`,
		text, at, at, id,
	)
	return err
}

func (r *ConversationRepository) UpdateAssigneeRaw(ctx context.Context, userID, teamID *string, now time.Time, id string) (interface{}, error) {
	q := `UPDATE "wzConversations" SET "assigneeId"=$1, "teamId"=$2, "lastActivityAt"=$3, "updatedAt"=$4 WHERE "id"=$5`
	_, err := r.db.Exec(ctx, q, userID, teamID, now, now, id)
	return nil, err
}

func (r *ConversationRepository) UpdateAssignee(ctx context.Context, id string, userID, teamID *string) error {
	now := time.Now()
	q := `UPDATE "wzConversations" SET "assigneeId"=$1, "teamId"=$2, "lastActivityAt"=$3, "updatedAt"=$4 WHERE "id"=$5`
	_, err := r.db.Exec(ctx, q, userID, teamID, now, now, id)
	return err
}

func (r *ConversationRepository) UpdatePriority(ctx context.Context, id string, priority domain.ConversationPriority) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzConversations" SET "priority"=$1,"updatedAt"=NOW() WHERE "id"=$2`, priority, id)
	return err
}

func (r *ConversationRepository) Snooze(ctx context.Context, id string, until time.Time) error {
	_, err := r.db.Exec(ctx,
		`UPDATE "wzConversations" SET "status"='snoozed', "snoozedUntil"=$1, "updatedAt"=NOW() WHERE "id"=$2`,
		until, id,
	)
	return err
}

func (r *ConversationRepository) UnSnooze(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE "wzConversations" SET "status"='open', "snoozedUntil"=NULL, "updatedAt"=NOW() WHERE "id"=$1`,
		id,
	)
	return err
}

func (r *ConversationRepository) UpdateMuted(ctx context.Context, id string, muted bool) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzConversations" SET "muted"=$1,"updatedAt"=NOW() WHERE "id"=$2`, muted, id)
	return err
}

func (r *ConversationRepository) FindByAssignee(ctx context.Context, userID string, limit, offset int) ([]domain.Conversation, error) {
	q := `SELECT ` + conversationCols + ` FROM "wzConversations"
		  WHERE "assigneeId"=$1 ORDER BY COALESCE("lastActivityAt","lastMessageAt","createdAt") DESC
		  LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations by assignee: %w", err)
	}
	defer rows.Close()

	var list []domain.Conversation
	for rows.Next() {
		c, err := scanConversation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *c)
	}
	return list, rows.Err()
}

func (r *ConversationRepository) FindByTeam(ctx context.Context, teamID string, limit, offset int) ([]domain.Conversation, error) {
	q := `SELECT ` + conversationCols + ` FROM "wzConversations"
		  WHERE "teamId"=$1 ORDER BY COALESCE("lastActivityAt","lastMessageAt","createdAt") DESC
		  LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, q, teamID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations by team: %w", err)
	}
	defer rows.Close()

	var list []domain.Conversation
	for rows.Next() {
		c, err := scanConversation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *c)
	}
	return list, rows.Err()
}

func (r *ConversationRepository) FindByStatus(ctx context.Context, inboxID string, status domain.ConversationStatus, limit, offset int) ([]domain.Conversation, error) {
	q := `SELECT ` + conversationCols + ` FROM "wzConversations"
		  WHERE "inboxId"=$1 AND "status"=$2 ORDER BY COALESCE("lastActivityAt","lastMessageAt","createdAt") DESC
		  LIMIT $3 OFFSET $4`
	rows, err := r.db.Query(ctx, q, inboxID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations by status: %w", err)
	}
	defer rows.Close()

	var list []domain.Conversation
	for rows.Next() {
		c, err := scanConversation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *c)
	}
	return list, rows.Err()
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
