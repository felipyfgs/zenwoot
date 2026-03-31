package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type MessageHandler struct {
	svc *services.MessageService
}

func NewMessageHandler(svc *services.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) List(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	convID, err := helpers.ParseID(c, "conversation_id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	limit := helpers.ParseQueryInt(c, "limit", 25)
	var before *int64
	if b := c.Query("before"); b != "" {
		v, _ := helpers.ParseID(c, "before")
		before = &v
	}

	msgs, err := h.svc.List(c.Context(), accountID, convID, before, limit)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": msgs})
}

func (h *MessageHandler) Create(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	userID := helpers.GetUserID(c)
	convID, err := helpers.ParseID(c, "conversation_id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	var body struct {
		Content      string         `json:"content"`
		Private      bool           `json:"private"`
		MessageType  int            `json:"message_type"`
		ContentAttrs map[string]any `json:"content_attributes"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
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
		return helpers.Unprocessable(c, err.Error())
	}
	return helpers.Created(c, msg)
}

func (h *MessageHandler) Get(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid message id")
	}

	msg, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "message not found")
	}
	return c.JSON(msg)
}

func (h *MessageHandler) Update(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid message id")
	}

	existing, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "message not found")
	}

	var body models.Message
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	existing.Content = body.Content
	existing.ContentAttributes = body.ContentAttributes

	updated, err := h.svc.Update(c.Context(), existing)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *MessageHandler) Delete(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid message id")
	}

	if err := h.svc.Delete(c.Context(), accountID, id); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
