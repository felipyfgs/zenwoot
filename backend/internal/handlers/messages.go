package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type MessageHandler struct {
	svc *services.MessageService
}

func NewMessageHandler(svc *services.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) Register(rg fiber.Router) {
	rg.Get("/conversations/:conversation_id/messages", h.List)
	rg.Post("/conversations/:conversation_id/messages", h.Create)
	rg.Get("/messages/:id", h.Get)
	rg.Patch("/messages/:id", h.Update)
	rg.Delete("/messages/:id", h.Delete)
}

func (h *MessageHandler) List(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	convID, _ := strconv.ParseInt(c.Params("conversation_id"), 10, 64)
	limit, _ := strconv.Atoi(c.Query("limit", "25"))

	var before *int64
	if b := c.Query("before"); b != "" {
		v, _ := strconv.ParseInt(b, 10, 64)
		before = &v
	}

	msgs, err := h.svc.List(c.Context(), accountID, convID, before, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": msgs})
}

func (h *MessageHandler) Create(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	userID := c.Locals("user_id").(int64)
	convID, _ := strconv.ParseInt(c.Params("conversation_id"), 10, 64)

	var body struct {
		Content      string         `json:"content"`
		Private      bool           `json:"private"`
		MessageType  int            `json:"message_type"`
		ContentAttrs map[string]any `json:"content_attributes"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	msg, err := h.svc.Create(c.Context(), services.CreateMessageInput{
		ConversationID: convID,
		AccountID:      accountID,
		SenderType:     "User",
		SenderID:       userID,
		Content:        body.Content,
		MessageType:    body.MessageType,
		Private:        body.Private,
		ContentAttrs:   body.ContentAttrs,
	})
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(msg)
}

func (h *MessageHandler) Get(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	msg, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "message not found"})
	}
	return c.JSON(msg)
}

func (h *MessageHandler) Update(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)

	existing, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "message not found"})
	}

	var body models.Message
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	existing.Content = body.Content
	existing.ContentAttributes = body.ContentAttributes

	updated, err := h.svc.Update(c.Context(), existing)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updated)
}

func (h *MessageHandler) Delete(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	if err := h.svc.Delete(c.Context(), accountID, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
