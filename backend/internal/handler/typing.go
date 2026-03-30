package handler

import (
	"github.com/gofiber/fiber/v2"

	"wzap/internal/dto"
	"wzap/internal/ws"
)

type TypingHandler struct {
	wsHub     *ws.Hub
	accountID string
}

func NewTypingHandler(wsHub *ws.Hub, accountID string) *TypingHandler {
	return &TypingHandler{wsHub: wsHub, accountID: accountID}
}

// TypingReq represents a typing indicator request
type TypingReq struct {
	ConversationID string `json:"conversationId"`
	InboxID        string `json:"inboxId"`
	UserID         string `json:"userId"`
	UserName       string `json:"userName,omitempty"`
	IsTyping       bool   `json:"isTyping"`
}

// SetTyping godoc
// @Summary     Set typing indicator
// @Description Broadcasts typing status to other clients
// @Tags        Conversations
// @Accept      json
// @Produce     json
// @Param       body body TypingReq true "Typing data"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/typing [post]
func (h *TypingHandler) SetTyping(c *fiber.Ctx) error {
	var req TypingReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	// Get user info from context if available
	userID := req.UserID
	if uid, ok := c.Locals("userId").(string); ok && uid != "" {
		userID = uid
	}
	userName := req.UserName
	if name, ok := c.Locals("userName").(string); ok && name != "" {
		userName = name
	}

	// Broadcast typing event via WebSocket
	h.wsHub.PublishTyping(h.accountID, req.InboxID, req.ConversationID, userID, userName, req.IsTyping)

	return c.JSON(dto.SuccessResp(nil))
}
