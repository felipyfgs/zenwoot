package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"wzap/internal/domain"
)

type NoteRepository struct {
	db *pgxpool.Pool
}

func NewNoteRepository(db *pgxpool.Pool) *NoteRepository {
	return &NoteRepository{db: db}
}

const noteCols = `"id","accountId","contactId",COALESCE("userId",''),"content","createdAt","updatedAt"`

func scanNote(row interface{ Scan(...any) error }) (*domain.Note, error) {
	var n domain.Note
	err := row.Scan(&n.ID, &n.AccountID, &n.ContactID, &n.UserID, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NoteRepository) Create(ctx context.Context, n *domain.Note) error {
	q := `INSERT INTO "wzNotes" ("id","accountId","contactId","userId","content","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, q,
		n.ID, n.AccountID, n.ContactID, nilIfEmpty(n.UserID), n.Content, n.CreatedAt, n.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}
	return nil
}

func (r *NoteRepository) FindByID(ctx context.Context, id string) (*domain.Note, error) {
	q := `SELECT ` + noteCols + ` FROM "wzNotes" WHERE "id"=$1`
	return scanNote(r.db.QueryRow(ctx, q, id))
}

func (r *NoteRepository) FindByContactID(ctx context.Context, contactID string) ([]domain.Note, error) {
	q := `SELECT ` + noteCols + ` FROM "wzNotes" WHERE "contactId"=$1 ORDER BY "createdAt" DESC`
	rows, err := r.db.Query(ctx, q, contactID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var list []domain.Note
	for rows.Next() {
		n, err := scanNote(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *n)
	}
	return list, rows.Err()
}

func (r *NoteRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzNotes" WHERE "id"=$1`, id)
	return err
}

// ConversationParticipant Repository

type ConversationParticipantRepository struct {
	db *pgxpool.Pool
}

func NewConversationParticipantRepository(db *pgxpool.Pool) *ConversationParticipantRepository {
	return &ConversationParticipantRepository{db: db}
}

func (r *ConversationParticipantRepository) Create(ctx context.Context, p *domain.ConversationParticipant) error {
	q := `INSERT INTO "wzConversationParticipants" ("id","accountId","conversationId","userId","createdAt","updatedAt")
		  VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(ctx, q, p.ID, p.AccountID, p.ConversationID, p.UserID, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}
	return nil
}

func (r *ConversationParticipantRepository) FindByConversationID(ctx context.Context, conversationID string) ([]domain.ConversationParticipant, error) {
	q := `SELECT "id","accountId","conversationId","userId","createdAt","updatedAt"
		  FROM "wzConversationParticipants" WHERE "conversationId"=$1`
	rows, err := r.db.Query(ctx, q, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query participants: %w", err)
	}
	defer rows.Close()

	var list []domain.ConversationParticipant
	for rows.Next() {
		var p domain.ConversationParticipant
		if err := rows.Scan(&p.ID, &p.AccountID, &p.ConversationID, &p.UserID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *ConversationParticipantRepository) Delete(ctx context.Context, conversationID, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "wzConversationParticipants" WHERE "conversationId"=$1 AND "userId"=$2`, conversationID, userID)
	return err
}

func (r *ConversationParticipantRepository) FindByUserID(ctx context.Context, userID string) ([]domain.ConversationParticipant, error) {
	q := `SELECT "id","accountId","conversationId","userId","createdAt","updatedAt"
		  FROM "wzConversationParticipants" WHERE "userId"=$1`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query participant conversations: %w", err)
	}
	defer rows.Close()

	var list []domain.ConversationParticipant
	for rows.Next() {
		var p domain.ConversationParticipant
		if err := rows.Scan(&p.ID, &p.AccountID, &p.ConversationID, &p.UserID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}
