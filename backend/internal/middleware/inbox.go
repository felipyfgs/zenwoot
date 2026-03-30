package middleware

import (
	"wzap/internal/dto"
	"wzap/internal/repo"

	"github.com/gofiber/fiber/v2"
)

// RequireInbox resolves an inbox from :inboxId (by UUID or name).
// Injects into c.Locals:
//   - "sessionId" → inbox.ChannelID  (consumed by all WhatsApp service calls)
//   - "inboxId"   → inbox.ID         (consumed by domain handlers)
func RequireInbox(inboxRepo *repo.InboxRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		param := c.Params("inboxId")
		if param == "" {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResp("Bad Request", "inboxId is required"))
		}

		inbox, err := inboxRepo.FindByNameOrID(c.Context(), param)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResp("Not Found", "Inbox not found"))
		}

		// Session-scoped ApiKey: verify the key belongs to this inbox.
		// Auth middleware stores the session.ID in Locals("sessionId") when role == "session".
		// Since session.ID == inbox.ID (created together in SessionService.Create), we can compare directly.
		if c.Locals("authRole") == "session" {
			if c.Locals("sessionId") != inbox.ID {
				return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResp("Forbidden", "ApiKey not authorized for this inbox"))
			}
		}

		c.Locals("sessionId", inbox.ChannelID) // channelID for engine.GetClient
		c.Locals("inboxId", inbox.ID)
		return c.Next()
	}
}
