.DEFAULT_GOAL := help

.PHONY: help run build build-prod test test-race lint fmt tidy migrate-up migrate-down db-reset

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"}; /^[a-zA-Z_-]+:.*##/ {printf "%-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Start the API (loads .env when exported by your shell)
	@set -a; [ ! -f .env ] || . ./.env; set +a; go run ./cmd/api

build: ## Build the API binary
	go build -o bin/scholarbuddy-api ./cmd/api

build-prod: ## Build production binary without Swagger/OpenAPI routes
	go build -tags production -o bin/scholarbuddy-api ./cmd/api

test: ## Run all unit tests
	go test ./...

test-race: ## Run tests with the race detector
	go test -race ./...

fmt: ## Format Go source files
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

tidy: ## Synchronize Go module dependencies
	go mod tidy

migrate-up: ## Apply the schema migration using the built-in Go runner
	@set -a; [ ! -f .env ] || . ./.env; set +a; go run ./cmd/migrate up

migrate-down: ## Roll back the schema migration using the built-in Go runner
	@set -a; [ ! -f .env ] || . ./.env; set +a; go run ./cmd/migrate down

db-reset: ## DESTRUCTIVE: erase ScholarBuddy tables and migration history
	@set -a; [ ! -f .env ] || . ./.env; set +a; go run ./cmd/migrate reset
