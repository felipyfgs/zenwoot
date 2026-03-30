.PHONY: help dev backend frontend up down install build

# Variables
BACKEND_DIR=backend
FRONTEND_DIR=frontend

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Run backend and frontend in dev mode
	@$(MAKE) -C $(BACKEND_DIR) dev &
	@$(MAKE) -C $(FRONTEND_DIR) dev

backend: ## Run backend only
	$(MAKE) -C $(BACKEND_DIR) dev

frontend: ## Run frontend only
	cd $(FRONTEND_DIR) && npm run dev

build: ## Build backend binary
	$(MAKE) -C $(BACKEND_DIR) build

install: ## Install all dependencies
	cd $(FRONTEND_DIR) && npm install
	cd $(BACKEND_DIR) && go mod download

up: ## Start all infrastructure services (Postgres, MinIO, NATS)
	docker compose up -d

down: ## Stop all infrastructure services
	docker compose down

down-clean: ## Stop services and remove volumes (DESTRUCTIVE)
	docker compose down -v

tidy: ## Tidy go modules
	$(MAKE) -C $(BACKEND_DIR) tidy

docs: ## Generate Swagger documentation
	$(MAKE) -C $(BACKEND_DIR) docs

lint: ## Run linter
	cd $(BACKEND_DIR) && golangci-lint run ./...
