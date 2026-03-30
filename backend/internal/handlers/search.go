package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type SearchHandler struct {
	contactRepo *repo.ContactRepo
	convRepo    *repo.ConversationRepo
}

func NewSearchHandler(contactRepo *repo.ContactRepo, convRepo *repo.ConversationRepo) *SearchHandler {
	return &SearchHandler{contactRepo: contactRepo, convRepo: convRepo}
}

func (h *SearchHandler) Register(rg fiber.Router) {
	rg.Get("/search", h.Search)
}

func (h *SearchHandler) Search(c fiber.Ctx) error {
	accountID := c.Locals("account_id").(int64)
	q := c.Query("q", "")
	if q == "" {
		return c.JSON(fiber.Map{"contacts": []any{}, "conversations": []any{}})
	}
	contacts, _, _ := h.contactRepo.Search(c.Context(), accountID, q, 1, 10)
	return c.JSON(fiber.Map{"contacts": contacts})
}
