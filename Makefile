include .env
export

# ==============================================================================
# VARIABLES
# ==============================================================================
BINARY_NAME      ?= azgher
CMD_PATH          = ./cmd/server
MIGRATIONS_PATH   = internal/postgres/migrations
DATABASE_URL     ?= postgres://postgres:postgres@localhost:5432/azgher?sslmode=disable
MIGRATE          ?= migrate
GOLANGCI_LINT    ?= golangci-lint
GO               ?= go

# ==============================================================================
# .PHONY TARGETS
# ==============================================================================
.PHONY: help deps sqlc build run test lint fmt vet tidy \
        migrate-install migrate-up migrate-down migrate-force migrate-new \
        docker-up docker-down clean

# ==============================================================================
# HELP
# ==============================================================================
help: ## Show this help message
	@printf "\033[36m%-30s\033[0m %s\n" "TARGET" "DESCRIPTION"
	@printf "\033[36m%-30s\033[0m %s\n" "------" "-----------"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ==============================================================================
# DEPENDENCIES & CODE GENERATION
# ==============================================================================
deps: ## Download and tidy go modules
	$(GO) mod download
	$(GO) mod tidy

sqlc: ## Generate type-safe code from sqlc.yaml
	sqlc generate

# ==============================================================================
# BUILD & RUN
# ==============================================================================
build: ## Build the binary into bin/
	@mkdir -p bin
	$(GO) build -o bin/$(BINARY_NAME) $(CMD_PATH)

run: ## Run the server (development)
	$(GO) run $(CMD_PATH)

# ==============================================================================
# TESTING & QUALITY
# ==============================================================================
test: ## Run all tests with race detector
	$(GO) test -v -race -count=1 ./...

lint: ## Run golangci-lint
	$(GOLANGCI_LINT) run ./...

fmt: ## Format code with gofmt
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

# ==============================================================================
# DATABASE MIGRATIONS (requires golang-migrate/migrate)
# ==============================================================================
migrate-install: ## Install golang-migrate CLI locally
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: ## Apply all pending migrations
	$(MIGRATE) -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) up

migrate-down: ## Rollback the last migration
	$(MIGRATE) -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) down 1

migrate-force: ## Force a specific migration version (usage: make migrate-force VERSION=1)
ifndef VERSION
	$(error VERSION is required. Example: make migrate-force VERSION=1)
endif
	$(MIGRATE) -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) force $(VERSION)

migrate-new: ## Create a new migration (usage: make migrate-new NAME=add_user_preferences)
ifndef NAME
	$(error NAME is required. Example: make migrate-new NAME=add_user_preferences)
endif
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

# ==============================================================================
# DOCKER / CONTAINERS
# ==============================================================================
docker-up: ## Start containers with docker-compose
	docker-compose up -d

docker-down: ## Stop containers with docker-compose
	docker-compose down

docker-logs: ## Tail logs from containers
	docker-compose logs -f

# ==============================================================================
# MAINTENANCE
# ==============================================================================
clean: ## Remove build artifacts and log files
	rm -rf bin/
	rm -f logs.json
