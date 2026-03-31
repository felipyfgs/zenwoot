package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type InboxHandler struct {
	repo *repo.InboxRepo
}

func NewInboxHandler(r *repo.InboxRepo) *InboxHandler {
	return &InboxHandler{repo: r}
}

func (h *InboxHandler) List(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	items, err := h.repo.ListByAccount(c.Context(), accountID)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *InboxHandler) Get(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid inbox id")
	}

	inbox, err := h.repo.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "inbox not found")
	}
	return c.JSON(inbox)
}

func (h *InboxHandler) Create(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	var body models.Inbox
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	body.AccountID = accountID
	if err := h.repo.Create(c.Context(), &body); err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return helpers.Created(c, body)
}

func (h *InboxHandler) Update(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid inbox id")
	}

	existing, err := h.repo.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "inbox not found")
	}

	if err := c.Bind().JSON(existing); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	if err := h.repo.Update(c.Context(), existing); err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return c.JSON(existing)
}

func (h *InboxHandler) Delete(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid inbox id")
	}

	if err := h.repo.Delete(c.Context(), accountID, id); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
