package handler

import (
	"wzap/internal/broker"
	"wzap/internal/db"
	"wzap/internal/dto"
	"wzap/internal/storage"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
	database *db.DB
	nats     *broker.Nats
	minio    *storage.Minio
}

func NewHealthHandler(database *db.DB, nats *broker.Nats, minio *storage.Minio) *HealthHandler {
	return &HealthHandler{database: database, nats: nats, minio: minio}
}

// Check godoc
// @Summary     Health check
// @Description Returns the health status of the API and its dependencies
// @Tags        Health
// @Produce     json
// @Success     200 {object} dto.APIResponse
// @Router      /health [get]
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	dbOk := false
	if h.database != nil {
		dbOk = h.database.Health(c.Context()) == nil
	}

	natsOk := false
	if h.nats != nil {
		natsOk = h.nats.Health() == nil
	}

	minioOk := false
	if h.minio != nil {
		minioOk = h.minio.Health(c.Context()) == nil
	}

	overall := "UP"
	if !dbOk {
		overall = "DEGRADED"
	}

	status := map[string]interface{}{
		"status": overall,
		"services": map[string]bool{
			"database": dbOk,
			"nats":     natsOk,
			"minio":    minioOk,
		},
	}
	return c.JSON(dto.SuccessResp(status))
}
