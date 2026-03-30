package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"

	"wzap/internal/broker"
	"wzap/internal/config"
	"wzap/internal/db"
	"wzap/internal/middleware"
	"wzap/internal/storage"
)

type Server struct {
	App    *fiber.App
	Config *config.Config

	database *db.DB
	nats     *broker.Nats
	minio    *storage.Minio
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(cfg *config.Config, database *db.DB, n *broker.Nats, m *storage.Minio) *Server {
	app := fiber.New(fiber.Config{
		ServerHeader:          "wzap",
		DisableStartupMessage: true,
		BodyLimit:             50 * 1024 * 1024, // 50 MB max body size
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Middlewares
	app.Use(middleware.Recovery())
	app.Use(middleware.Logger())
	allowOrigins := cfg.AllowedOrigins
	if allowOrigins == "" {
		allowOrigins = "*"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, ApiKey",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		App:      app,
		Config:   cfg,
		database: database,
		nats:     n,
		minio:    m,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.Config.ServerHost, s.Config.Port)
	log.Info().Str("addr", addr).Msg("Starting API server")
	return s.App.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down API server")
	s.cancel()

	// Fiber shutdown might block, we wrap it in a channel with context timeout
	done := make(chan error, 1)
	go func() {
		done <- s.App.Shutdown()
	}()

	select {
	case <-ctx.Done():
		log.Warn().Msg("API server shutdown timed out")
		return ctx.Err()
	case err := <-done:
		log.Info().Msg("API server stopped gracefully")
		return err
	}
}
