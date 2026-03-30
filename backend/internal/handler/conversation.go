package handler

import (
	"github.com/gofiber/fiber/v2"

	"wzap/internal/dto"
	"wzap/internal/service"
)

type ConversationHandler struct {
	svc *service.ConversationService
}

func NewConversationHandler(svc *service.ConversationService) *ConversationHandler {
	return &ConversationHandler{svc: svc}
}

// ListByInbox godoc
// @Summary     List conversations for an inbox
// @Description Returns paginated conversations for a specific inbox
// @Tags        Conversations
// @Produce     json
// @Param       inboxId path     string true  "Inbox ID"
// @Param       page    query    int    false "Page number (default 1)"
// @Param       limit   query    int    false "Items per page (default 50)"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/conversations [get]
func (h *ConversationHandler) ListByInbox(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	var req dto.PaginationReq
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	list, err := h.svc.ListByInbox(c.Context(), inboxID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}

	return c.JSON(dto.SuccessResp(list))
}

// Get godoc
// @Summary     Get a conversation
// @Tags        Conversations
// @Produce     json
// @Param       conversationId path     string true "Conversation ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId} [get]
func (h *ConversationHandler) Get(c *fiber.Ctx) error {
	id := c.Params("conversationId")
	conv, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", "Conversation not found"))
	}
	return c.JSON(dto.SuccessResp(conv))
}

// ToggleStatus godoc
// @Summary     Toggle conversation status (open ↔ resolved)
// @Tags        Conversations
// @Produce     json
// @Param       conversationId path     string true "Conversation ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/toggle_status [post]
func (h *ConversationHandler) ToggleStatus(c *fiber.Ctx) error {
	id := c.Params("conversationId")
	conv, err := h.svc.ToggleStatus(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", err.Error()))
	}
	return c.JSON(dto.SuccessResp(conv))
}

// MarkRead godoc
// @Summary     Mark conversation as read
// @Tags        Conversations
// @Produce     json
// @Param       conversationId path     string true "Conversation ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/read [post]
func (h *ConversationHandler) MarkRead(c *fiber.Ctx) error {
	id := c.Params("conversationId")
	if err := h.svc.MarkRead(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// ListMessages godoc
// @Summary     List messages in a conversation
// @Tags        Conversations
// @Produce     json
// @Param       conversationId path     string true  "Conversation ID"
// @Param       page           query    int    false "Page"
// @Param       limit          query    int    false "Limit"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/messages [get]
func (h *ConversationHandler) ListMessages(c *fiber.Ctx) error {
	id := c.Params("conversationId")
	var req dto.PaginationReq
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	msgs, err := h.svc.ListMessages(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}

	return c.JSON(dto.SuccessResp(msgs))
}
