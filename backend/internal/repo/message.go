package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{db: db}
}

const messageCols = `"id","accountId","conversationId","inboxId",
	COALESCE("contactId",''),COALESCE("externalId",''),
	"direction","contentType",COALESCE("content",''),
	COALESCE("mediaUrl",''),COALESCE("mediaType",''),
	"metadata","status","createdAt"`

func scanMessage(row interface{ Scan(...any) error }) (*domain.Message, error) {
	var m domain.Message
	err := row.Scan(
		&m.ID, &m.AccountID, &m.ConversationID, &m.InboxID,
		&m.ContactID, &m.ExternalID,
		&m.Direction, &m.ContentType, &m.Content,
		&m.MediaURL, &m.MediaType,
		&m.Metadata, &m.Status, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MessageRepository) Create(ctx context.Context, msg *domain.Message) error {
	q := `INSERT INTO "wzMessages" ("id","accountId","conversationId","inboxId","contactId","externalId","direction","contentType","content","mediaUrl","mediaType","metadata","status","createdAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := r.db.Exec(ctx, q,
		msg.ID, msg.AccountID, msg.ConversationID, msg.InboxID,
		nilIfEmpty(msg.ContactID), nilIfEmpty(msg.ExternalID),
		msg.Direction, msg.ContentType, msg.Content,
		msg.MediaURL, msg.MediaType,
		msg.Metadata, msg.Status, msg.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}
	return nil
}

func (r *MessageRepository) FindByConversation(ctx context.Context, conversationID string, limit, offset int) ([]domain.Message, error) {
	q := `SELECT ` + messageCols + ` FROM "wzMessages"
		  WHERE "conversationId"=$1 ORDER BY "createdAt" ASC
		  LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, q, conversationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var list []domain.Message
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *msg)
	}
	return list, rows.Err()
}

func (r *MessageRepository) FindByExternalID(ctx context.Context, externalID string) (*domain.Message, error) {
	q := `SELECT ` + messageCols + ` FROM "wzMessages" WHERE "externalId"=$1 LIMIT 1`
	msg, err := scanMessage(r.db.QueryRow(ctx, q, externalID))
	if err != nil {
		return nil, fmt.Errorf("message not found: %w", err)
	}
	return msg, nil
}

func (r *MessageRepository) UpdateStatus(ctx context.Context, externalID string, status domain.MessageStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzMessages" SET "status"=$1 WHERE "externalId"=$2`, status, externalID)
	return err
}
