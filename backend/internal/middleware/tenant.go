package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"

	"github.com/felipyfgs/zenwoot/backend/internal/models"
)

func TenantMiddleware(db *bun.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		accountIDStr := c.Params("accountId")
		accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid accountId"})
		}
		userID, ok := c.Locals("user_id").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
		}
		var au models.AccountUser
		err = db.NewSelect().Model(&au).
			Where(`"account_id" = ? AND "user_id" = ?`, accountID, userID).
			Scan(c.Context())
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
		}
		c.Locals("account_id", accountID)
		c.Locals("user_role", au.Role)
		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c fiber.Ctx) error {
		role, _ := c.Locals("user_role").(int)
		if role != 1 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin required"})
		}
		return c.Next()
	}
}
