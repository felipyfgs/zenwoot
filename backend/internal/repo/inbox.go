package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type InboxRepository struct {
	db *pgxpool.Pool
}

func NewInboxRepository(db *pgxpool.Pool) *InboxRepository {
	return &InboxRepository{db: db}
}

func (r *InboxRepository) Create(ctx context.Context, inbox *domain.Inbox) error {
	q := `INSERT INTO "wzInboxes" ("id","accountId","name","channelType","channelId","status","settings","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.db.Exec(ctx, q,
		inbox.ID, inbox.AccountID, inbox.Name, inbox.ChannelType,
		inbox.ChannelID, inbox.Status, inbox.Settings, inbox.CreatedAt, inbox.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert inbox: %w", err)
	}
	return nil
}

func (r *InboxRepository) FindByID(ctx context.Context, id string) (*domain.Inbox, error) {
	q := `SELECT "id","accountId","name","channelType","channelId","status","settings","createdAt","updatedAt"
		  FROM "wzInboxes" WHERE "id"=$1`
	var i domain.Inbox
	err := r.db.QueryRow(ctx, q, id).Scan(
		&i.ID, &i.AccountID, &i.Name, &i.ChannelType,
		&i.ChannelID, &i.Status, &i.Settings, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("inbox not found: %w", err)
	}
	return &i, nil
}

// FindByNameOrID resolves an inbox by UUID or by name (case-sensitive).
func (r *InboxRepository) FindByNameOrID(ctx context.Context, nameOrID string) (*domain.Inbox, error) {
	q := `SELECT "id","accountId","name","channelType","channelId","status","settings","createdAt","updatedAt"
		  FROM "wzInboxes" WHERE "id"=$1 OR "name"=$1 LIMIT 1`
	var i domain.Inbox
	err := r.db.QueryRow(ctx, q, nameOrID).Scan(
		&i.ID, &i.AccountID, &i.Name, &i.ChannelType,
		&i.ChannelID, &i.Status, &i.Settings, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("inbox not found: %w", err)
	}
	return &i, nil
}

func (r *InboxRepository) FindByAccountID(ctx context.Context, accountID string) ([]domain.Inbox, error) {
	q := `SELECT "id","accountId","name","channelType","channelId","status","settings","createdAt","updatedAt"
		  FROM "wzInboxes" WHERE "accountId"=$1 ORDER BY "createdAt" DESC`
	rows, err := r.db.Query(ctx, q, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inboxes: %w", err)
	}
	defer rows.Close()

	var list []domain.Inbox
	for rows.Next() {
		var i domain.Inbox
		if err := rows.Scan(&i.ID, &i.AccountID, &i.Name, &i.ChannelType, &i.ChannelID, &i.Status, &i.Settings, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, i)
	}
	return list, rows.Err()
}

func (r *InboxRepository) FindByChannelID(ctx context.Context, channelType domain.ChannelType, channelID string) (*domain.Inbox, error) {
	q := `SELECT "id","accountId","name","channelType","channelId","status","settings","createdAt","updatedAt"
		  FROM "wzInboxes" WHERE "channelType"=$1 AND "channelId"=$2 LIMIT 1`
	var i domain.Inbox
	err := r.db.QueryRow(ctx, q, channelType, channelID).Scan(
		&i.ID, &i.AccountID, &i.Name, &i.ChannelType,
		&i.ChannelID, &i.Status, &i.Settings, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("inbox not found: %w", err)
	}
	return &i, nil
}

func (r *InboxRepository) UpdateStatus(ctx context.Context, id string, status domain.InboxStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzInboxes" SET "status"=$1 WHERE "id"=$2`, status, id)
	return err
}

func (r *InboxRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzInboxes" WHERE "id"=$1`, id)
	return err
}
