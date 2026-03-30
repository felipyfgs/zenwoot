package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"wzap/internal/domain"
	"wzap/internal/dto"
	"wzap/internal/service"
)

type AssignmentHandler struct {
	assignSvc *service.AssignmentService
}

func NewAssignmentHandler(assignSvc *service.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{assignSvc: assignSvc}
}

// Assign godoc
// @Summary     Assign conversation
// @Description Assigns a conversation to a user or team
// @Tags        Conversations
// @Accept      json
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Param       body body dto.AssignReq true "Assignment data"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/assign [post]
func (h *AssignmentHandler) Assign(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")

	var req dto.AssignReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	var conv *domain.Conversation
	var err error

	if req.UserID != "" {
		conv, err = h.assignSvc.AssignToUser(c.Context(), conversationID, req.UserID)
	} else if req.TeamID != "" {
		conv, err = h.assignSvc.AssignToTeam(c.Context(), conversationID, req.TeamID)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", "userId or teamId required"))
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(conv))
}

// Unassign godoc
// @Summary     Unassign conversation
// @Description Removes assignment from conversation
// @Tags        Conversations
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/unassign [post]
func (h *AssignmentHandler) Unassign(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")

	conv, err := h.assignSvc.Unassign(c.Context(), conversationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(conv))
}

// UpdatePriority godoc
// @Summary     Update priority
// @Description Updates conversation priority
// @Tags        Conversations
// @Accept      json
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Param       body body dto.PriorityReq true "Priority data"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/priority [post]
func (h *AssignmentHandler) UpdatePriority(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")

	var req dto.PriorityReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	conv, err := h.assignSvc.UpdatePriority(c.Context(), conversationID, domain.ConversationPriority(req.Priority))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(conv))
}

// Snooze godoc
// @Summary     Snooze conversation
// @Description Snoozes conversation until specified time
// @Tags        Conversations
// @Accept      json
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Param       body body dto.SnoozeReq true "Snooze data"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/snooze [post]
func (h *AssignmentHandler) Snooze(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")

	var req dto.SnoozeReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	until, err := time.Parse(time.RFC3339, req.Until)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", "invalid timestamp format"))
	}

	conv, err := h.assignSvc.Snooze(c.Context(), conversationID, until)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(conv))
}

// UnSnooze godoc
// @Summary     Unsnooze conversation
// @Description Reopens a snoozed conversation
// @Tags        Conversations
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/unsnooze [post]
func (h *AssignmentHandler) UnSnooze(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")

	conv, err := h.assignSvc.UnSnooze(c.Context(), conversationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(conv))
}
