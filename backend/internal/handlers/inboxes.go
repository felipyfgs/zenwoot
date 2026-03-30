package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type InboxHandler struct {
	repo *repo.InboxRepo
}

func NewInboxHandler(r *repo.InboxRepo) *InboxHandler {
	return &InboxHandler{repo: r}
}

func (h *InboxHandler) Register(rg fiber.Router) {
	rg.Get("/inboxes", h.List)
	rg.Post("/inboxes", h.Create)
	rg.Get("/inboxes/:id", h.Get)
	rg.Patch("/inboxes/:id", h.Update)
	rg.Delete("/inboxes/:id", h.Delete)
}

func (h *InboxHandler) List(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	items, err := h.repo.ListByAccount(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *InboxHandler) Get(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	inbox, err := h.repo.GetByID(c.Context(), accountID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(inbox)
}

func (h *InboxHandler) Create(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	var body models.Inbox
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	body.AccountID = accountID
	if err := h.repo.Create(c.Context(), &body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(body)
}

func (h *InboxHandler) Update(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	existing, err := h.repo.GetByID(c.Context(), accountID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	if err := c.Bind().JSON(existing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if err := h.repo.Update(c.Context(), existing); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(existing)
}

func (h *InboxHandler) Delete(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	if err := h.repo.Delete(c.Context(), accountID, id); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
