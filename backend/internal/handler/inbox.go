package handler

import (
	"encoding/base64"

	"wzap/internal/dto"
	"wzap/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

type InboxHandler struct {
	inboxSvc *service.InboxService
}

func NewInboxHandler(inboxSvc *service.InboxService) *InboxHandler {
	return &InboxHandler{inboxSvc: inboxSvc}
}

// Create godoc
// @Summary     Create a new inbox
// @Description Creates a new WhatsApp inbox
// @Tags        Inboxes
// @Accept      json
// @Produce     json
// @Param       body body     dto.InboxCreateReq true "Inbox data"
// @Success     200  {object} dto.APIResponse
// @Failure     400  {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes [post]
func (h *InboxHandler) Create(c *fiber.Ctx) error {
	if c.Locals("authRole") != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResp("Forbidden", "Admin access required"))
	}

	var req dto.InboxCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	inbox, err := h.inboxSvc.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResp(inbox))
}

// List godoc
// @Summary     List inboxes
// @Description Returns all inboxes for the account
// @Tags        Inboxes
// @Produce     json
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes [get]
func (h *InboxHandler) List(c *fiber.Ctx) error {
	inboxes, err := h.inboxSvc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(inboxes))
}

// Get godoc
// @Summary     Get inbox
// @Description Returns the inbox identified by :inboxId
// @Tags        Inboxes
// @Produce     json
// @Param       inboxId   path string true "Inbox ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId} [get]
func (h *InboxHandler) Get(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	inbox, err := h.inboxSvc.Get(c.Context(), inboxID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", err.Error()))
	}
	return c.JSON(dto.SuccessResp(inbox))
}

// Delete godoc
// @Summary     Delete inbox
// @Description Disconnects and deletes the inbox
// @Tags        Inboxes
// @Produce     json
// @Param       inboxId   path string true "Inbox ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId} [delete]
func (h *InboxHandler) Delete(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	if err := h.inboxSvc.Delete(c.Context(), inboxID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// Connect godoc
// @Summary     Connect inbox
// @Description Connects a WhatsApp inbox (starts pairing if new)
// @Tags        Inboxes
// @Produce     json
// @Param       inboxId   path string true "Inbox ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/connect [post]
func (h *InboxHandler) Connect(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	status, err := h.inboxSvc.Connect(c.Context(), inboxID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(map[string]string{"status": status}))
}

// Disconnect godoc
// @Summary     Disconnect inbox
// @Description Disconnects the active WhatsApp inbox
// @Tags        Inboxes
// @Produce     json
// @Param       inboxId   path string true "Inbox ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/disconnect [post]
func (h *InboxHandler) Disconnect(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	if err := h.inboxSvc.Disconnect(inboxID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// QR godoc
// @Summary     Get QR code for pairing
// @Description Returns a QR code for pairing a new WhatsApp device
// @Tags        Inboxes
// @Produce     json
// @Param       inboxId   path string true "Inbox ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/qr [get]
func (h *InboxHandler) QR(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	qrCode, err := h.inboxSvc.GetQRCode(c.Context(), inboxID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", err.Error()))
	}

	if qrCode == "" {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", "No QR code available. Call connect first, then poll this endpoint."))
	}

	imageBytes, imgErr := qrcode.Encode(qrCode, qrcode.Medium, 256)
	qrBase64 := ""
	if imgErr == nil {
		qrBase64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageBytes)
	}

	return c.JSON(dto.SuccessResp(map[string]interface{}{
		"qr":    qrCode,
		"image": qrBase64,
	}))
}
