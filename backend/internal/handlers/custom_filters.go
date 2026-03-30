package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type CustomFilterHandler struct {
	repo *repo.CustomFilterRepo
}

func NewCustomFilterHandler(repo *repo.CustomFilterRepo) *CustomFilterHandler {
	return &CustomFilterHandler{repo: repo}
}

func (h *CustomFilterHandler) Register(rg fiber.Router) {
	rg.Get("/custom_filters", h.List)
	rg.Post("/custom_filters", h.Create)
	rg.Put("/custom_filters/:id", h.Update)
	rg.Delete("/custom_filters/:id", h.Delete)
}

func (h *CustomFilterHandler) List(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	userID := c.Locals("user_id").(int64)

	filters, err := h.repo.List(c.Context(), accountID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": filters})
}

func (h *CustomFilterHandler) Create(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	userID := c.Locals("user_id").(int64)

	var body models.CustomFilter
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	body.AccountID = accountID
	body.UserID = userID

	err := h.repo.Create(c.Context(), &body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(body)
}

func (h *CustomFilterHandler) Update(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	userID := c.Locals("user_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)

	var body models.CustomFilter
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	body.ID = id
	body.AccountID = accountID
	body.UserID = userID

	err := h.repo.Update(c.Context(), &body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(body)
}

func (h *CustomFilterHandler) Delete(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	userID := c.Locals("user_id").(int64)
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)

	err := h.repo.Delete(c.Context(), id, accountID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
