package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/helpers"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
	"github.com/felipyfgs/zenwoot/backend/internal/validators"
)

type AuthHandler struct {
	authSvc *services.AuthService
}

func NewAuthHandler(authSvc *services.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) SignIn(c fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	if !validators.ValidateEmail(body.Email) {
		return helpers.BadRequest(c, "invalid email format")
	}

	result, err := h.authSvc.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		return helpers.Unauthorized(c, "invalid credentials")
	}
	return c.JSON(result)
}

func (h *AuthHandler) SignUp(c fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return helpers.BadRequest(c, "invalid request body")
	}

	if ok, msg := validators.ValidateRequired(body.Name, "name"); !ok {
		return helpers.BadRequest(c, msg)
	}

	if !validators.ValidateEmail(body.Email) {
		return helpers.BadRequest(c, "invalid email format")
	}

	if ok, msg := validators.ValidatePassword(body.Password); !ok {
		return helpers.BadRequest(c, msg)
	}

	user, err := h.authSvc.Register(c.Context(), body.Name, body.Email, body.Password)
	if err != nil {
		return helpers.Unprocessable(c, err.Error())
	}
	return helpers.Created(c, user)
}
