.PHONY: build run docker-up docker-down docker-restart db-up db-down

# Application name
APP_NAME=go-authentication

# Go commands
GO=go
GO_BUILD=$(GO) build
GO_RUN=$(GO) run

# Docker compose command
DOCKER_COMPOSE=docker-compose
DOCKER_COMPOSE_DB=docker-compose -f docker-compose.db.yml

# Build the application
build:
	$(GO_BUILD) -o $(APP_NAME) ./cmd/main.go

# Run the application locally
run:
	$(GO_RUN) ./cmd/main.go

# Start Docker containers
docker-up:
	$(DOCKER_COMPOSE) up -d

# Stop Docker containers
docker-down:
	$(DOCKER_COMPOSE) down

# Rebuild and restart Docker containers
docker-restart:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) up --build -d

# Start just the database
db-up:
	$(DOCKER_COMPOSE_DB) up -d

# Stop just the database
db-down:
	$(DOCKER_COMPOSE_DB) down

# Show help
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application locally"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-restart - Rebuild and restart Docker containers"
	@echo "  make db-up         - Start just the PostgreSQL database"
	@echo "  make db-down       - Stop just the PostgreSQL database"

# Default target
.DEFAULT_GOAL := help 