package handler

import (
	"github.com/gofiber/fiber/v2"

	"wzap/internal/dto"
	"wzap/internal/service"
)

type NoteHandler struct {
	noteSvc        *service.NoteService
	participantSvc *service.ConversationParticipantService
}

func NewNoteHandler(noteSvc *service.NoteService, participantSvc *service.ConversationParticipantService) *NoteHandler {
	return &NoteHandler{noteSvc: noteSvc, participantSvc: participantSvc}
}

// CreateNote godoc
// @Summary     Create contact note
// @Description Adds a private note to a contact
// @Tags        Contacts
// @Accept      json
// @Produce     json
// @Param       contactId path string true "Contact ID"
// @Param       body      body  dto.NoteCreateReq true "Note data"
// @Success     201 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/contacts/{contactId}/notes [post]
func (h *NoteHandler) CreateNote(c *fiber.Ctx) error {
	contactID := c.Params("contactId")

	var req dto.NoteCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	userID, _ := c.Locals("userId").(string)

	note, err := h.noteSvc.Create(c.Context(), contactID, service.CreateNoteReq{
		Content: req.Content,
		UserID:  userID,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResp(note))
}

// ListNotes godoc
// @Summary     List contact notes
// @Description Returns all notes for a contact
// @Tags        Contacts
// @Produce     json
// @Param       contactId path string true "Contact ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/contacts/{contactId}/notes [get]
func (h *NoteHandler) ListNotes(c *fiber.Ctx) error {
	contactID := c.Params("contactId")

	notes, err := h.noteSvc.ListByContact(c.Context(), contactID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(notes))
}

// DeleteNote godoc
// @Summary     Delete contact note
// @Description Removes a note from a contact
// @Tags        Contacts
// @Produce     json
// @Param       contactId path string true "Contact ID"
// @Param       noteId    path string true "Note ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/contacts/{contactId}/notes/{noteId} [delete]
func (h *NoteHandler) DeleteNote(c *fiber.Ctx) error {
	noteID := c.Params("noteId")
	if err := h.noteSvc.Delete(c.Context(), noteID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// AddParticipant godoc
// @Summary     Add conversation participant
// @Description Subscribes a user to a conversation
// @Tags        Conversations
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Param       userId         path string true "User ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/participants/{userId} [post]
func (h *NoteHandler) AddParticipant(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")
	userID := c.Params("userId")
	inboxID := c.Query("inboxId")

	if err := h.participantSvc.Add(c.Context(), conversationID, userID, inboxID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// RemoveParticipant godoc
// @Summary     Remove conversation participant
// @Description Unsubscribes a user from a conversation
// @Tags        Conversations
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Param       userId         path string true "User ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/participants/{userId} [delete]
func (h *NoteHandler) RemoveParticipant(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")
	userID := c.Params("userId")
	inboxID := c.Query("inboxId")

	if err := h.participantSvc.Remove(c.Context(), conversationID, userID, inboxID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// ListParticipants godoc
// @Summary     List conversation participants
// @Description Returns all participants of a conversation
// @Tags        Conversations
// @Produce     json
// @Param       conversationId path string true "Conversation ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/conversations/{conversationId}/participants [get]
func (h *NoteHandler) ListParticipants(c *fiber.Ctx) error {
	conversationID := c.Params("conversationId")

	participants, err := h.participantSvc.List(c.Context(), conversationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(participants))
}
