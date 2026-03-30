package server

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/swagger"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	waChannel "wzap/internal/channel/whatsapp"
	"wzap/internal/dispatcher"
	"wzap/internal/handler"
	"wzap/internal/middleware"
	"wzap/internal/repo"
	"wzap/internal/service"
	"wzap/internal/ws"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) SetupRoutes() error {
	// ─── Repositories ────────────────────────────────────────────────────────────
	accountRepo := repo.NewAccountRepository(s.database.Pool)
	webhookRepo := repo.NewWebhookRepository(s.database.Pool)
	inboxRepo := repo.NewInboxRepository(s.database.Pool)
	waChannelRepo := repo.NewWhatsAppChannelRepository(s.database.Pool)
	contactRepo := repo.NewContactRepository(s.database.Pool)
	convRepo := repo.NewConversationRepository(s.database.Pool)
	msgRepo := repo.NewMessageRepository(s.database.Pool)
	contactInboxRepo := repo.NewContactInboxRepository(s.database.Pool)
	userRepo := repo.NewUserRepository(s.database.Pool)
	accountUserRepo := repo.NewAccountUserRepository(s.database.Pool)
	teamRepo := repo.NewTeamRepository(s.database.Pool)
	teamMemberRepo := repo.NewTeamMemberRepository(s.database.Pool)
	inboxMemberRepo := repo.NewInboxMemberRepository(s.database.Pool)
	labelRepo := repo.NewLabelRepository(s.database.Pool)
	conversationLabelRepo := repo.NewConversationLabelRepository(s.database.Pool)
	cannedResponseRepo := repo.NewCannedResponseRepository(s.database.Pool)
	noteRepo := repo.NewNoteRepository(s.database.Pool)
	participantRepo := repo.NewConversationParticipantRepository(s.database.Pool)

	// ─── Infrastructure ────────────────────────────────────────────────────────────
	disp := dispatcher.New(webhookRepo, s.nats)
	go disp.StartConsumer(s.ctx)

	// WebSocket Hub
	wsHub := ws.NewHub()
	go wsHub.Run()

	// Default account ID for single-tenant mode
	accountID := "default"

	// Incoming message processor
	msgProcessor := service.NewIncomingMessageProcessor(contactRepo, contactInboxRepo, convRepo, msgRepo, inboxRepo, wsHub, accountID)

	// Engine: connects WhatsApp to domain model
	engine, err := waChannel.NewInboxEngine(s.ctx, s.Config, waChannelRepo, inboxRepo, s.nats, disp, wsHub, msgProcessor, accountID)
	if err != nil {
		return err
	}

	if err := engine.ReconnectAll(s.ctx); err != nil {
		log.Warn().Err(err).Msg("Failed to reconnect inboxes on startup")
	}

	// ─── Services ──────────────────────────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo, s.Config.JWTSecret)
	inboxSvc := service.NewInboxService(webhookRepo, waChannelRepo, inboxRepo, engine, accountID)
	webhookSvc := service.NewWebhookService(webhookRepo)
	convSvc := service.NewConversationService(convRepo, msgRepo)
	contactCatalogSvc := service.NewContactCatalogService(contactRepo, convRepo)
	userSvc := service.NewUserService(userRepo, accountUserRepo)
	teamSvc := service.NewTeamService(teamRepo, teamMemberRepo, inboxMemberRepo)
	labelSvc := service.NewLabelService(labelRepo, conversationLabelRepo)
	assignSvc := service.NewAssignmentService(convRepo, inboxMemberRepo, teamMemberRepo, wsHub, accountID)
	cannedResponseSvc := service.NewCannedResponseService(cannedResponseRepo)
	noteSvc := service.NewNoteService(noteRepo, wsHub, accountID)
	participantSvc := service.NewConversationParticipantService(participantRepo, wsHub, accountID)

	// ─── Handlers ───────────────────────────────────────────────────────────────────
	authHandler := handler.NewAuthHandler(authSvc, userRepo)
	setupHandler := handler.NewSetupHandler(userRepo, userSvc, authSvc, accountID)
	healthHandler := handler.NewHealthHandler(s.database, s.nats, s.minio)
	inboxHandler := handler.NewInboxHandler(inboxSvc)
	webhookHandler := handler.NewWebhookHandler(webhookSvc)
	convHandler := handler.NewConversationHandler(convSvc)
	contactCatalogHandler := handler.NewContactCatalogHandler(contactCatalogSvc)
	userHandler := handler.NewUserHandler(userSvc, teamSvc)
	teamHandler := handler.NewTeamHandler(teamSvc)
	labelHandler := handler.NewLabelHandler(labelSvc)
	assignHandler := handler.NewAssignmentHandler(assignSvc)
	typingHandler := handler.NewTypingHandler(wsHub, accountID)
	cannedResponseHandler := handler.NewCannedResponseHandler(cannedResponseSvc)
	noteHandler := handler.NewNoteHandler(noteSvc, participantSvc)

	// ─── Middleware ─────────────────────────────────────────────────────────────────
	authMW := middleware.Auth(s.Config, accountRepo, accountID)
	requireInbox := middleware.RequireInbox(inboxRepo)

	// ─── Public ─────────────────────────────────────────────────────────────────────
	s.App.Get("/swagger/*", swagger.HandlerDefault)
	s.App.Get("/health", healthHandler.Check)
	s.App.Post("/api/v1/auth/login", authHandler.Login)
	s.App.Get("/api/v1/setup", setupHandler.Status)
	s.App.Post("/api/v1/setup", setupHandler.Create)

	// ─── /api/v1 ─────────────────────────────────────────────────────────────────────
	v1 := s.App.Group("/api/v1", authMW)
	v1.Get("/auth/me", authHandler.Me)

	// ── Inboxes ────────────────────────────────────────────────────────────────────
	v1.Post("/inboxes", inboxHandler.Create)
	v1.Get("/inboxes", inboxHandler.List)

	inbox := v1.Group("/inboxes/:inboxId", requireInbox)
	inbox.Get("/", inboxHandler.Get)
	inbox.Delete("/", inboxHandler.Delete)
	inbox.Post("/connect", inboxHandler.Connect)
	inbox.Post("/disconnect", inboxHandler.Disconnect)
	inbox.Get("/qr", inboxHandler.QR)

	// ── Webhooks (per inbox) ────────────────────────────────────────────────────────
	inbox.Post("/webhooks", webhookHandler.Create)
	inbox.Get("/webhooks", webhookHandler.List)
	inbox.Delete("/webhooks/:wid", webhookHandler.Delete)

	// ── Conversations (per inbox) ───────────────────────────────────────────────────
	inbox.Get("/conversations", convHandler.ListByInbox)

	// ── /api/v1/conversations ───────────────────────────────────────────────────────
	conv := v1.Group("/conversations")
	conv.Get("/:conversationId", convHandler.Get)
	conv.Post("/:conversationId/toggle_status", convHandler.ToggleStatus)
	conv.Post("/:conversationId/read", convHandler.MarkRead)
	conv.Get("/:conversationId/messages", convHandler.ListMessages)
	conv.Post("/:conversationId/assign", assignHandler.Assign)
	conv.Post("/:conversationId/unassign", assignHandler.Unassign)
	conv.Post("/:conversationId/priority", assignHandler.UpdatePriority)
	conv.Post("/:conversationId/snooze", assignHandler.Snooze)
	conv.Post("/:conversationId/unsnooze", assignHandler.UnSnooze)
	conv.Get("/:conversationId/labels", labelHandler.GetConversationLabels)
	conv.Post("/:conversationId/labels/:labelId", labelHandler.AddToConversation)
	conv.Delete("/:conversationId/labels/:labelId", labelHandler.RemoveFromConversation)
	conv.Put("/:conversationId/labels", labelHandler.SetConversationLabels)

	// ── /api/v1/contacts (catalog) ───────────────────────────────────────────────────
	contacts := v1.Group("/contacts")
	contacts.Get("/", contactCatalogHandler.List)
	contacts.Get("/:contactId", contactCatalogHandler.Get)
	contacts.Get("/:contactId/conversations", contactCatalogHandler.ListConversations)
	contacts.Post("/:contactId/notes", noteHandler.CreateNote)
	contacts.Get("/:contactId/notes", noteHandler.ListNotes)
	contacts.Delete("/:contactId/notes/:noteId", noteHandler.DeleteNote)

	// ── /api/v1/users ────────────────────────────────────────────────────────────────
	users := v1.Group("/users")
	users.Post("/", userHandler.Create)
	users.Get("/", userHandler.List)
	users.Get("/:userId", userHandler.Get)
	users.Put("/:userId", userHandler.Update)
	users.Delete("/:userId", userHandler.Delete)

	// ── /api/v1/teams ─────────────────────────────────────────────────────────────────
	teams := v1.Group("/teams")
	teams.Post("/", teamHandler.Create)
	teams.Get("/", teamHandler.List)
	teams.Get("/:teamId", teamHandler.Get)
	teams.Put("/:teamId", teamHandler.Update)
	teams.Delete("/:teamId", teamHandler.Delete)
	teams.Get("/:teamId/members", teamHandler.ListMembers)
	teams.Post("/:teamId/members/:userId", teamHandler.AddMember)
	teams.Delete("/:teamId/members/:userId", teamHandler.RemoveMember)

	// ── /api/v1/labels ────────────────────────────────────────────────────────────────
	labels := v1.Group("/labels")
	labels.Post("/", labelHandler.Create)
	labels.Get("/", labelHandler.List)
	labels.Get("/:labelId", labelHandler.Get)
	labels.Put("/:labelId", labelHandler.Update)
	labels.Delete("/:labelId", labelHandler.Delete)

	// ── /api/v1/canned-responses ─────────────────────────────────────────────────────
	canned := v1.Group("/canned-responses")
	canned.Post("/", cannedResponseHandler.Create)
	canned.Get("/", cannedResponseHandler.List)
	canned.Put("/:id", cannedResponseHandler.Update)
	canned.Delete("/:id", cannedResponseHandler.Delete)

	// ── Conversation participants ─────────────────────────────────────────────────────
	conv.Get("/:conversationId/participants", noteHandler.ListParticipants)
	conv.Post("/:conversationId/participants/:userId", noteHandler.AddParticipant)
	conv.Delete("/:conversationId/participants/:userId", noteHandler.RemoveParticipant)

	// ── Inbox members ────────────────────────────────────────────────────────────────
	inbox.Get("/members", userHandler.ListInboxMembers)
	inbox.Post("/members/:userId", userHandler.AddInboxMember)
	inbox.Delete("/members/:userId", userHandler.RemoveInboxMember)

	// ── Typing indicators ──────────────────────────────────────────────────────────
	v1.Post("/typing", typingHandler.SetTyping)

	// ─── WebSocket endpoint ──────────────────────────────────────────────────────────
	v1.Get("/ws", websocketHandler(wsHub, accountID, s.Config.AllowedOrigins))

	return nil
}

// websocketHandler returns a Fiber handler for WebSocket connections.
func websocketHandler(hub *ws.Hub, defaultAccountID, allowedOrigins string) fiber.Handler {
	originChecker := func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true // No origin header = direct connection
		}
		if allowedOrigins == "" {
			return true // No restriction configured
		}
		if allowedOrigins == "*" {
			return true
		}
		// Check if origin is in allowed list (comma-separated)
		for _, allowed := range splitAndTrim(allowedOrigins) {
			if origin == allowed {
				return true
			}
		}
		return false
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     originChecker,
	}

	return adaptor.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("WebSocket upgrade failed")
			return
		}
		defer conn.Close()

		// Extract accountID and inboxIDs from query params
		accountID := r.URL.Query().Get("accountId")
		if accountID == "" {
			accountID = defaultAccountID
		}
		inboxIDs := r.URL.Query()["inboxId"]

		client := hub.NewClient(accountID, inboxIDs)
		defer hub.RemoveClient(client)

		// Read loop (keep connection alive)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}

		// Write loop
		for msg := range client.Send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}))
}

func splitAndTrim(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
