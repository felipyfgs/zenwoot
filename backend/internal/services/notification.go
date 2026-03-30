package services

import (
	"context"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type NotificationService struct {
	repo *repo.NotificationRepo
}

func NewNotificationService(repo *repo.NotificationRepo) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Create(ctx context.Context, notification *models.Notification) error {
	return s.repo.Create(ctx, notification)
}

func (s *NotificationService) NotifyUser(ctx context.Context, userID, accountID int64, notificationType, title, body string) error {
	notification := &models.Notification{
		UserID:           userID,
		AccountID:        accountID,
		NotificationType: notificationType,
		Title:            title,
		Body:             &body,
	}
	return s.repo.Create(ctx, notification)
}

func (s *NotificationService) ListByUser(ctx context.Context, userID, accountID int64, page, limit int) ([]models.Notification, int, error) {
	offset := (page - 1) * limit
	notifications, err := s.repo.ListByUser(ctx, userID, accountID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, _ := s.repo.UnreadCount(ctx, userID, accountID)
	return notifications, count, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, userID, accountID int64, notificationID int64) error {
	return s.repo.MarkAsRead(ctx, userID, accountID, notificationID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID, accountID int64) error {
	return s.repo.MarkAllAsRead(ctx, userID, accountID)
}

func (s *NotificationService) GetSettings(ctx context.Context, userID, accountID int64) (*models.NotificationSetting, error) {
	return s.repo.GetSettings(ctx, userID, accountID)
}

func (s *NotificationService) UpdateSettings(ctx context.Context, userID, accountID int64, emailEnabled, pushEnabled bool) error {
	setting := &models.NotificationSetting{
		UserID:       userID,
		AccountID:    accountID,
		EmailEnabled: emailEnabled,
		PushEnabled:  pushEnabled,
	}
	return s.repo.UpdateSettings(ctx, setting)
}
