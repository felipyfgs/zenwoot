package services

import (
	"context"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type NoteService struct {
	noteRepo *repo.NoteRepo
}

func NewNoteService(noteRepo *repo.NoteRepo) *NoteService {
	return &NoteService{noteRepo: noteRepo}
}

func (s *NoteService) ListByContact(ctx context.Context, accountID, contactID int64) ([]*models.Note, error) {
	return s.noteRepo.ListByContact(ctx, accountID, contactID)
}

func (s *NoteService) GetByID(ctx context.Context, accountID, id int64) (*models.Note, error) {
	return s.noteRepo.GetByID(ctx, accountID, id)
}

func (s *NoteService) Create(ctx context.Context, m *models.Note) (*models.Note, error) {
	if err := s.noteRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *NoteService) Update(ctx context.Context, m *models.Note) (*models.Note, error) {
	if err := s.noteRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *NoteService) Delete(ctx context.Context, accountID, id int64) error {
	return s.noteRepo.Delete(ctx, accountID, id)
}
