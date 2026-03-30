package handler

import (
	"wzap/internal/dto"
	"wzap/internal/service"

	"github.com/gofiber/fiber/v2"
)

type WebhookHandler struct {
	webhookSvc *service.WebhookService
}

func NewWebhookHandler(webhookSvc *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookSvc: webhookSvc}
}

// Create godoc
// @Summary     Create a webhook
// @Description Registers a new webhook for the inbox
// @Tags        Webhooks
// @Accept      json
// @Produce     json
// @Param       inboxId   path     string                 true "Inbox ID"
// @Param       body        body     dto.CreateWebhookReq true "Webhook data"
// @Success     200  {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/webhooks [post]
func (h *WebhookHandler) Create(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	var req dto.CreateWebhookReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	webhook, err := h.webhookSvc.Create(c.Context(), inboxID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResp(webhook))
}

// List godoc
// @Summary     List webhooks
// @Description Returns all webhooks for the inbox
// @Tags        Webhooks
// @Produce     json
// @Param       inboxId   path string true "Inbox ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/webhooks [get]
func (h *WebhookHandler) List(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	webhooks, err := h.webhookSvc.List(c.Context(), inboxID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}

	return c.JSON(dto.SuccessResp(webhooks))
}

// Delete godoc
// @Summary     Delete a webhook
// @Description Removes a webhook from the inbox
// @Tags        Webhooks
// @Produce     json
// @Param       inboxId   path string true "Inbox ID"
// @Param       wid         path string true "Webhook ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/webhooks/{wid} [delete]
func (h *WebhookHandler) Delete(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	webhookID := c.Params("wid")

	if err := h.webhookSvc.Delete(c.Context(), inboxID, webhookID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}

	return c.JSON(dto.SuccessResp(nil))
}
