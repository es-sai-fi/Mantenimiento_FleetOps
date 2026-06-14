# =============================================================================
# Makefile — FleetOps Maintenance Microservice
# =============================================================================

.PHONY: all build run test test-coverage lint fmt clean docker-up docker-down migrate-up migrate-down

# Variables
BINARY_NAME=maintenance-service
BINARY_PATH=bin/$(BINARY_NAME)
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
COVERAGE_DIR=coverage
COVERAGE_FILE=$(COVERAGE_DIR)/coverage.out
COVERAGE_HTML=$(COVERAGE_DIR)/coverage.html

# Default target
all: lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GO) build -ldflags="-s -w" -o $(BINARY_PATH) ./cmd/server

# Run locally (requires DATABASE_URL env var)
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_PATH)

# Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./internal/...

# Run tests with coverage report (Rule R3)
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic \
		./internal/domain/... \
		./internal/service/... \
		./internal/handler/...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	$(GO) tool cover -func=$(COVERAGE_FILE)
	@echo ""
	@echo "Coverage report: $(COVERAGE_HTML)"

# Lint with golangci-lint (ADR-13)
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Format with gofumpt (ADR-13)
fmt:
	@echo "Formatting code..."
	gofumpt -w .

# Vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Generate mocks with mockery
mocks:
	@echo "Generating mocks..."
	mockery --all --dir=internal/port --output=internal/mocks --outpkg=mocks

# Docker Compose up (Rule R4)
docker-up:
	@echo "Starting Docker services..."
	docker compose up --build -d

# Docker Compose down
docker-down:
	@echo "Stopping Docker services..."
	docker compose down -v

# Database migrations — up (Archetype Convention: golang-migrate)
migrate-up:
	@echo "Running migrations up..."
	migrate -path migrations -database "$${DATABASE_URL}" up

# Database migrations — down
migrate-down:
	@echo "Running migrations down..."
	migrate -path migrations -database "$${DATABASE_URL}" down

# Tidy dependencies
tidy:
	$(GO) mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/ $(COVERAGE_DIR)/
