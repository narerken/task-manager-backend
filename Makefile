# conf
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_USER     ?= postgres
DB_PASSWORD ?= Ernar17042006
DB_NAME     ?= todo
DB_SSLMODE  ?= disable

DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATIONS_DIR := migrations
MIGRATE        := migrate
MIGRATE_ARGS   := -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)"

# help message
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

# run
.PHONY: run build

run:
	go run ./...

build:
	go build -o bin/task-manager ./...

# commands
.PHONY: migrate-up migrate-down migrate-down-all migrate-to migrate-force migrate-version migrate-create

migrate-up:
	$(MIGRATE) $(MIGRATE_ARGS) up

migrate-down:
	$(MIGRATE) $(MIGRATE_ARGS) down 1

migrate-down-all:
	$(MIGRATE) $(MIGRATE_ARGS) down -all

migrate-to:
	$(MIGRATE) $(MIGRATE_ARGS) goto $(VERSION)

migrate-force:
	$(MIGRATE) $(MIGRATE_ARGS) force $(VERSION)

migrate-version:
	$(MIGRATE) $(MIGRATE_ARGS) version

migrate-create:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

# install
.PHONY: install-migrate-linux install-migrate-mac install-migrate-go

install-migrate-go:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# update go modules
.PHONY: tidy

tidy:
	go mod tidy
