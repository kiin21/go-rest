.PHONY: help build run stop restart logs clean test migrate

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose restart

logs:
	docker compose logs -f

logs-app:
	docker compose logs -f app

logs-mysql:
	docker compose logs -f mysql

logs-kafka:
	docker compose logs -f kafka

logs-es:
	docker compose logs -f elasticsearch

clean:
	docker compose down -v

ps:
	docker compose ps

exec-app:
	docker compose exec app sh

exec-mysql:
	docker compose exec mysql mysql -u gouser -pgopassword intern_app

migrate:
	docker compose exec mysql mysql -u gouser -pgopassword intern_app < migrations/001_create_starter_table.sql

rebuild:
	docker compose down
	docker compose build --no-cache
	docker compose up -d

dev: ## Run in development mode (without Docker)
	go run cmd/api/main.go
