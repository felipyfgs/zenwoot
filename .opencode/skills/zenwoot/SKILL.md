---
name: zenwoot-backend
description: Backend Go + Fiber + Bun ORM para Zenwoot (fork Chatwoot). Use para implementar/revisar: models, repos, services, handlers, workers, migrations SQL, multi-tenancy, JWT auth. Convenções: snake_case em todas as camadas, zerolog para logs.
compatibility: opencode
---

# Zenwoot Backend — Skill

Stack: **Go + Fiber v3 + Bun ORM + PostgreSQL + NATS + Redis**

---

## Convenções de Código

### snake_case — REGRA ABSOLUTA

| Camada | Convenção | Exemplo |
|--------|-----------|---------|
| Coluna PostgreSQL | `snake_case` | `account_id`, `created_at` |
| SQL explícito | aspas duplas | `WHERE "account_id" = ?` |
| Bun struct tag | `snake_case` | `` `bun:"account_id,notnull"` `` |
| Go struct field | `PascalCase` | `AccountID int64` |
| JSON response | `snake_case` | `` `json:"account_id"` `` |

**Exemplo struct Go:**
```go
type Message struct {
    bun.BaseModel       `bun:"table:messages"`
    ID                  int64     `bun:"id,pk,autoincrement" json:"id"`
    AccountID           int64     `bun:"account_id,notnull"   json:"account_id"`
    ConversationID      int64     `bun:"conversation_id,notnull" json:"conversation_id"`
    SenderType          *string   `bun:"sender_type"         json:"sender_type"`
    SenderID            *int64    `bun:"sender_id"           json:"sender_id"`
    Content             *string   `bun:"content"             json:"content"`
    CreatedAt           time.Time `bun:"created_at,notnull"  json:"created_at"`
}
```

### Estrutura de Diretórios

```
backend/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go
│   ├── db/
│   │   ├── db.go              # Conexão Bun + zerolog
│   │   ├── migrator.go        # Bun migrations
│   │   └── migrations/*.sql
│   ├── models/                # Structs Bun
│   ├── repo/                  # Repositories
│   ├── services/              # Lógica de negócio
│   ├── handlers/              # Handlers Fiber
│   ├── middleware/            # JWT + Tenant
│   ├── workers/               # NATS consumers
│   └── logger/logger.go       # Zerolog centralizado
├── docker-compose.yaml
└── go.mod
```

---

## Camadas do Sistema

### Models (`internal/models/`)

Cada entidade tem seu próprio arquivo Go com struct Bun:
- `user.go`, `account.go`, `account_user.go`
- `inbox.go`, `channel_web_widget.go`
- `contact.go`, `conversation.go`, `message.go`
- `label.go`, `team.go`, `webhook.go`, `automation_rule.go`

### Repositories (`internal/repo/`)

Padrão: `repo/<entidade>.go` com métodos CRUD + filtros por `account_id`.

### Services (`internal/services/`)

Padrão: `svc.<Entidade>.Method(ctx, accountID, ...)`. Publica eventos para NATS.

### Handlers (`internal/handlers/`)

Padrão Fiber: `func (h *Handler) Register(app *fiber.Group)`. middlewares: `authMw`, `tenantMw`.

### Workers (`internal/workers/`)

Consumers NATS JetStream:
- `webhook.go` — dispatch webhooks
- `automation.go` — rules engine
- `auto_resolve.go` — wake snoozed conversations

---

## Multi-tenancy

**REGRA:** Todo repo DEVE filtrar por `account_id`.

```go
q.Where(`"account_id" = ?`, accountID)
```

Middleware `tenant.go` injeta `accountID` no context.

---

## Eventos NATS

| Subject | Publisher | Consumers |
|---------|-----------|-----------|
| `zenwoot.conversation.created` | ConversationService | WebhookWorker, AutomationWorker |
| `zenwoot.message.created` | MessageService | WebhookWorker, AutomationWorker |
| `zenwoot.contact.created` | ContactService | WebhookWorker |

---

## Logger (zerolog)

Log centralizado em `internal/logger/logger.go`:

```go
import "github.com/felipyfgs/zenwoot/backend/internal/logger"

logger.Info().Str("addr", addr).Msg("server listening")
logger.Error().Err(err).Msg("failed to process")
```

Bun queries são logadas automaticamente via bundebug.

---

## Migrations

Arquivos SQL em `internal/db/migrations/`. Formato: `YYYYMMDDHHMMSS_<nome>.up.sql`.

Executadas automaticamente ao iniciar o servidor via `db.RunMigrations()`.

---

## Quando Usar

- Criar/modificar models Bun
- Implementar repositories com filtros de tenant
- Adicionar novos handlers Fiber
- Escrever workers NATS
- Criar migrations SQL
- Revisar convenções de código
