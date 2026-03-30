package repo

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

type NotificationRepo struct {
	db *bun.DB
}

func NewNotificationRepo(db *bun.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(ctx context.Context, notification *models.Notification) error {
	_, err := r.db.NewInsert().Model(notification).Exec(ctx)
	return err
}

func (r *NotificationRepo) ListByUser(ctx context.Context, userID, accountID int64, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.NewSelect().
		Model(&notifications).
		Where("user_id = ? AND account_id = ?", userID, accountID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return notifications, err
}

func (r *NotificationRepo) UnreadCount(ctx context.Context, userID, accountID int64) (int, error) {
	count, err := r.db.NewSelect().
		Model(&models.Notification{}).
		Where("user_id = ? AND account_id = ? AND read_at IS NULL", userID, accountID).
		Count(ctx)
	return count, err
}

func (r *NotificationRepo) MarkAsRead(ctx context.Context, userID, accountID int64, notificationID int64) error {
	_, err := r.db.NewUpdate().
		Model(&models.Notification{}).
		Set("read_at = NOW()").
		Where("id = ? AND user_id = ? AND account_id = ?", notificationID, userID, accountID).
		Exec(ctx)
	return err
}

func (r *NotificationRepo) MarkAllAsRead(ctx context.Context, userID, accountID int64) error {
	_, err := r.db.NewUpdate().
		Model(&models.Notification{}).
		Set("read_at = NOW()").
		Where("user_id = ? AND account_id = ? AND read_at IS NULL", userID, accountID).
		Exec(ctx)
	return err
}

func (r *NotificationRepo) GetSettings(ctx context.Context, userID, accountID int64) (*models.NotificationSetting, error) {
	var setting models.NotificationSetting
	err := r.db.NewSelect().
		Model(&setting).
		Where("user_id = ? AND account_id = ?", userID, accountID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *NotificationRepo) UpdateSettings(ctx context.Context, setting *models.NotificationSetting) error {
	_, err := r.db.NewInsert().
		Model(setting).
		On("conflict (user_id, account_id) DO UPDATE").
		Set("email_enabled = EXCLUDED.email_enabled, push_enabled = EXCLUDED.push_enabled, updated_at = NOW()").
		Exec(ctx)
	return err
}
