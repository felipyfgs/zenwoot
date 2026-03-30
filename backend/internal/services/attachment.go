package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/storage"
)

type AttachmentService struct {
	db    *bun.DB
	store *storage.S3Storage
}

func NewAttachmentService(db *bun.DB, store *storage.S3Storage) *AttachmentService {
	return &AttachmentService{db: db, store: store}
}

type CreateAttachmentInput struct {
	MessageID int64
	AccountID int64
	File      *multipart.FileHeader
	FileType  int
	Extension string
}

func (s *AttachmentService) Create(ctx context.Context, input CreateAttachmentInput) (*models.Attachment, error) {
	objectName, err := s.store.Upload(ctx, input.File, "attachments")
	if err != nil {
		return nil, fmt.Errorf("upload: %w", err)
	}

	attachment := &models.Attachment{
		MessageID:   input.MessageID,
		AccountID:   input.AccountID,
		FileType:    input.FileType,
		ExternalURL: &objectName,
		Extension:   &input.Extension,
	}

	_, err = s.db.NewInsert().Model(attachment).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}

	return attachment, nil
}

func (s *AttachmentService) GetByMessageID(ctx context.Context, messageID int64) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := s.db.NewSelect().
		Model(&attachments).
		Where("message_id = ?", messageID).
		Scan(ctx)
	return attachments, err
}

func (s *AttachmentService) GetByID(ctx context.Context, id int64) (*models.Attachment, error) {
	if id == 0 {
		return nil, nil
	}
	var attachment models.Attachment
	err := s.db.NewSelect().
		Model(&attachment).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	if attachment.ID == 0 {
		return nil, nil
	}
	return &attachment, nil
}

func (s *AttachmentService) GetDownloadURL(ctx context.Context, attachment *models.Attachment) (string, error) {
	if attachment == nil {
		return "", errors.New("attachment is nil")
	}
	if attachment.ExternalURL == nil {
		return "", nil
	}
	return s.store.GetURL(ctx, *attachment.ExternalURL)
}
