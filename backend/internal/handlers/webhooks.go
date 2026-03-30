package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

type WebhookHandler struct{}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

func (h *WebhookHandler) Register(rg fiber.Router) {
	rg.Post("/webhooks/whatsapp", h.WhatsApp)
	rg.Post("/webhooks/telegram", h.Telegram)
	rg.Post("/webhooks/email", h.Email)
	rg.Post("/webhooks/debug", h.Debug)
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
