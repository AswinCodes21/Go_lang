.PHONY: build run docker-up docker-down docker-restart

# Application name
APP_NAME=go-authentication

# Go commands
GO=go
GO_BUILD=$(GO) build
GO_RUN=$(GO) run

# Docker compose command
DOCKER_COMPOSE=docker-compose

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

# Show help
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application locally"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-restart - Rebuild and restart Docker containers"

# Default target
.DEFAULT_GOAL := help 