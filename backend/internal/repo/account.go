package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, acc *domain.Account) error {
	q := `INSERT INTO "wzAccounts" ("id","name","domain","apiKey","settings","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, q, acc.ID, acc.Name, acc.Domain, acc.APIKey, acc.Settings, acc.CreatedAt, acc.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}
	return nil
}

func (r *AccountRepository) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	q := `SELECT "id","name",COALESCE("domain",''),"apiKey","settings","createdAt","updatedAt"
		  FROM "wzAccounts" WHERE "id"=$1`
	var a domain.Account
	err := r.db.QueryRow(ctx, q, id).Scan(&a.ID, &a.Name, &a.Domain, &a.APIKey, &a.Settings, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return &a, nil
}

func (r *AccountRepository) FindByAPIKey(ctx context.Context, apiKey string) (*domain.Account, error) {
	q := `SELECT "id","name",COALESCE("domain",''),"apiKey","settings","createdAt","updatedAt"
		  FROM "wzAccounts" WHERE "apiKey"=$1`
	var a domain.Account
	err := r.db.QueryRow(ctx, q, apiKey).Scan(&a.ID, &a.Name, &a.Domain, &a.APIKey, &a.Settings, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("account not found for apiKey: %w", err)
	}
	return &a, nil
}

func (r *AccountRepository) FindAll(ctx context.Context) ([]domain.Account, error) {
	q := `SELECT "id","name",COALESCE("domain",''),"apiKey","settings","createdAt","updatedAt"
		  FROM "wzAccounts" ORDER BY "createdAt" DESC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer rows.Close()

	var list []domain.Account
	for rows.Next() {
		var a domain.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Domain, &a.APIKey, &a.Settings, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzAccounts" WHERE "id"=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete account %s: %w", id, err)
	}
	return nil
}

func (r *AccountRepository) UpdatedAt(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `UPDATE "wzAccounts" SET "updatedAt"=$1 WHERE "id"=$2`, time.Now(), id)
	return err
}
