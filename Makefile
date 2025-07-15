# Makefile untuk Social Media API

# Variables
APP_NAME=social-media-api
DB_CONTAINER=social_media_db
DB_NAME=social_media

# Build aplikasi
build:
	go build -o bin/$(APP_NAME) .

# Run aplikasi
run:
	go run .

# Run dengan hot reload (butuh air package)
dev:
	air

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Start PostgreSQL dengan Docker
db-start:
	docker-compose up -d postgres

# Stop PostgreSQL
db-stop:
	docker-compose down

# Connect ke PostgreSQL
db-connect:
	docker exec -it $(DB_CONTAINER) psql -U admin -d $(DB_NAME)

# Reset database
db-reset:
	docker-compose down -v
	docker-compose up -d postgres

# Clean up
clean:
	rm -rf bin/
	go clean

# Install air untuk hot reload
install-air:
	go install github.com/cosmtrek/air@latest

# Setup development environment
setup: deps db-start
	@echo "Development environment setup complete!"
	@echo "Run 'make run' to start the application"

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  db-start      - Start PostgreSQL with Docker"
	@echo "  db-stop       - Stop PostgreSQL"
	@echo "  db-connect    - Connect to PostgreSQL"
	@echo "  db-reset      - Reset database"
	@echo "  clean         - Clean up build files"
	@echo "  install-air   - Install air for hot reload"
	@echo "  setup         - Setup development environment"
	@echo "  help          - Show this help message"

.PHONY: build run dev deps fmt test test-coverage db-start db-stop db-connect db-reset clean install-air setup help
