# Project Variables
APP_NAME=api
CMD_DIR=./cmd/api
BUILD_DIR=./bin
MIGRATIONS_DIR=./database/migrations

# Load .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Database URL for migrations
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: all build run test ginkgo deps clean migrate-up migrate-down migrate-status migrate-create seed

all: build

## build: Build the application binary
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

## run: Run the application directly
run:
	@go run $(CMD_DIR)/main.go

## test: Run all tests using standard go test
test:
	@go test ./...

## ginkgo: Run all tests using ginkgo
ginkgo:
	@ginkgo ./...

## test-coverage: Run tests with coverage and show summary
test-coverage:
	@go test -coverprofile=coverage.out ./... || true
	@if [ -f coverage.out ]; then \
		go tool cover -func=coverage.out; \
		rm coverage.out; \
	fi

## test-coverage-html: Run tests with coverage and generate HTML report
test-coverage-html:
	@go test -coverprofile=coverage.out ./... || true
	@if [ -f coverage.out ]; then \
		go tool cover -html=coverage.out -o coverage.html; \
		rm coverage.out; \
		echo "Coverage report generated at coverage.html"; \
	fi

## deps: Download and tidy dependencies
deps:
	@go mod download
	@go mod tidy

## clean: Remove build artifacts
clean:
	@rm -rf $(BUILD_DIR)

## migrate-up: Run all pending migrations
migrate-up:
	@echo "Running migrations up..."
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

## migrate-down: Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

## migrate-status: Show current migration version
migrate-status:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

## migrate-create name=migration_name: Create a new migration file
migrate-create:
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## seed: Run database seeds
seed:
	@echo "Running database seeds..."
	@for file in $(shell ls database/seeds/*.sql); do \
		echo "Applying $$file..."; \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f $$file; \
	done

## swagger: Generate swagger documentation
swagger:
	swag init -g cmd/api/main.go -o ./docs/swagger --parseDependency --parseInternal

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## .*' $(MAKEFILE_LIST) | sed 's/^## //' | sort | awk 'BEGIN {FS = ": "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
