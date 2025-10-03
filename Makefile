# Migration Commands
.PHONY: migrate-up migrate-down migrate-status migrate-reset migrate-help

# Run database migrations (up)
migrate-up:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go -action up

# Rollback database migrations
migrate-down:
	@echo "Rolling back database migrations..."
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make migrate-down VERSION=<target_version>"; \
		echo "Example: make migrate-down VERSION=1.0.0"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go -action down -version $(VERSION)

# Show migration status
migrate-status:
	@echo "Checking migration status..."
	go run cmd/migrate/main.go -action status

# Reset database (drop all tables and recreate)
migrate-reset:
	@echo "Resetting database..."
	go run cmd/migrate/main.go -action reset

# Show migration help
migrate-help:
	@echo "Available migration commands:"
	@echo "  make migrate-up              - Run all pending migrations"
	@echo "  make migrate-down VERSION=x  - Rollback to specific version"
	@echo "  make migrate-status          - Show migration status"
	@echo "  make migrate-reset           - Drop all tables and recreate"
	@echo "  make migrate-help            - Show this help"

# License management commands
.PHONY: license-generate license-list license-status license-validate license-help

# Generate a new license
license-generate:
	@echo "Generating license..."
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make license-generate NAME='License Name' [MAX_USERS=1] [DAYS=365]"; \
		echo "Example: make license-generate NAME='Personal License' MAX_USERS=1 DAYS=365"; \
		exit 1; \
	fi
	go run cmd/license/main.go -action generate -name "$(NAME)" -description "$(DESCRIPTION)" -max-users $(MAX_USERS) -days $(DAYS) -created-by "$(CREATED_BY)"

# List all licenses
license-list:
	@echo "Listing all licenses..."
	go run cmd/license/main.go -action list

# Check license status
license-status:
	@echo "Checking license status..."
	@if [ -n "$(KEY)" ]; then \
		go run cmd/license/main.go -action status -key "$(KEY)"; \
	elif [ -n "$(USER_ID)" ]; then \
		go run cmd/license/main.go -action status -user-id "$(USER_ID)"; \
	else \
		go run cmd/license/main.go -action status; \
	fi

# Validate license key
license-validate:
	@echo "Validating license..."
	@if [ -z "$(KEY)" ]; then \
		echo "Usage: make license-validate KEY=FBT-XXXX-XXXX-XXXX-XXXX"; \
		exit 1; \
	fi
	go run cmd/license/main.go -action validate -key "$(KEY)"

# Show license help
license-help:
	@echo "Available license commands:"
	@echo "  make license-generate NAME='License Name' [MAX_USERS=1] [DAYS=365] - Generate new license"
	@echo "  make license-list                                         - List all licenses"
	@echo "  make license-status [KEY=xxx] [USER_ID=xxx]               - Check license status"
	@echo "  make license-validate KEY=xxx                            - Validate license key"
	@echo "  make license-help                                        - Show this help"

# Build commands
.PHONY: build run clean test

# Build the application
build:
	@echo "Building application..."
	go build -o bin/bot cmd/bot/main.go
	go build -o bin/migrate cmd/migrate/main.go
	go build -o bin/license cmd/license/main.go

# Run the application
run:
	@echo "Starting bot..."
	go run cmd/bot/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Development commands
.PHONY: dev-setup dev-migrate dev-run

# Setup development environment
dev-setup:
	@echo "Setting up development environment..."
	@echo "Creating .env file if it doesn't exist..."
	@if [ ! -f .env ]; then \
		cp .env.example .env 2>/dev/null || echo "# Environment variables" > .env; \
		echo "DATABASE_URL=postgresql://user:password@localhost/dbname?sslmode=disable" >> .env; \
		echo "BOT_TOKEN=your_telegram_bot_token" >> .env; \
		echo "Please update .env with your actual configuration"; \
	fi
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run migrations in development
dev-migrate: migrate-up

# Run application in development mode
dev-run:
	@echo "Running in development mode..."
	go run cmd/bot/main.go

# Docker commands (if you use Docker)
.PHONY: docker-build docker-run docker-migrate

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t chatbot .

# Run with Docker
docker-run:
	@echo "Running with Docker..."
	docker run --env-file .env chatbot

# Run migrations with Docker
docker-migrate:
	@echo "Running migrations with Docker..."
	docker run --env-file .env chatbot go run cmd/migrate/main.go -action up

# Default target
.DEFAULT_GOAL := migrate-help