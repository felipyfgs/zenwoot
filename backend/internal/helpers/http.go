package helpers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func GetAccountID(c fiber.Ctx) int64 {
	if id, ok := c.Locals("account_id").(int64); ok {
		return id
	}
	return 0
}

func GetUserID(c fiber.Ctx) int64 {
	if id, ok := c.Locals("user_id").(int64); ok {
		return id
	}
	return 0
}

func ParseID(c fiber.Ctx, param string) (int64, error) {
	return strconv.ParseInt(c.Params(param), 10, 64)
}

func ParseQueryInt(c fiber.Ctx, key string, defaultVal int) int {
	val, err := strconv.Atoi(c.Query(key))
	if err != nil {
		return defaultVal
	}
	return val
}

func ParseQueryInt64(c fiber.Ctx, key string) (int64, error) {
	return strconv.ParseInt(c.Query(key), 10, 64)
}

func ParsePageParams(c fiber.Ctx) (page, limit int) {
	page = ParseQueryInt(c, "page", 1)
	limit = ParseQueryInt(c, "limit", 25)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func Success(c fiber.Ctx, data any) error {
	return c.JSON(fiber.Map{"data": data})
}

func Created(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

func BadRequest(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": message})
}

func NotFound(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": message})
}

func InternalError(c fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}

func Unprocessable(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": message})
}

func Unauthorized(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": message})
}
