package middleware

import (
	"strings"
	"wzap/internal/config"
	"wzap/internal/dto"
	"wzap/internal/repo"
	"wzap/internal/service"

	"github.com/gofiber/fiber/v2"
)

// Auth validates the ApiKey header.
// - If ApiKey matches cfg.APIKey → admin role (full access)
// - If ApiKey matches an account's apiKey → account role (account-scoped access)
// - Otherwise → 401 Unauthorized
func Auth(cfg *config.Config, accountRepo *repo.AccountRepository, accountID string) fiber.Handler {
	authSvc := service.NewAuthService(nil, cfg.JWTSecret)

	return func(c *fiber.Ctx) error {
		// ── 1. Bearer JWT (primary — frontend) ──────────────────────────────────
		if authHeader := c.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authSvc.ParseToken(tokenStr)
			if err == nil {
				c.Locals("authRole", string(claims.Role))
				c.Locals("userId", claims.UserID)
				c.Locals("accountId", accountID)
				return c.Next()
			}
		}

		// ── 2. ApiKey fallback (integrations / webhooks) ─────────────────────────
		if cfg.APIKey == "" {
			c.Locals("authRole", "admin")
			return c.Next()
		}

		apiKey := c.Get("ApiKey")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResp("Unauthorized", "Missing credentials"))
		}

		if apiKey == cfg.APIKey {
			c.Locals("authRole", "admin")
			return c.Next()
		}

		account, err := accountRepo.FindByAPIKey(c.Context(), apiKey)
		if err == nil {
			c.Locals("authRole", "account")
			c.Locals("accountId", account.ID)
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResp("Unauthorized", "Invalid credentials"))
	}
}
