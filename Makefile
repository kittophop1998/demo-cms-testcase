# Makefile for demo-notion-api

# Variables
IMAGE_NAME := demo-notion-api
VERSION := latest
PORT := 8080

# Help target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
.PHONY: dev
dev: ## Run development server with hot reload
	docker-compose -f docker-compose.dev.yml up --build

.PHONY: dev-down
dev-down: ## Stop development server
	docker-compose -f docker-compose.dev.yml down

# Production targets
.PHONY: build
build: ## Build production Docker image
	docker build -t $(IMAGE_NAME):$(VERSION) .

.PHONY: run
run: ## Run production server using docker-compose
	docker-compose up --build -d

.PHONY: stop
stop: ## Stop production server
	docker-compose down

.PHONY: logs
logs: ## View logs
	docker-compose logs -f

# Local development
.PHONY: local
local: ## Run locally without Docker
	go run main.go

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: clean
clean: ## Clean up Docker resources
	docker-compose down --volumes --remove-orphans
	docker-compose -f docker-compose.dev.yml down --volumes --remove-orphans
	docker system prune -f

# Security scan
.PHONY: security-scan
security-scan: ## Scan Docker image for vulnerabilities
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:latest image $(IMAGE_NAME):$(VERSION)

# Show image size
.PHONY: image-size
image-size: ## Show Docker image size
	docker images $(IMAGE_NAME):$(VERSION)

.DEFAULT_GOAL := help