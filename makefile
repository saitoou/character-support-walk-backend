SHELL := /bin/bash

include docker/.example.env

DOCKER_COMPOSE_FILE := docker/docker-compose.yml

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATION_PATH := migration

.PHONY: help setup up down reset logs migrate-up migrate-down test lint vulncheck generate

help: ## display this help screen
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## install develop tools
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/golang/mock/mockgen@v1.6.0
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	

up: ## docker compose up with build
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build

down: ## docker compose down
	docker compose -f $(DOCKER_COMPOSE_FILE) down

reset: ## docker compose down with volume
	docker compose -f $(DOCKER_COMPOSE_FILE) down -v

logs: ## docker compose logs
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

migrate-up: ## run migrate up
	migrate \
		-path $(MIGRATION_PATH) \
		-database "$(DB_URL)" \
		up

migrate-down: ## run migrate down
	migrate \
		-path $(MIGRATION_PATH) \
		-database "$(DB_URL)" \
		down

test: ## run go test
	go test ./...

lint: ## run staticcheck
	staticcheck ./...

vulncheck: ## run govulncheck
	govulncheck ./...

check: ## run test, lint, and vulncheck
	go test ./...
	staticcheck ./...
	govulncheck ./...

define oapi-codegen
	oapi-codegen -config api/config/servers.config.yaml api/openapi.yaml
	oapi-codegen -config api/config/models.config.yaml api/openapi.yaml
endef

generate: ## generate openapi code
	$(call oapi-codegen)