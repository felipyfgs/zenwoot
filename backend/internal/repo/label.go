package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type LabelRepository struct {
	db *pgxpool.Pool
}

func NewLabelRepository(db *pgxpool.Pool) *LabelRepository {
	return &LabelRepository{db: db}
}

func (r *LabelRepository) Create(ctx context.Context, l *domain.Label) error {
	q := `INSERT INTO "wzLabels" ("id","accountId","title","color","description","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, q, l.ID, l.AccountID, l.Title, l.Color, l.Description, l.CreatedAt, l.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create label: %w", err)
	}
	return nil
}

func (r *LabelRepository) FindByID(ctx context.Context, id string) (*domain.Label, error) {
	q := `SELECT "id","accountId","title","color",COALESCE("description",''),"createdAt","updatedAt"
		  FROM "wzLabels" WHERE "id"=$1`
	var l domain.Label
	err := r.db.QueryRow(ctx, q, id).Scan(
		&l.ID, &l.AccountID, &l.Title, &l.Color, &l.Description, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("label not found: %w", err)
	}
	return &l, nil
}

func (r *LabelRepository) FindByAccountID(ctx context.Context, accountID string) ([]domain.Label, error) {
	q := `SELECT "id","accountId","title","color",COALESCE("description",''),"createdAt","updatedAt"
		  FROM "wzLabels" WHERE "accountId"=$1 ORDER BY "title" ASC`
	rows, err := r.db.Query(ctx, q, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query labels: %w", err)
	}
	defer rows.Close()

	var list []domain.Label
	for rows.Next() {
		var l domain.Label
		if err := rows.Scan(&l.ID, &l.AccountID, &l.Title, &l.Color, &l.Description, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, rows.Err()
}

func (r *LabelRepository) Update(ctx context.Context, l *domain.Label) error {
	q := `UPDATE "wzLabels" SET "title"=$1,"color"=$2,"description"=$3,"updatedAt"=NOW() WHERE "id"=$4`
	_, err := r.db.Exec(ctx, q, l.Title, l.Color, l.Description, l.ID)
	return err
}

func (r *LabelRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzLabels" WHERE "id"=$1`, id)
	return err
}

// ConversationLabel Repository

type ConversationLabelRepository struct {
	db *pgxpool.Pool
}

func NewConversationLabelRepository(db *pgxpool.Pool) *ConversationLabelRepository {
	return &ConversationLabelRepository{db: db}
}

func (r *ConversationLabelRepository) Create(ctx context.Context, cl *domain.ConversationLabel) error {
	q := `INSERT INTO "wzConversationLabels" ("id","conversationId","labelId","createdAt")
		  VALUES ($1,$2,$3,$4)`
	_, err := r.db.Exec(ctx, q, cl.ID, cl.ConversationID, cl.LabelID, cl.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to add label to conversation: %w", err)
	}
	return nil
}

func (r *ConversationLabelRepository) FindByConversationID(ctx context.Context, conversationID string) ([]domain.Label, error) {
	q := `SELECT l."id",l."accountId",l."title",l."color",COALESCE(l."description",''),l."createdAt",l."updatedAt"
		  FROM "wzLabels" l
		  INNER JOIN "wzConversationLabels" cl ON l."id" = cl."labelId"
		  WHERE cl."conversationId"=$1
		  ORDER BY l."title" ASC`
	rows, err := r.db.Query(ctx, q, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversation labels: %w", err)
	}
	defer rows.Close()

	var list []domain.Label
	for rows.Next() {
		var l domain.Label
		if err := rows.Scan(&l.ID, &l.AccountID, &l.Title, &l.Color, &l.Description, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, rows.Err()
}

func (r *ConversationLabelRepository) Delete(ctx context.Context, conversationID, labelID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzConversationLabels" WHERE "conversationId"=$1 AND "labelId"=$2`, conversationID, labelID)
	return err
}

func (r *ConversationLabelRepository) DeleteAllByConversation(ctx context.Context, conversationID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzConversationLabels" WHERE "conversationId"=$1`, conversationID)
	return err
}
