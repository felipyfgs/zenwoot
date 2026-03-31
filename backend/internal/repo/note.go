package repo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type NoteRepo struct {
	BaseRepo[models.Note]
}

func NewNoteRepo(db *bun.DB) *NoteRepo {
	return &NoteRepo{BaseRepo: *NewBaseRepo[models.Note](db)}
}

func (r *NoteRepo) ListByContact(ctx context.Context, accountID, contactID int64) ([]*models.Note, error) {
	var items []*models.Note
	err := r.DB().NewSelect().Model(&items).
		Where(`"account_id" = ? AND "contact_id" = ?`, accountID, contactID).
		OrderExpr(`"created_at" DESC`).
		Scan(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("noteRepo.ListByContact: %w", err)
	}
	return items, nil
}

func (r *NoteRepo) Create(ctx context.Context, m *models.Note) error {
	_, err := r.DB().NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return fmt.Errorf("noteRepo.Create: %w", err)
	}
	return nil
}

func (r *NoteRepo) Update(ctx context.Context, m *models.Note) error {
	_, err := r.DB().NewUpdate().Model(m).
		Where(`"id" = ? AND "account_id" = ?`, m.ID, m.AccountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("noteRepo.Update: %w", err)
	}
	return nil
}

func (r *NoteRepo) Delete(ctx context.Context, accountID, id int64) error {
	_, err := r.DB().NewDelete().TableExpr(`"notes"`).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("noteRepo.Delete: %w", err)
	}
	return nil
}

func (r *NoteRepo) GetByID(ctx context.Context, accountID, id int64) (*models.Note, error) {
	var m models.Note
	err := r.DB().NewSelect().Model(&m).
		Where(`"id" = ? AND "account_id" = ?`, id, accountID).
		Scan(ctx, &m)
	if err != nil {
		return nil, fmt.Errorf("noteRepo.GetByID: %w", err)
	}
	return &m, nil
}
