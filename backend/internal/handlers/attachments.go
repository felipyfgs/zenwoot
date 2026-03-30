package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type AttachmentHandler struct {
	svc *services.AttachmentService
}

func NewAttachmentHandler(svc *services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{svc: svc}
}

func (h *AttachmentHandler) Register(rg fiber.Router) {
	rg.Post("/attachments", h.Upload)
	rg.Get("/attachments/:id", h.Get)
	rg.Get("/attachments/:id/download", h.Download)
}

func (h *AttachmentHandler) Upload(c fiber.Ctx) error {
	accountID := helpers.GetAccountID(c)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file required"})
	}

	messageID, err := strconv.ParseInt(c.FormValue("message_id"), 10, 64)
	if err != nil {
		messageID = 0
	}

	ext := ""
	if len(file.Filename) > 0 {
		for i := len(file.Filename) - 1; i >= 0; i-- {
			if file.Filename[i] == '.' {
				ext = file.Filename[i+1:]
				break
			}
		}
	}

	attachment, err := h.svc.Create(c.Context(), services.CreateAttachmentInput{
		MessageID: messageID,
		AccountID: accountID,
		File:      file,
		FileType:  0,
		Extension: ext,
	})
	if err != nil {
		return helpers.InternalError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(attachment)
}

func (h *AttachmentHandler) Get(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return helpers.BadRequest(c, "invalid id")
	}

	attachment, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return helpers.NotFound(c, "attachment not found")
	}
	if attachment == nil {
		return helpers.NotFound(c, "attachment not found")
	}

	url, err := h.svc.GetDownloadURL(c.Context(), attachment)
	if err != nil {
		return helpers.InternalError(c, err)
	}

	return c.JSON(fiber.Map{
		"id":        attachment.ID,
		"url":       url,
		"file_type": attachment.FileType,
		"extension": attachment.Extension,
	})
}

func (h *AttachmentHandler) Download(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return helpers.BadRequest(c, "invalid id")
	}

	attachment, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return helpers.NotFound(c, "attachment not found")
	}
	if attachment == nil {
		return helpers.NotFound(c, "attachment not found")
	}

	url, err := h.svc.GetDownloadURL(c.Context(), attachment)
	if err != nil {
		return helpers.InternalError(c, err)
	}

	return c.Status(fiber.StatusFound).JSON(fiber.Map{"url": url})
}
