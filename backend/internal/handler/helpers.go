package handler

import "github.com/gofiber/fiber/v2"

// getInboxID returns the inbox domain ID injected by RequireInbox middleware.
func getInboxID(c *fiber.Ctx) string {
	if val := c.Locals("inboxId"); val != nil {
		return val.(string)
	}
	return ""
}
