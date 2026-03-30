package handler

import (
	"github.com/gofiber/fiber/v2"

	"wzap/internal/dto"
	"wzap/internal/repo"
	"wzap/internal/service"
)

type AuthHandler struct {
	authSvc  *service.AuthService
	userRepo *repo.UserRepository
}

func NewAuthHandler(authSvc *service.AuthService, userRepo *repo.UserRepository) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, userRepo: userRepo}
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginReq true "Credentials"
// @Success 200 {object} map[string]any
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad request", "invalid body"))
	}
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad request", "email and password are required"))
	}

	token, user, err := h.authSvc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResp("Unauthorized", "invalid credentials"))
	}

	return c.JSON(dto.SuccessResp(fiber.Map{
		"token": token,
		"user":  user,
	}))
}

// Me godoc
// @Summary Get current authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResp("Unauthorized", "missing user context"))
	}

	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not found", "user not found"))
	}

	return c.JSON(dto.SuccessResp(user))
}
