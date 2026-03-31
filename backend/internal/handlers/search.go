package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
)

type SearchHandler struct {
	contactRepo *repo.ContactRepo
	convRepo    *repo.ConversationRepo
}

func NewSearchHandler(contactRepo *repo.ContactRepo, convRepo *repo.ConversationRepo) *SearchHandler {
	return &SearchHandler{contactRepo: contactRepo, convRepo: convRepo}
}

func (h *SearchHandler) Search(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)
	q := c.Query("q", "")

	if q == "" {
		return c.JSON(fiber.Map{"contacts": []any{}, "conversations": []any{}})
	}

	contacts, _, err := h.contactRepo.Search(c.Context(), accountID, q, 1, 10)
	if err != nil {
		return helpers.InternalError(c, err)
	}

	conversations, _, err := h.convRepo.Search(c.Context(), accountID, q, 1, 10)
	if err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"contacts":      contacts,
		"conversations": conversations,
	})
}
