package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"wzap/internal/domain"
	"wzap/internal/dto"
	"wzap/internal/service"
)

type UserHandler struct {
	userSvc *service.UserService
	teamSvc *service.TeamService
}

func NewUserHandler(userSvc *service.UserService, teamSvc *service.TeamService) *UserHandler {
	return &UserHandler{userSvc: userSvc, teamSvc: teamSvc}
}

// Create godoc
// @Summary     Create a new user
// @Description Creates a new user and adds to account
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       body body     dto.UserCreateReq true "User data"
// @Success     200  {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/users [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	if c.Locals("authRole") != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResp("Forbidden", "Admin access required"))
	}

	accountID := c.Locals("accountId").(string)

	var req dto.UserCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	user, err := h.userSvc.Create(c.Context(), accountID, service.CreateUserReq{
		Email:       req.Email,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Role:        domain.UserRole(req.Role),
		Provider:    req.Provider,
		UID:         req.UID,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResp(user))
}

// List godoc
// @Summary     List users
// @Description Returns all users for the account
// @Tags        Users
// @Produce     json
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/users [get]
func (h *UserHandler) List(c *fiber.Ctx) error {
	accountID := c.Locals("accountId").(string)

	users, err := h.userSvc.ListByAccount(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(users))
}

// Get godoc
// @Summary     Get user
// @Description Returns user by ID
// @Tags        Users
// @Produce     json
// @Param       userId path string true "User ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/users/{userId} [get]
func (h *UserHandler) Get(c *fiber.Ctx) error {
	userID := c.Params("userId")
	user, err := h.userSvc.Get(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", err.Error()))
	}
	return c.JSON(dto.SuccessResp(user))
}

// Update godoc
// @Summary     Update user
// @Description Updates user data
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       userId path string true "User ID"
// @Param       body body     dto.UserUpdateReq true "User data"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/users/{userId} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var req dto.UserUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}

	user, err := h.userSvc.Update(c.Context(), userID, service.UpdateUserReq{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		Role:        domain.UserRole(req.Role),
		Status:      domain.UserStatus(req.Status),
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(user))
}

// Delete godoc
// @Summary     Delete user
// @Description Removes user from system
// @Tags        Users
// @Produce     json
// @Param       userId path string true "User ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/users/{userId} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	if c.Locals("authRole") != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResp("Forbidden", "Admin access required"))
	}

	userID := c.Params("userId")
	if err := h.userSvc.Delete(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// AddInboxMember godoc
// @Summary     Add user to inbox
// @Description Assigns a user as inbox member
// @Tags        Users
// @Produce     json
// @Param       inboxId path string true "Inbox ID"
// @Param       userId  path string true "User ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/members/{userId} [post]
func (h *UserHandler) AddInboxMember(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	userID := c.Params("userId")

	if err := h.teamSvc.AddInboxMember(c.Context(), inboxID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// RemoveInboxMember godoc
// @Summary     Remove user from inbox
// @Description Removes user from inbox members
// @Tags        Users
// @Produce     json
// @Param       inboxId path string true "Inbox ID"
// @Param       userId  path string true "User ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/members/{userId} [delete]
func (h *UserHandler) RemoveInboxMember(c *fiber.Ctx) error {
	inboxID := getInboxID(c)
	userID := c.Params("userId")

	if err := h.teamSvc.RemoveInboxMember(c.Context(), inboxID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", err.Error()))
	}
	return c.JSON(dto.SuccessResp(nil))
}

// ListInboxMembers godoc
// @Summary     List inbox members
// @Description Returns all members of an inbox
// @Tags        Users
// @Produce     json
// @Param       inboxId path string true "Inbox ID"
// @Success     200 {object} dto.APIResponse
// @Security    ApiKey
// @Router      /api/v1/inboxes/{inboxId}/members [get]
func (h *UserHandler) ListInboxMembers(c *fiber.Ctx) error {
	inboxID := getInboxID(c)

	members, err := h.teamSvc.ListInboxMembers(c.Context(), inboxID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal Server Error", err.Error()))
	}
	return c.JSON(dto.SuccessResp(members))
}

func getAccountID(c *fiber.Ctx) string {
	if id, ok := c.Locals("accountId").(string); ok {
		return id
	}
	return "default"
}

func parseUUID(s string) (string, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", err
	}
	return s, nil
}
