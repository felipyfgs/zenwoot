package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type ConversationHandler struct {
	svc *services.ConversationService
}

func NewConversationHandler(svc *services.ConversationService) *ConversationHandler {
	return &ConversationHandler{svc: svc}
}

func (h *ConversationHandler) List(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	page := helpers.ParseQueryInt(c, "page", 1)
	pageSize := helpers.ParseQueryInt(c, "pageSize", 25)

	f := repo.ConversationFilter{Page: page, PageSize: pageSize}
	if s := c.Query("status"); s != "" {
		v := helpers.ParseQueryInt(c, "status", 0)
		f.Status = &v
	}
	if s := c.Query("assignee_id"); s != "" {
		v, _ := helpers.ParseID(c, "assignee_id")
		f.AssigneeID = &v
	}
	if s := c.Query("inbox_id"); s != "" {
		v, _ := helpers.ParseID(c, "inbox_id")
		f.InboxID = &v
	}
	if s := c.Query("team_id"); s != "" {
		v, _ := helpers.ParseID(c, "team_id")
		f.TeamID = &v
	}

	items, total, err := h.svc.List(c.Context(), accountID, f)
	if err != nil {
		return helpers.InternalError(c, err)
	}
	return c.JSON(fiber.Map{"data": items, "total": total, "page": page, "pageSize": pageSize})
}

func (h *ConversationHandler) Get(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	conv, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return helpers.NotFound(c, "conversation not found")
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) Create(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	var body models.Conversation
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}
	body.AccountID = accountID
	conv, err := h.svc.Create(c.Context(), &body)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return helpers.Created(c, conv)
}

func (h *ConversationHandler) Resolve(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	conv, err := h.svc.Resolve(c.Context(), accountID, id)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) Reopen(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	conv, err := h.svc.Reopen(c.Context(), accountID, id)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) Snooze(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	var body struct {
		SnoozedUntil time.Time `json:"snoozed_until"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	conv, err := h.svc.Snooze(c.Context(), accountID, id, body.SnoozedUntil)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) Assign(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	var body struct {
		AssigneeID *int64 `json:"assignee_id"`
		TeamID     *int64 `json:"team_id"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	conv, err := h.svc.Assign(c.Context(), accountID, id, body.AssigneeID, body.TeamID)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) Delete(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	id, err := helpers.ParseID(c, "id")
	if err != nil {
		return helpers.BadRequest(c, "invalid conversation id")
	}

	if err := h.svc.Delete(c.Context(), accountID, id); err != nil {
		return helpers.InternalError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
