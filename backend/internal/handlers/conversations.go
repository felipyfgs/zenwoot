package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

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

func (h *ConversationHandler) Register(rg fiber.Router) {
	rg.Get("/conversations", h.List)
	rg.Post("/conversations", h.Create)
	rg.Get("/conversations/:id", h.Get)
	rg.Post("/conversations/:id/resolve", h.Resolve)
	rg.Post("/conversations/:id/reopen", h.Reopen)
	rg.Post("/conversations/:id/snooze", h.Snooze)
	rg.Post("/conversations/:id/assignments", h.Assign)
}

func (h *ConversationHandler) List(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "25"))

	f := repo.ConversationFilter{Page: page, PageSize: pageSize}
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		f.Status = &v
	}
	if s := c.Query("assignee_id"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 64)
		f.AssigneeID = &v
	}
	if s := c.Query("inbox_id"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 64)
		f.InboxID = &v
	}
	if s := c.Query("team_id"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 64)
		f.TeamID = &v
	}

	items, total, err := h.svc.List(c.Context(), accountID, f)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": items, "total": total, "page": page, "pageSize": pageSize})
}

func (h *ConversationHandler) Get(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	conv, err := h.svc.GetByID(c.Context(), accountID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) Create(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	var body models.Conversation
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	body.AccountID = accountID
	conv, err := h.svc.Create(c.Context(), &body)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(conv)
}

func (h *ConversationHandler) Resolve(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	conv, err := h.svc.Resolve(c.Context(), accountID, id)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) Reopen(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	conv, err := h.svc.Reopen(c.Context(), accountID, id)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) Snooze(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	var body struct {
		SnoozedUntil time.Time `json:"snoozed_until"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	conv, err := h.svc.Snooze(c.Context(), accountID, id, body.SnoozedUntil)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) Assign(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	var body struct {
		AssigneeID *int64 `json:"assignee_id"`
		TeamID     *int64 `json:"team_id"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	conv, err := h.svc.Assign(c.Context(), accountID, id, body.AssigneeID, body.TeamID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(conv)
}
