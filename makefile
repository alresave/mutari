APP_NAME=mutari-api

DB_URL?=postgres://postgres:postgres@localhost:5432/mutari?sslmode=disable

.PHONY: help dev build run watch test tidy fmt clean seed bootstrap logs docker-up docker-down rebuild

help:
	@echo "Mutari Makefile"
	@echo ""
	@echo "make dev           Run API locally"
	@echo "make watch         Run with hot reload"
	@echo "make build         Build binary"
	@echo "make test          Run tests"
	@echo "make tidy          Go mod tidy"
	@echo "make fmt           Format code"
	@echo "make docker-up     Start containers"
	@echo "make docker-down   Stop containers"
	@echo "make seed          Seed database"
	@echo "make rebuild       Rebuild identity vectors"
	@echo "make bootstrap     Full setup for demo"

dev:
	go run cmd/api/main.go

watch:
	air

build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...

clean:
	rm -rf bin

docker-up:
	docker compose up -d

docker-down:
	docker compose down

logs:
	docker compose logs -f

seed:
	psql $(DB_URL) -f db/seed_demo.sql

rebuild:
	curl -X POST http://localhost:8080/admin/rebuild-profiles

bootstrap: docker-up seed rebuild
	@echo "Mutari ready"