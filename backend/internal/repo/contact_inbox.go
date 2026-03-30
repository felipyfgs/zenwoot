package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type ContactInboxRepository struct {
	db *pgxpool.Pool
}

func NewContactInboxRepository(db *pgxpool.Pool) *ContactInboxRepository {
	return &ContactInboxRepository{db: db}
}

func (r *ContactInboxRepository) Upsert(ctx context.Context, ci *domain.ContactInbox) error {
	if ci.ID == "" {
		ci.ID = uuid.NewString()
	}
	now := time.Now()
	ci.CreatedAt = now
	ci.UpdatedAt = now

	q := `INSERT INTO "wzContactInboxes" ("id","contactId","inboxId","sourceId","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6)
		  ON CONFLICT ("inboxId","sourceId") DO UPDATE SET
		    "contactId" = EXCLUDED."contactId",
		    "updatedAt" = NOW()
		  RETURNING "id"`
	err := r.db.QueryRow(ctx, q, ci.ID, ci.ContactID, ci.InboxID, ci.SourceID, ci.CreatedAt, ci.UpdatedAt).Scan(&ci.ID)
	if err != nil {
		return fmt.Errorf("failed to upsert contact_inbox: %w", err)
	}
	return nil
}

func (r *ContactInboxRepository) FindByInboxAndSourceID(ctx context.Context, inboxID, sourceID string) (*domain.ContactInbox, error) {
	q := `SELECT "id","contactId","inboxId","sourceId","createdAt","updatedAt"
		  FROM "wzContactInboxes" WHERE "inboxId"=$1 AND "sourceId"=$2`
	var ci domain.ContactInbox
	err := r.db.QueryRow(ctx, q, inboxID, sourceID).Scan(&ci.ID, &ci.ContactID, &ci.InboxID, &ci.SourceID, &ci.CreatedAt, &ci.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("contact_inbox not found: %w", err)
	}
	return &ci, nil
}

func (r *ContactInboxRepository) FindByContactID(ctx context.Context, contactID string) ([]domain.ContactInbox, error) {
	q := `SELECT "id","contactId","inboxId","sourceId","createdAt","updatedAt"
		  FROM "wzContactInboxes" WHERE "contactId"=$1 ORDER BY "createdAt" DESC`
	rows, err := r.db.Query(ctx, q, contactID)
	if err != nil {
		return nil, fmt.Errorf("failed to query contact_inboxes: %w", err)
	}
	defer rows.Close()

	var list []domain.ContactInbox
	for rows.Next() {
		var ci domain.ContactInbox
		if err := rows.Scan(&ci.ID, &ci.ContactID, &ci.InboxID, &ci.SourceID, &ci.CreatedAt, &ci.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, ci)
	}
	return list, rows.Err()
}

func (r *ContactInboxRepository) FindByInboxID(ctx context.Context, inboxID string) ([]domain.ContactInbox, error) {
	q := `SELECT "id","contactId","inboxId","sourceId","createdAt","updatedAt"
		  FROM "wzContactInboxes" WHERE "inboxId"=$1 ORDER BY "createdAt" DESC`
	rows, err := r.db.Query(ctx, q, inboxID)
	if err != nil {
		return nil, fmt.Errorf("failed to query contact_inboxes: %w", err)
	}
	defer rows.Close()

	var list []domain.ContactInbox
	for rows.Next() {
		var ci domain.ContactInbox
		if err := rows.Scan(&ci.ID, &ci.ContactID, &ci.InboxID, &ci.SourceID, &ci.CreatedAt, &ci.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, ci)
	}
	return list, rows.Err()
}

func (r *ContactInboxRepository) DeleteByInboxID(ctx context.Context, inboxID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzContactInboxes" WHERE "inboxId"=$1`, inboxID)
	return err
}
