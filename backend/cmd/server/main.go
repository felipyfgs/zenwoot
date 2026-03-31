package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/nats-io/nats.go"

	"github.com/felipyfgs/zenwoot/backend/internal/config"
	"github.com/felipyfgs/zenwoot/backend/internal/db"
	"github.com/felipyfgs/zenwoot/backend/internal/logger"
	"github.com/felipyfgs/zenwoot/backend/internal/middleware"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
	"github.com/felipyfgs/zenwoot/backend/internal/router"
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

	bunDB := db.New(cfg.DB.DSN, cfg.App.Env == "development", log)
	if err := db.RunMigrations(context.Background(), bunDB); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	store, err := storage.NewS3Storage(cfg.MinIO)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to MinIO, attachments will be disabled")
		store = nil
	}

	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer nc.Close()

	convRepo := repo.NewConversationRepo(bunDB)
	contactRepo := repo.NewContactRepo(bunDB)
	msgRepo := repo.NewMessageRepo(bunDB)
	inboxRepo := repo.NewInboxRepo(bunDB)
	labelRepo := repo.NewLabelRepo(bunDB)
	teamRepo := repo.NewTeamRepo(bunDB)
	webhookRepo := repo.NewWebhookRepo(bunDB)
	autoRepo := repo.NewAutomationRuleRepo(bunDB)
	notificationRepo := repo.NewNotificationRepo(bunDB)
	customFilterRepo := repo.NewCustomFilterRepo(bunDB)
	companyRepo := repo.NewCompanyRepo(bunDB)
	noteRepo := repo.NewNoteRepo(bunDB)
	cannedResponseRepo := repo.NewCannedResponseRepo(bunDB)
	campaignRepo := repo.NewCampaignRepo(bunDB)
	agentBotRepo := repo.NewAgentBotRepo(bunDB)
	macroRepo := repo.NewMacroRepo(bunDB)

	authSvc := services.NewAuthService(bunDB, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	convSvc := services.NewConversationService(convRepo, nc)
	msgSvc := services.NewMessageService(msgRepo, convRepo, bunDB, nc)
	contactSvc := services.NewContactService(contactRepo, nc)
	labelSvc := services.NewLabelService(labelRepo, bunDB)
	teamSvc := services.NewTeamService(teamRepo)
	automationSvc := services.NewAutomationService(autoRepo)
	webhookSvc := services.NewWebhookService(webhookRepo)
	attachmentSvc := services.NewAttachmentService(bunDB, store)
	notificationSvc := services.NewNotificationService(notificationRepo)
	companySvc := services.NewCompanyService(companyRepo)
	noteSvc := services.NewNoteService(noteRepo)
	cannedResponseSvc := services.NewCannedResponseService(cannedResponseRepo)
	campaignSvc := services.NewCampaignService(campaignRepo)
	agentBotSvc := services.NewAgentBotService(agentBotRepo)
	macroSvc := services.NewMacroService(macroRepo)

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

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	authMw := middleware.AuthMiddleware(authSvc)
	tenantMw := middleware.TenantMiddleware(bunDB)

	router.Setup(app, router.Config{
		AuthSvc:           authSvc,
		ConvSvc:           convSvc,
		MsgSvc:            msgSvc,
		ContactSvc:        contactSvc,
		LabelSvc:          labelSvc,
		TeamSvc:           teamSvc,
		AutomationSvc:     automationSvc,
		WebhookSvc:        webhookSvc,
		AttachmentSvc:     attachmentSvc,
		NotificationSvc:   notificationSvc,
		CompanySvc:        companySvc,
		NoteSvc:           noteSvc,
		CannedResponseSvc: cannedResponseSvc,
		CampaignSvc:       campaignSvc,
		AgentBotSvc:       agentBotSvc,
		MacroSvc:          macroSvc,
		ConvRepo:          convRepo,
		ContactRepo:       contactRepo,
		InboxRepo:         inboxRepo,
		CustomFilterRepo:  customFilterRepo,
	}, authMw, tenantMw)

	addr := cfg.App.Host + ":" + cfg.App.Port
	log.Info().Str("addr", addr).Msg("server listening")
	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
}
