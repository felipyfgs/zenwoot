package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type NoteHandler struct {
	svc *services.NoteService
}

func NewNoteHandler(svc *services.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h *NoteHandler) List(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	contactID, err := helpers.ParseID(c, "contact_id")
	if err != nil {
		return helpers.BadRequest(c, "invalid contact id")
	}

	items, err := h.svc.ListByContact(c.Context(), accountID, contactID)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *NoteHandler) Get(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid note id")
	}

	item, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "note not found")
	}
	return c.JSON(item)
}

func (h *NoteHandler) Create(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	contactID, err := helpers.ParseID(c, "contact_id")
	if err != nil {
		return helpers.BadRequest(c, "invalid contact id")
	}
	userID := helpers.GetUserID(c)

	var body struct {
		Content string `json:"content"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	item := &models.Note{
		AccountID: accountID,
		ContactID: contactID,
		UserID:    userID,
		Content:   body.Content,
	}

	created, err := h.svc.Create(c.Context(), item)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return helpers.Created(c, created)
}

func (h *NoteHandler) Update(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid note id")
	}

	existing, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "note not found")
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

func (h *NoteHandler) Delete(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid note id")
	}

	if err := h.svc.Delete(c.Context(), accountID, id); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
