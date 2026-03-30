package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type WebhookHandler struct {
	svc *services.WebhookService
}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

func NewWebhookHandlerWithService(svc *services.WebhookService) *WebhookHandler {
	return &WebhookHandler{svc: svc}
}

func (h *WebhookHandler) Register(rg fiber.Router) {
	rg.Post("/webhooks/whatsapp", h.WhatsApp)
	rg.Post("/webhooks/telegram", h.Telegram)
	rg.Post("/webhooks/email", h.Email)
	rg.Post("/webhooks/debug", h.Debug)
}

func (h *WebhookHandler) RegisterCRUD(rg fiber.Router) {
	rg.Get("/webhooks", h.List)
	rg.Post("/webhooks", h.Create)
	rg.Get("/webhooks/:id", h.Get)
	rg.Patch("/webhooks/:id", h.Update)
	rg.Delete("/webhooks/:id", h.Delete)
}

func (h *WebhookHandler) List(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)

	items, err := h.svc.List(c.Context(), accountID)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *WebhookHandler) Get(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid webhook id")
	}

	item, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "webhook not found")
	}
	return c.JSON(item)
}

func (h *WebhookHandler) Create(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	var body models.Webhook
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	body.AccountID = accountID
	item, err := h.svc.Create(c.Context(), &body)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return helpers.Created(c, item)
}

func (h *WebhookHandler) Update(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid webhook id")
	}

	existing, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "webhook not found")
	}

	if err := c.Bind().JSON(existing); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	updated, err := h.svc.Update(c.Context(), existing)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *WebhookHandler) Delete(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid webhook id")
	}

	if err := h.svc.Delete(c.Context(), accountID, id); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *WebhookHandler) WhatsApp(c fiber.Ctx) error {
	var payload map[string]any
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	entry, ok := payload["entry"].([]any)
	if !ok || len(entry) == 0 {
		return c.JSON(fiber.Map{"status": "ok"})
	}

	entryMap, ok := entry[0].(map[string]any)
	if !ok {
		return c.JSON(fiber.Map{"status": "ok"})
	}

	changes, ok := entryMap["changes"].([]any)
	if !ok || len(changes) == 0 {
		return c.JSON(fiber.Map{"status": "ok"})
	}

	return c.JSON(fiber.Map{"status": "received"})
}

func (h *WebhookHandler) Telegram(c fiber.Ctx) error {
	var payload map[string]any
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	message, ok := payload["message"].(map[string]any)
	if !ok {
		return c.JSON(fiber.Map{"status": "ok"})
	}

	chat, _ := message["chat"].(map[string]any)
	chatID, _ := chat["id"].(float64)

	text, _ := message["text"].(string)
	if text == "" {
		text, _ = message["caption"].(string)
	}

	return c.JSON(fiber.Map{
		"status":  "received",
		"chat_id": chatID,
		"text":    text,
	})
}

func (h *WebhookHandler) Email(c fiber.Ctx) error {
	var payload map[string]any
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	return c.JSON(fiber.Map{"status": "received"})
}

func (h *WebhookHandler) Debug(c fiber.Ctx) error {
	var payload map[string]any
	c.Bind().JSON(&payload)
	jsonBytes, _ := json.MarshalIndent(payload, "", "  ")
	return c.SendString(string(jsonBytes))
}
