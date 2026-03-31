package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type CustomFilterHandler struct {
	repo *repo.CustomFilterRepo
}

func NewCustomFilterHandler(repo *repo.CustomFilterRepo) *CustomFilterHandler {
	return &CustomFilterHandler{repo: repo}
}

func (h *CustomFilterHandler) List(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	userID := helpers.GetUserID(c)

	filters, err := h.repo.List(c.Context(), accountID, userID)
	if err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(fiber.Map{"data": filters})
}

func (h *CustomFilterHandler) Create(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	userID := helpers.GetUserID(c)

	var body models.CustomFilter
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	body.AccountID = accountID
	body.UserID = userID

	if err := h.repo.Create(c.Context(), &body); err != nil {
		return helpers.InternalError(c, err)
	}

	return helpers.Created(c, body)
}

func (h *CustomFilterHandler) Update(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	userID := helpers.GetUserID(c)

	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid filter id")
	}

	var body models.CustomFilter
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	body.ID = id
	body.AccountID = accountID
	body.UserID = userID

	if err := h.repo.Update(c.Context(), &body); err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(body)
}

func (h *CustomFilterHandler) Delete(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	userID := helpers.GetUserID(c)

	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid filter id")
	}

	if err := h.repo.Delete(c.Context(), id, accountID, userID); err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(fiber.Map{"success": true})
}
