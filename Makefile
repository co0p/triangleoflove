SHELL := /bin/sh

COMPOSE := docker compose
TEST_PROFILE := tests
BACKEND_DIR := services/backend
FRONTEND_DIR := services/frontend

.DEFAULT_GOAL := help

.PHONY: help backend-test frontend-test frontend-build acceptance-test docker-build dev

help:
	@printf '%s\n' \
	  'Supported targets:' \
	  '  make backend-test     Run backend Go tests' \
	  '  make frontend-test    Run frontend Vitest suite' \
	  '  make frontend-build   Build the frontend bundle' \
	  '  make acceptance-test  Run Docker acceptance tests and clean up afterward' \
	  '  make docker-build     Build the local Docker images used by the project workflows' \
	  '  make dev              Start the full stack locally, frontend served at http://localhost:5173'

$(FRONTEND_DIR)/node_modules: $(FRONTEND_DIR)/package-lock.json $(FRONTEND_DIR)/package.json
	npm --prefix $(FRONTEND_DIR) ci

backend-test:
	cd $(BACKEND_DIR) && go test ./...
frontend-test: $(FRONTEND_DIR)/node_modules
	npm --prefix $(FRONTEND_DIR) test

frontend-build: $(FRONTEND_DIR)/node_modules
	npm --prefix $(FRONTEND_DIR) run build

acceptance-test:
	@set -eu; \
	trap '$(COMPOSE) --profile $(TEST_PROFILE) down -v' EXIT; \
	$(COMPOSE) --profile $(TEST_PROFILE) up --build --abort-on-container-exit --exit-code-from api-tests api-tests

docker-build:
	$(COMPOSE) --profile $(TEST_PROFILE) build frontend backend db api-tests

dev:
	$(COMPOSE) up --build