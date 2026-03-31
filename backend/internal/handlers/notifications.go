package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type NotificationHandler struct {
	svc *services.NotificationService
}

func NewNotificationHandler(svc *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) List(c fiber.Ctx) error {
	userID := helpers.GetUserID(c)
	accountID := helpers.GetAccountID(c)
	page := helpers.ParseQueryInt(c, "page", 1)
	limit := helpers.ParseQueryInt(c, "limit", 20)

	notifications, unreadCount, err := h.svc.ListByUser(c.Context(), userID, accountID, page, limit)
	if err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

func (h *NotificationHandler) UnreadCount(c fiber.Ctx) error {
	userID := helpers.GetUserID(c)
	accountID := helpers.GetAccountID(c)

	notifications, unreadCount, err := h.svc.ListByUser(c.Context(), userID, accountID, 1, 1)
	if err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"unread_count":  unreadCount,
		"notifications": notifications,
	})
}

func (h *NotificationHandler) MarkAsRead(c fiber.Ctx) error {
	userID := helpers.GetUserID(c)
	accountID := helpers.GetAccountID(c)

	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid notification id")
	}

	if err := h.svc.MarkAsRead(c.Context(), userID, accountID, id); err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *NotificationHandler) MarkAllAsRead(c fiber.Ctx) error {
	userID := helpers.GetUserID(c)
	accountID := helpers.GetAccountID(c)

	if err := h.svc.MarkAllAsRead(c.Context(), userID, accountID); err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *NotificationHandler) GetSettings(c fiber.Ctx) error {
	userID := helpers.GetUserID(c)
	accountID := helpers.GetAccountID(c)

	settings, err := h.svc.GetSettings(c.Context(), userID, accountID)
	if err != nil {
		return helpers.NotFound(c, "settings not found")
	}

	return c.JSON(settings)
}

func (h *NotificationHandler) UpdateSettings(c fiber.Ctx) error {
	userID := helpers.GetUserID(c)
	accountID := helpers.GetAccountID(c)

	var body struct {
		EmailEnabled bool `json:"email_enabled"`
		PushEnabled  bool `json:"push_enabled"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	if err := h.svc.UpdateSettings(c.Context(), userID, accountID, body.EmailEnabled, body.PushEnabled); err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(fiber.Map{"success": true})
}
