package router

import (
	"github.com/gofiber/fiber/v3"

	"github.com/felipyfgs/zenwoot/backend/internal/handlers"
	"github.com/felipyfgs/zenwoot/backend/internal/repo"
	"github.com/felipyfgs/zenwoot/backend/internal/services"
)

type Config struct {
	AuthSvc           *services.AuthService
	ConvSvc           *services.ConversationService
	MsgSvc            *services.MessageService
	ContactSvc        *services.ContactService
	LabelSvc          *services.LabelService
	TeamSvc           *services.TeamService
	AutomationSvc     *services.AutomationService
	WebhookSvc        *services.WebhookService
	AttachmentSvc     *services.AttachmentService
	NotificationSvc   *services.NotificationService
	CompanySvc        *services.CompanyService
	NoteSvc           *services.NoteService
	CannedResponseSvc *services.CannedResponseService
	CampaignSvc       *services.CampaignService
	AgentBotSvc       *services.AgentBotService
	MacroSvc          *services.MacroService

	ConvRepo         *repo.ConversationRepo
	ContactRepo      *repo.ContactRepo
	InboxRepo        *repo.InboxRepo
	CustomFilterRepo *repo.CustomFilterRepo
}

func Setup(app *fiber.App, cfg Config, authMw fiber.Handler, tenantMw fiber.Handler) {
	// Auth routes (public)
	authH := handlers.NewAuthHandler(cfg.AuthSvc)
	app.Post("/api/v1/auth/sign_in", authH.SignIn)
	app.Post("/api/v1/auth/sign_up", authH.SignUp)

	// Webhook routes (public - incoming webhooks)
	webhookH := handlers.NewWebhookHandler()
	app.Post("/webhooks/whatsapp", webhookH.WhatsApp)
	app.Post("/webhooks/telegram", webhookH.Telegram)
	app.Post("/webhooks/email", webhookH.Email)
	app.Post("/webhooks/debug", webhookH.Debug)

	// Protected routes
	account := app.Group("/api/v1/accounts/:accountId", authMw, tenantMw)

	// Conversations
	convH := handlers.NewConversationHandler(cfg.ConvSvc)
	account.Get("/conversations", convH.List)
	account.Post("/conversations", convH.Create)
	account.Get("/conversations/:id", convH.Get)
	account.Post("/conversations/:id/resolve", convH.Resolve)
	account.Post("/conversations/:id/reopen", convH.Reopen)
	account.Post("/conversations/:id/snooze", convH.Snooze)
	account.Post("/conversations/:id/assignments", convH.Assign)
	account.Delete("/conversations/:id", convH.Delete)

	// Messages
	msgH := handlers.NewMessageHandler(cfg.MsgSvc)
	account.Get("/conversations/:conversation_id/messages", msgH.List)
	account.Post("/conversations/:conversation_id/messages", msgH.Create)
	account.Get("/messages/:id", msgH.Get)
	account.Patch("/messages/:id", msgH.Update)
	account.Delete("/messages/:id", msgH.Delete)

	// Contacts
	contactH := handlers.NewContactHandler(cfg.ContactSvc)
	account.Get("/contacts", contactH.List)
	account.Post("/contacts", contactH.Create)
	account.Get("/contacts/:id", contactH.Get)
	account.Patch("/contacts/:id", contactH.Update)
	account.Delete("/contacts/:id", contactH.Delete)

	// Inboxes
	inboxH := handlers.NewInboxHandler(cfg.InboxRepo)
	account.Get("/inboxes", inboxH.List)
	account.Post("/inboxes", inboxH.Create)
	account.Get("/inboxes/:id", inboxH.Get)
	account.Patch("/inboxes/:id", inboxH.Update)
	account.Delete("/inboxes/:id", inboxH.Delete)

	// Labels
	labelH := handlers.NewLabelHandler(cfg.LabelSvc)
	account.Get("/labels", labelH.List)
	account.Post("/labels", labelH.Create)
	account.Get("/labels/:id", labelH.Get)
	account.Patch("/labels/:id", labelH.Update)
	account.Delete("/labels/:id", labelH.Delete)
	account.Post("/conversations/:conversation_id/labels", labelH.AddToConversation)
	account.Delete("/conversations/:conversation_id/labels/:label_id", labelH.RemoveFromConversation)

	// Teams
	teamH := handlers.NewTeamHandler(cfg.TeamSvc)
	account.Get("/teams", teamH.List)
	account.Post("/teams", teamH.Create)
	account.Get("/teams/:id", teamH.Get)
	account.Patch("/teams/:id", teamH.Update)
	account.Delete("/teams/:id", teamH.Delete)
	account.Get("/teams/:id/members", teamH.ListMembers)
	account.Post("/teams/:id/members", teamH.AddMember)
	account.Delete("/teams/:id/members/:user_id", teamH.RemoveMember)

	// Automations
	autoH := handlers.NewAutomationHandler(cfg.AutomationSvc)
	account.Get("/automations", autoH.List)
	account.Post("/automations", autoH.Create)
	account.Get("/automations/:id", autoH.Get)
	account.Patch("/automations/:id", autoH.Update)
	account.Delete("/automations/:id", autoH.Delete)

	// Search
	searchH := handlers.NewSearchHandler(cfg.ContactRepo, cfg.ConvRepo)
	account.Get("/search", searchH.Search)

	// Attachments
	attachH := handlers.NewAttachmentHandler(cfg.AttachmentSvc)
	account.Post("/attachments", attachH.Upload)
	account.Get("/attachments/:id", attachH.Get)
	account.Get("/attachments/:id/download", attachH.Download)

	// Notifications
	notifH := handlers.NewNotificationHandler(cfg.NotificationSvc)
	account.Get("/notifications", notifH.List)
	account.Get("/notifications/unread_count", notifH.UnreadCount)
	account.Post("/notifications/:id/read", notifH.MarkAsRead)
	account.Post("/notifications/read_all", notifH.MarkAllAsRead)
	account.Get("/notification_settings", notifH.GetSettings)
	account.Put("/notification_settings", notifH.UpdateSettings)

	// Custom Filters
	cfH := handlers.NewCustomFilterHandler(cfg.CustomFilterRepo)
	account.Get("/custom_filters", cfH.List)
	account.Post("/custom_filters", cfH.Create)
	account.Put("/custom_filters/:id", cfH.Update)
	account.Delete("/custom_filters/:id", cfH.Delete)

	// Companies
	companyH := handlers.NewCompanyHandler(cfg.CompanySvc)
	account.Get("/companies", companyH.List)
	account.Post("/companies", companyH.Create)
	account.Get("/companies/search", companyH.Search)
	account.Get("/companies/:id", companyH.Get)
	account.Patch("/companies/:id", companyH.Update)
	account.Delete("/companies/:id", companyH.Delete)

	// Notes
	noteH := handlers.NewNoteHandler(cfg.NoteSvc)
	account.Get("/contacts/:contact_id/notes", noteH.List)
	account.Post("/contacts/:contact_id/notes", noteH.Create)
	account.Get("/notes/:id", noteH.Get)
	account.Patch("/notes/:id", noteH.Update)
	account.Delete("/notes/:id", noteH.Delete)

	// Canned Responses
	cannedH := handlers.NewCannedResponseHandler(cfg.CannedResponseSvc)
	account.Get("/canned_responses", cannedH.List)
	account.Post("/canned_responses", cannedH.Create)
	account.Get("/canned_responses/:id", cannedH.Get)
	account.Patch("/canned_responses/:id", cannedH.Update)
	account.Delete("/canned_responses/:id", cannedH.Delete)

	// Campaigns
	campaignH := handlers.NewCampaignHandler(cfg.CampaignSvc)
	account.Get("/campaigns", campaignH.List)
	account.Post("/campaigns", campaignH.Create)
	account.Get("/campaigns/:id", campaignH.Get)
	account.Patch("/campaigns/:id", campaignH.Update)
	account.Delete("/campaigns/:id", campaignH.Delete)

	// Agent Bots
	agentBotH := handlers.NewAgentBotHandler(cfg.AgentBotSvc)
	account.Get("/agent_bots", agentBotH.List)
	account.Post("/agent_bots", agentBotH.Create)
	account.Get("/agent_bots/:id", agentBotH.Get)
	account.Patch("/agent_bots/:id", agentBotH.Update)
	account.Delete("/agent_bots/:id", agentBotH.Delete)

	// Macros
	macroH := handlers.NewMacroHandler(cfg.MacroSvc)
	account.Get("/macros", macroH.List)
	account.Post("/macros", macroH.Create)
	account.Get("/macros/:id", macroH.Get)
	account.Patch("/macros/:id", macroH.Update)
	account.Delete("/macros/:id", macroH.Delete)

	// Webhook CRUD (protected)
	webhookCrudH := handlers.NewWebhookHandlerWithService(cfg.WebhookSvc)
	account.Get("/webhooks", webhookCrudH.List)
	account.Post("/webhooks", webhookCrudH.Create)
	account.Get("/webhooks/:id", webhookCrudH.Get)
	account.Patch("/webhooks/:id", webhookCrudH.Update)
	account.Delete("/webhooks/:id", webhookCrudH.Delete)
}
