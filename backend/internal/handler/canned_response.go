package handler

import (
	"github.com/gofiber/fiber/v2"

	"wzap/internal/dto"
	"wzap/internal/service"
)

type CannedResponseHandler struct {
	svc *service.CannedResponseService
}

func NewCannedResponseHandler(svc *service.CannedResponseService) *CannedResponseHandler {
	return &CannedResponseHandler{svc: svc}
}

// Create godoc
// @Summary     Create canned response
// @Description Creates a new preset reply shortcut
// @Tags        CannedResponses
// @Accept      json
// @Produce     json
// @Param       body body dto.CannedResponseCreateReq true "Canned response data"
// @Success     201 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/canned-responses [post]
func (h *CannedResponseHandler) Create(c *fiber.Ctx) error {
	accountID := getAccountID(c)

	var req dto.CannedResponseCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	cr, err := h.svc.Create(c.Context(), accountID, service.CreateCannedResponseReq{
		ShortCode: req.ShortCode,
		Content:   req.Content,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResp(cr))
}

// List godoc
// @Summary     List canned responses
// @Description Returns canned responses, optionally filtered by search query
// @Tags        CannedResponses
// @Produce     json
// @Param       search query string false "Search by short code or content"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/canned-responses [get]
func (h *CannedResponseHandler) List(c *fiber.Ctx) error {
	accountID := getAccountID(c)
	search := c.Query("search")

	var list interface{}
	var err error
	if search != "" {
		list, err = h.svc.Search(c.Context(), accountID, search)
	} else {
		list, err = h.svc.List(c.Context(), accountID)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(list))
}

// Update godoc
// @Summary     Update canned response
// @Description Updates a canned response
// @Tags        CannedResponses
// @Accept      json
// @Produce     json
// @Param       id   path string true "Canned Response ID"
// @Param       body body dto.CannedResponseUpdateReq true "Canned response data"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/canned-responses/{id} [put]
func (h *CannedResponseHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.CannedResponseUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	cr, err := h.svc.Update(c.Context(), id, service.UpdateCannedResponseReq{
		ShortCode: req.ShortCode,
		Content:   req.Content,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(cr))
}

// Delete godoc
// @Summary     Delete canned response
// @Description Removes a canned response
// @Tags        CannedResponses
// @Produce     json
// @Param       id path string true "Canned Response ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/canned-responses/{id} [delete]
func (h *CannedResponseHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}
