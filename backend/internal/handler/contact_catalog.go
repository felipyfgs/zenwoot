package handler

import (
	"github.com/gofiber/fiber/v2"

	"wzap/internal/dto"
	"wzap/internal/service"
)

// ContactCatalogHandler serves the persisted contacts catalog (wzContacts).
// Not to be confused with ContactHandler which operates on live WhatsApp state.
type ContactCatalogHandler struct {
	svc *service.ContactCatalogService
}

func NewContactCatalogHandler(svc *service.ContactCatalogService) *ContactCatalogHandler {
	return &ContactCatalogHandler{svc: svc}
}

// List godoc
// @Summary     List all contacts
// @Description Returns all persisted contacts, paginated
// @Tags        Contacts
// @Produce     json
// @Param       page  query int false "Page"
// @Param       limit query int false "Limit"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/contacts [get]
func (h *ContactCatalogHandler) List(c *fiber.Ctx) error {
	var req dto.PaginationReq
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	// Single-tenant: empty accountID means "all"
	list, err := h.svc.List(c.Context(), "", req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}

	return c.JSON(dto.SuccessResp(list))
}

// Get godoc
// @Summary     Get a contact by ID
// @Tags        Contacts
// @Produce     json
// @Param       contactId path string true "Contact ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/contacts/{contactId} [get]
func (h *ContactCatalogHandler) Get(c *fiber.Ctx) error {
	id := c.Params("contactId")
	contact, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", "Contact not found"))
	}
	return c.JSON(dto.SuccessResp(contact))
}

// ListConversations godoc
// @Summary     List conversations for a contact
// @Tags        Contacts
// @Produce     json
// @Param       contactId path string true  "Contact ID"
// @Param       page      query int    false "Page"
// @Param       limit     query int    false "Limit"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/contacts/{contactId}/conversations [get]
func (h *ContactCatalogHandler) ListConversations(c *fiber.Ctx) error {
	contactID := c.Params("contactId")
	var req dto.PaginationReq
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	list, err := h.svc.ListConversations(c.Context(), contactID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}

	return c.JSON(dto.SuccessResp(list))
}
