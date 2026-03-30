package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type ContactRepository struct {
	db *pgxpool.Pool
}

func NewContactRepository(db *pgxpool.Pool) *ContactRepository {
	return &ContactRepository{db: db}
}

const contactCols = `"id","accountId","identifier",COALESCE("name",''),
	COALESCE("pushName",''),COALESCE("avatarUrl",''),
	"isBlocked","metadata","createdAt","updatedAt"`

func scanContact(row interface{ Scan(...any) error }) (*domain.Contact, error) {
	var c domain.Contact
	err := row.Scan(
		&c.ID, &c.AccountID, &c.Identifier, &c.Name,
		&c.PushName, &c.AvatarURL, &c.IsBlocked,
		&c.Metadata, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContactRepository) Upsert(ctx context.Context, c *domain.Contact) error {
	q := `INSERT INTO "wzContacts" ("id","accountId","identifier","name","pushName","avatarUrl","isBlocked","metadata","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		  ON CONFLICT ("accountId","identifier") DO UPDATE SET
		    "name"=COALESCE(EXCLUDED."name","wzContacts"."name"),
		    "pushName"=COALESCE(EXCLUDED."pushName","wzContacts"."pushName"),
		    "updatedAt"=NOW()
		  RETURNING "id"`
	err := r.db.QueryRow(ctx, q,
		c.ID, c.AccountID, c.Identifier, c.Name,
		c.PushName, c.AvatarURL, c.IsBlocked, c.Metadata, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID)
	if err != nil {
		return fmt.Errorf("failed to upsert contact: %w", err)
	}
	return nil
}

func (r *ContactRepository) FindByID(ctx context.Context, id string) (*domain.Contact, error) {
	q := `SELECT ` + contactCols + ` FROM "wzContacts" WHERE "id"=$1`
	c, err := scanContact(r.db.QueryRow(ctx, q, id))
	if err != nil {
		return nil, fmt.Errorf("contact not found: %w", err)
	}
	return c, nil
}

func (r *ContactRepository) FindByIdentifier(ctx context.Context, accountID, identifier string) (*domain.Contact, error) {
	q := `SELECT ` + contactCols + ` FROM "wzContacts" WHERE "accountId"=$1 AND "identifier"=$2`
	c, err := scanContact(r.db.QueryRow(ctx, q, accountID, identifier))
	if err != nil {
		return nil, fmt.Errorf("contact not found: %w", err)
	}
	return c, nil
}

func (r *ContactRepository) FindByAccountID(ctx context.Context, accountID string, limit, offset int) ([]domain.Contact, error) {
	q := `SELECT ` + contactCols + ` FROM "wzContacts"
		  WHERE "accountId"=$1 ORDER BY "name" ASC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, q, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query contacts: %w", err)
	}
	defer rows.Close()

	var list []domain.Contact
	for rows.Next() {
		c, err := scanContact(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *c)
	}
	return list, rows.Err()
}
