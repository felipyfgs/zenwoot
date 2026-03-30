package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type NotificationHandler struct {
	svc *services.NotificationService
}

func NewNotificationHandler(svc *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) Register(rg fiber.Router) {
	rg.Get("/notifications", h.List)
	rg.Get("/notifications/unread_count", h.UnreadCount)
	rg.Post("/notifications/:id/read", h.MarkAsRead)
	rg.Post("/notifications/read_all", h.MarkAllAsRead)
	rg.Get("/notification_settings", h.GetSettings)
	rg.Put("/notification_settings", h.UpdateSettings)
}

func (h *NotificationHandler) List(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	accountID := c.Locals("account_id").(int64)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	notifications, unreadCount, err := h.svc.ListByUser(c.Context(), userID, accountID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

func (h *NotificationHandler) UnreadCount(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	accountID := c.Locals("account_id").(int64)

	notifications, unreadCount, err := h.svc.ListByUser(c.Context(), userID, accountID, 1, 1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"unread_count":  unreadCount,
		"notifications": notifications,
	})
}

func (h *NotificationHandler) MarkAsRead(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)

	err := h.svc.MarkAsRead(c.Context(), userID, accountID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *NotificationHandler) MarkAllAsRead(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	accountID := c.Locals("account_id").(int64)

	err := h.svc.MarkAllAsRead(c.Context(), userID, accountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *NotificationHandler) GetSettings(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	accountID := c.Locals("account_id").(int64)

	settings, err := h.svc.GetSettings(c.Context(), userID, accountID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "settings not found"})
	}

	return c.JSON(settings)
}

func (h *NotificationHandler) UpdateSettings(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	accountID := c.Locals("account_id").(int64)

	var body struct {
		EmailEnabled bool `json:"email_enabled"`
		PushEnabled  bool `json:"push_enabled"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	err := h.svc.UpdateSettings(c.Context(), userID, accountID, body.EmailEnabled, body.PushEnabled)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
