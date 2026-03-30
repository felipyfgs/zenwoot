package handler

import (
	"github.com/gofiber/fiber/v2"

	"wzap/internal/domain"
	"wzap/internal/dto"
	"wzap/internal/repo"
	"wzap/internal/service"
)

type SetupHandler struct {
	userRepo  *repo.UserRepository
	userSvc   *service.UserService
	authSvc   *service.AuthService
	accountID string
}

func NewSetupHandler(
	userRepo *repo.UserRepository,
	userSvc *service.UserService,
	authSvc *service.AuthService,
	accountID string,
) *SetupHandler {
	return &SetupHandler{
		userRepo:  userRepo,
		userSvc:   userSvc,
		authSvc:   authSvc,
		accountID: accountID,
	}
}

// Status godoc
// @Summary Check if initial setup is needed
// @Tags setup
// @Produce json
// @Success 200 {object} map[string]any
// @Router /api/v1/setup [get]
func (h *SetupHandler) Status(c *fiber.Ctx) error {
	count, err := h.userRepo.CountAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal error", "failed to check setup status"))
	}
	return c.JSON(dto.SuccessResp(fiber.Map{
		"needsSetup": count == 0,
	}))
}

type setupReq struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Create godoc
// @Summary Create first admin user (only when no users exist)
// @Tags setup
// @Accept json
// @Produce json
// @Param body body setupReq true "Admin credentials"
// @Success 200 {object} map[string]any
// @Router /api/v1/setup [post]
func (h *SetupHandler) Create(c *fiber.Ctx) error {
	count, err := h.userRepo.CountAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal error", "failed to check setup status"))
	}
	if count > 0 {
		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResp("Forbidden", "setup already completed"))
	}

	var req setupReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad request", "invalid body"))
	}
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad request", "email, name and password are required"))
	}

	_, err = h.userSvc.Create(c.Context(), h.accountID, service.CreateUserReq{
		Email:    req.Email,
		Name:     req.Name,
		Role:     domain.RoleAdmin,
		Password: req.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal error", err.Error()))
	}

	token, user, err := h.authSvc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResp("Internal error", "user created but login failed"))
	}

	return c.JSON(dto.SuccessResp(fiber.Map{
		"token": token,
		"user":  user,
	}))
}
