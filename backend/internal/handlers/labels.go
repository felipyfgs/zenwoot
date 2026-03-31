package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type LabelHandler struct {
	svc *services.LabelService
}

func NewLabelHandler(svc *services.LabelService) *LabelHandler {
	return &LabelHandler{svc: svc}
}

func (h *LabelHandler) List(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)

	items, err := h.svc.List(c.Context(), accountID)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *LabelHandler) Get(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid label id")
	}

	item, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "label not found")
	}
	return c.JSON(item)
}

func (h *LabelHandler) Create(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	var body models.Label
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

func (h *LabelHandler) Update(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid label id")
	}

	existing, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "label not found")
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

func (h *LabelHandler) Delete(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid label id")
	}

	if err := h.svc.Delete(c.Context(), accountID, id); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *LabelHandler) AddToConversation(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	conversationID, err := helpers.ParseID(c, "conversation_id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	var body struct {
		LabelID int64 `json:"label_id"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	if err := h.svc.AddToConversation(c.Context(), accountID, conversationID, body.LabelID); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *LabelHandler) RemoveFromConversation(c fiber.Ctx) error {
	conversationID, err := helpers.ParseID(c, "conversation_id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	labelID, err := helpers.ParseID(c, "label_id")
	if err != nil {
		return helpers.BadRequest(c, "invalid label id")
	}

	if err := h.svc.RemoveFromConversation(c.Context(), conversationID, labelID); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
