package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type CannedResponseRepository struct {
	db *pgxpool.Pool
}

func NewCannedResponseRepository(db *pgxpool.Pool) *CannedResponseRepository {
	return &CannedResponseRepository{db: db}
}

func (r *CannedResponseRepository) Create(ctx context.Context, cr *domain.CannedResponse) error {
	q := `INSERT INTO "wzCannedResponses" ("id","accountId","shortCode","content","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(ctx, q, cr.ID, cr.AccountID, cr.ShortCode, cr.Content, cr.CreatedAt, cr.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create canned response: %w", err)
	}
	return nil
}

func (r *CannedResponseRepository) FindByID(ctx context.Context, id string) (*domain.CannedResponse, error) {
	q := `SELECT "id","accountId","shortCode","content","createdAt","updatedAt"
		  FROM "wzCannedResponses" WHERE "id"=$1`
	var cr domain.CannedResponse
	err := r.db.QueryRow(ctx, q, id).Scan(
		&cr.ID, &cr.AccountID, &cr.ShortCode, &cr.Content, &cr.CreatedAt, &cr.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("canned response not found: %w", err)
	}
	return &cr, nil
}

func (r *CannedResponseRepository) FindByAccountID(ctx context.Context, accountID string) ([]domain.CannedResponse, error) {
	q := `SELECT "id","accountId","shortCode","content","createdAt","updatedAt"
		  FROM "wzCannedResponses" WHERE "accountId"=$1 ORDER BY "shortCode" ASC`
	rows, err := r.db.Query(ctx, q, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query canned responses: %w", err)
	}
	defer rows.Close()

	var list []domain.CannedResponse
	for rows.Next() {
		var cr domain.CannedResponse
		if err := rows.Scan(&cr.ID, &cr.AccountID, &cr.ShortCode, &cr.Content, &cr.CreatedAt, &cr.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, cr)
	}
	return list, rows.Err()
}

// Search performs a case-insensitive search on short_code and content.
func (r *CannedResponseRepository) Search(ctx context.Context, accountID, query string) ([]domain.CannedResponse, error) {
	q := `SELECT "id","accountId","shortCode","content","createdAt","updatedAt"
		  FROM "wzCannedResponses"
		  WHERE "accountId"=$1
		    AND ("shortCode" ILIKE $2 OR "content" ILIKE $2)
		  ORDER BY
		    CASE WHEN "shortCode" ILIKE $3 THEN 0
		         WHEN "shortCode" ILIKE $2 THEN 1
		         ELSE 2
		    END,
		    "shortCode" ASC
		  LIMIT 15`
	pattern := "%" + query + "%"
	prefix := query + "%"
	rows, err := r.db.Query(ctx, q, accountID, pattern, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to search canned responses: %w", err)
	}
	defer rows.Close()

	var list []domain.CannedResponse
	for rows.Next() {
		var cr domain.CannedResponse
		if err := rows.Scan(&cr.ID, &cr.AccountID, &cr.ShortCode, &cr.Content, &cr.CreatedAt, &cr.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, cr)
	}
	return list, rows.Err()
}

func (r *CannedResponseRepository) Update(ctx context.Context, cr *domain.CannedResponse) error {
	q := `UPDATE "wzCannedResponses" SET "shortCode"=$1,"content"=$2,"updatedAt"=NOW() WHERE "id"=$3`
	_, err := r.db.Exec(ctx, q, cr.ShortCode, cr.Content, cr.ID)
	if err != nil {
		return fmt.Errorf("failed to update canned response: %w", err)
	}
	return nil
}

func (r *CannedResponseRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzCannedResponses" WHERE "id"=$1`, id)
	return err
}
