package handler

import (
	"github.com/gofiber/fiber/v2"

	"wzap/internal/dto"
	"wzap/internal/service"
)

type LabelHandler struct {
	labelSvc *service.LabelService
}

func NewLabelHandler(labelSvc *service.LabelService) *LabelHandler {
	return &LabelHandler{labelSvc: labelSvc}
}

// Create godoc
// @Summary     Create a new label
// @Description Creates a new label/tag for the account
// @Tags        Labels
// @Accept      json
// @Produce     json
// @Param       body body     dto.LabelCreateReq true "Label data"
// @Success     200  {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/labels [post]
func (h *LabelHandler) Create(c *fiber.Ctx) error {
	accountID := getAccountID(c)

	var req dto.LabelCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	label, err := h.labelSvc.Create(c.Context(), accountID, service.CreateLabelReq{
		Title:       req.Title,
		Color:       req.Color,
		Description: req.Description,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResp(label))
}

// List godoc
// @Summary     List labels
// @Description Returns all labels for the account
// @Tags        Labels
// @Produce     json
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/labels [get]
func (h *LabelHandler) List(c *fiber.Ctx) error {
	accountID := getAccountID(c)

	labels, err := h.labelSvc.ListByAccount(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(labels))
}

// Get godoc
// @Summary     Get label
// @Description Returns label by ID
// @Tags        Labels
// @Produce     json
// @Param       labelId path string true "Label ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/labels/{labelId} [get]
func (h *LabelHandler) Get(c *fiber.Ctx) error {
	labelID := c.Params("labelId")
	label, err := h.labelSvc.Get(c.Context(), labelID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", err.Error()))
	}
	return c.JSON(dto.SuccessResp(label))
}

// Update godoc
// @Summary     Update label
// @Description Updates label data
// @Tags        Labels
// @Accept      json
// @Produce     json
// @Param       labelId path string true "Label ID"
// @Param       body    body  dto.LabelUpdateReq true "Label data"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/labels/{labelId} [put]
func (h *LabelHandler) Update(c *fiber.Ctx) error {
	labelID := c.Params("labelId")

	var req dto.LabelUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	label, err := h.labelSvc.Update(c.Context(), labelID, service.UpdateLabelReq{
		Title:       req.Title,
		Color:       req.Color,
		Description: req.Description,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(label))
}

// Delete godoc
// @Summary     Delete label
// @Description Removes label from account
// @Tags        Labels
// @Produce     json
// @Param       labelId path string true "Label ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/labels/{labelId} [delete]
func (h *LabelHandler) Delete(c *fiber.Ctx) error {
	labelID := c.Params("labelId")
	if err := h.labelSvc.Delete(c.Context(), labelID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// AddToConversation godoc
// @Summary     Add label to conversation
// @Description Adds a label to a conversation
// @Tags        Labels
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Param       labelId        path string true "Label ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/labels/{labelId} [post]
func (h *LabelHandler) AddToConversation(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")
	labelID := c.Params("labelId")

	if err := h.labelSvc.AddToConversation(c.Context(), conversationID, labelID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// RemoveFromConversation godoc
// @Summary     Remove label from conversation
// @Description Removes a label from a conversation
// @Tags        Labels
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Param       labelId        path string true "Label ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/labels/{labelId} [delete]
func (h *LabelHandler) RemoveFromConversation(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")
	labelID := c.Params("labelId")

	if err := h.labelSvc.RemoveFromConversation(c.Context(), conversationID, labelID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// GetConversationLabels godoc
// @Summary     Get conversation labels
// @Description Returns all labels for a conversation
// @Tags        Labels
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/labels [get]
func (h *LabelHandler) GetConversationLabels(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")

	labels, err := h.labelSvc.GetConversationLabels(c.Context(), conversationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(labels))
}

// SetConversationLabels godoc
// @Summary     Set conversation labels
// @Description Replaces all labels for a conversation
// @Tags        Labels
// @Accept      json
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Param       body           body  dto.SetLabelsReq true "Label IDs"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/labels [put]
func (h *LabelHandler) SetConversationLabels(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")

	var req dto.SetLabelsReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	if err := h.labelSvc.SetConversationLabels(c.Context(), conversationID, req.LabelIDs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}
