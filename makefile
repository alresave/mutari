# =========================
# Mutari Makefile
# =========================

APP=mutari-api
GO=go
DATABASE_URL=postgres://postgres:mutari123@localhost:5432/mutari?sslmode=disable
GOCACHE_LOCAL=$(CURDIR)/.gocache
GOMODCACHE_LOCAL=$(CURDIR)/.gomodcache
GOLANGCI_LINT_CACHE_LOCAL=$(CURDIR)/.golangci-lint-cache

DB_CONTAINER=mutari-db
REDIS_CONTAINER=mutari-redis

DB_USER=postgres
DB_NAME=mutari
DB_PORT=5432

# =========================
# Phony Targets
# =========================

.PHONY: help dev watch build run clean \
up down restart logs ps \
db db-logs wait-db \
schema seed reset-db reseed \
swagger-init swagger-serve \
fmt tidy lint test test-integration test-all \
setup rebuild fresh prune

# =========================
# Help
# =========================

help:
	@echo ""
	@echo "Mutari Development Commands"
	@echo "---------------------------------"
	@echo "make dev            Run API locally"
	@echo "make watch          Run with hot reload (air)"
	@echo "make build          Build binary"
	@echo "make run            Run compiled binary"
	@echo "make clean          Remove binary"
	@echo ""
	@echo "Docker"
	@echo "make up             Start containers"
	@echo "make down           Stop containers"
	@echo "make restart        Restart containers"
	@echo "make logs           View logs"
	@echo "make ps             List containers"
	@echo ""
	@echo "Database"
	@echo "make db             Open psql shell"
	@echo "make schema         Load schema"
	@echo "make seed           Seed demo data"
	@echo "make reset-db       Recreate DB"
	@echo ""
	@echo "Utilities"
	@echo "make fmt            Format Go code"
	@echo "make tidy           Clean go.mod"
	@echo "make lint           Run linter"
	@echo "make test           Run tests"
	@echo "make test-integration Run integration tests with local Postgres"
	@echo "make test-all       Run unit/package tests and integration tests"
	@echo ""
	@echo "Project Setup"
	@echo "make setup          Full project setup"
	@echo "make fresh          Reset DB + build + run"
	@echo ""

# =========================
# Development
# =========================

dev:
	$(GO) run cmd/API/main.go

watch:
	air

build:
	$(GO) build -o $(APP) cmd/API/main.go

run:
	./$(APP)

clean:
	rm -f $(APP)

# =========================
# Docker
# =========================

build-docker:
	docker compose build --no-cache

up:
	docker compose up -d

down:
	docker compose down

restart: down up

logs:
	docker compose logs -f

ps:
	docker compose ps

# =========================
# Database
# =========================

db:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

db-logs:
	docker logs -f $(DB_CONTAINER)

wait-db:
	@echo "Waiting for database..."
	@until docker exec $(DB_CONTAINER) pg_isready -U $(DB_USER) > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "Database ready."

schema:
	docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < db/schema.sql

seed:
	docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < db/seed_base.sql

reseed: seed

reset-db:
	docker compose down -v
	docker compose up -d
	$(MAKE) wait-db
	$(MAKE) schema
	$(MAKE) seed

# =========================
# Swagger
# =========================

swagger-init:
	swag init --generalInfo cmd/API/main.go --output docs

swagger-serve:
	open http://localhost:8080/swagger/index.html

# =========================
# Go Utilities
# =========================

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

lint:
	env GOCACHE=$(GOCACHE_LOCAL) GOMODCACHE=$(GOMODCACHE_LOCAL) GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE_LOCAL) golangci-lint run ./...

test:
	$(GO) test ./...

test-integration:
	env 'DATABASE_URL=$(DATABASE_URL)' GOCACHE=$(GOCACHE_LOCAL) GOMODCACHE=$(GOMODCACHE_LOCAL) $(GO) test ./internal/modules/matching -v

test-all: test test-integration

# =========================
# Project Setup
# =========================

setup: up wait-db schema seed build

rebuild: build restart

fresh: reset-db build run

# =========================
# System cleanup
# =========================

prune:
	docker system prune -f
