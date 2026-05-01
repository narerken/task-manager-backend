# ─── Database config ──────────────────────────────────────────
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_USER     ?= postgres
DB_PASSWORD ?= Ernar17042006
DB_NAME     ?= todo
DB_SSLMODE  ?= disable

DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
M   := migrate -path migrations -database "$(DSN)"

# ─── App ──────────────────────────────────────────────────────
.PHONY: run build tidy
run:   go run ./...                      ## Run the app
build: go build -o bin/app ./...         ## Build binary
tidy:  go mod tidy                       ## Tidy modules

# ─── Migrations ───────────────────────────────────────────────
.PHONY: migrate-up migrate-down migrate-step migrate-back migrate-version migrate-create

migrate-up: ## Apply all pending migrations
	$(M) up

migrate-down: ## Roll back all migrations
	$(M) down -all

migrate-step: ## Apply 1 migration (or N): make migrate-step N=2
	$(M) up $(or $(N),1)

migrate-back: ## Roll back 1 migration (or N): make migrate-back N=2
	$(M) down $(or $(N),1)

migrate-version: ## Show current migration version
	$(M) version

migrate-create: ## New migration: make migrate-create NAME=add_something
	migrate create -ext sql -dir migrations -seq $(NAME)
