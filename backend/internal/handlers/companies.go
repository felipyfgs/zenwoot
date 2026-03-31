package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type CompanyHandler struct {
	svc *services.CompanyService
}

func NewCompanyHandler(svc *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

func (h *CompanyHandler) List(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)

	items, err := h.svc.List(c.Context(), accountID)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *CompanyHandler) Search(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	q := c.Query("q", "")

	items, err := h.svc.Search(c.Context(), accountID, q)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *CompanyHandler) Get(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid company id")
	}

	item, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "company not found")
	}
	return c.JSON(item)
}

func (h *CompanyHandler) Create(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	var body models.Company
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

func (h *CompanyHandler) Update(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid company id")
	}

	existing, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "company not found")
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

func (h *CompanyHandler) Delete(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid company id")
	}

	if err := h.svc.Delete(c.Context(), accountID, id); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
