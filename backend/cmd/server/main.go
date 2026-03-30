package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/nats-io/nats.go"

	"github.com/felipyfgs/zenwoot/backend/internal/config"
	"github.com/felipyfgs/zenwoot/backend/internal/db"
	"github.com/felipyfgs/zenwoot/backend/internal/handlers"
	"github.com/felipyfgs/zenwoot/backend/internal/logger"
	"github.com/felipyfgs/zenwoot/backend/internal/middleware"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
	"github.com/felipyfgs/zenwoot/backend/internal/storage"
	"github.com/felipyfgs/zenwoot/backend/internal/workers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	logger.Init(cfg.App.Env == "development")
	log := logger.Log

	// --- Database ---
	bunDB := db.New(cfg.DB.DSN, cfg.App.Env == "development", log)
	if err := db.RunMigrations(context.Background(), bunDB); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	// --- Storage ---
	store, err := storage.NewS3Storage(cfg.MinIO)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to MinIO, attachments will be disabled")
		store = nil
	}

	// --- NATS ---
	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer nc.Close()

	// --- Repos ---
	convRepo := repo.NewConversationRepo(bunDB)
	contactRepo := repo.NewContactRepo(bunDB)
	msgRepo := repo.NewMessageRepo(bunDB)
	inboxRepo := repo.NewInboxRepo(bunDB)
	labelRepo := repo.NewLabelRepo(bunDB)
	webhookRepo := repo.NewWebhookRepo(bunDB)
	autoRepo := repo.NewAutomationRuleRepo(bunDB)
	notificationRepo := repo.NewNotificationRepo(bunDB)
	customFilterRepo := repo.NewCustomFilterRepo(bunDB)

	// --- Services ---
	authSvc := services.NewAuthService(bunDB, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	convSvc := services.NewConversationService(convRepo, nc)
	msgSvc := services.NewMessageService(msgRepo, convRepo, bunDB, nc)
	contactSvc := services.NewContactService(contactRepo, nc)
	_ = services.NewLabelService(labelRepo, bunDB)
	webhookSvc := services.NewWebhookService(webhookRepo)
	attachmentSvc := services.NewAttachmentService(bunDB, store)
	notificationSvc := services.NewNotificationService(notificationRepo)

	// --- Workers ---
	webhookWorker := workers.NewWebhookWorker(nc, bunDB, webhookSvc)
	autoWorker := workers.NewAutomationWorker(nc, autoRepo)
	snoozeWorker := workers.NewAutoResolveWorker(bunDB, 60*time.Second)

	if err := webhookWorker.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start webhook worker")
	}
	if err := autoWorker.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start automation worker")
	}
	snoozeWorker.Start()
	defer func() {
		webhookWorker.Stop()
		autoWorker.Stop()
		snoozeWorker.Stop()
	}()

	// --- Fiber App ---
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Public routes
	handlers.NewAuthHandler(authSvc).Register(app)

	// Webhook routes (no auth)
	webhookHandler := handlers.NewWebhookHandler()
	webhookGroup := app.Group("")
	webhookHandler.Register(webhookGroup)

	// Protected routes — /api/v1/accounts/:accountId
	authMw := middleware.AuthMiddleware(authSvc)
	tenantMw := middleware.TenantMiddleware(bunDB)

	account := app.Group("/api/v1/accounts/:accountId", authMw, tenantMw)

	handlers.NewConversationHandler(convSvc).Register(account)
	handlers.NewMessageHandler(msgSvc).Register(account)
	handlers.NewContactHandler(contactSvc).Register(account)
	handlers.NewInboxHandler(inboxRepo).Register(account)
	handlers.NewSearchHandler(contactRepo, convRepo).Register(account)
	handlers.NewAttachmentHandler(attachmentSvc).Register(account)
	handlers.NewNotificationHandler(notificationSvc).Register(account)
	handlers.NewCustomFilterHandler(customFilterRepo).Register(account)

	addr := cfg.App.Host + ":" + cfg.App.Port
	log.Info().Str("addr", addr).Msg("server listening")
	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
}
