package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type TeamHandler struct {
	svc *services.TeamService
}

func NewTeamHandler(svc *services.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) List(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)

	items, err := h.svc.List(c.Context(), accountID)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *TeamHandler) Get(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid team id")
	}

	item, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "team not found")
	}
	return c.JSON(item)
}

func (h *TeamHandler) Create(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	var body models.Team
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

func (h *TeamHandler) Update(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid team id")
	}

	existing, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "team not found")
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

func (h *TeamHandler) Delete(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid team id")
	}

	if err := h.svc.Delete(c.Context(), accountID, id); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TeamHandler) ListMembers(c fiber.Ctx) error {
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid team id")
	}

	members, err := h.svc.ListMembers(c.Context(), id)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": members})
}

func (h *TeamHandler) AddMember(c fiber.Ctx) error {
	teamID, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid team id")
	}

	var body struct {
		UserID int64 `json:"user_id"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	if err := h.svc.AddMember(c.Context(), teamID, body.UserID); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *TeamHandler) RemoveMember(c fiber.Ctx) error {
	teamID, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid team id")
	}

	userID, err := helpers.ParseID(c, "user_id")
	if err != nil {
		return helpers.BadRequest(c, "invalid user id")
	}

	if err := h.svc.RemoveMember(c.Context(), teamID, userID); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
