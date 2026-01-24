# Vacation Tool Makefile
# Variables
BINARY_NAME=vt
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/main.go
DOCKER_IMAGE=vacation-tool
DOCKER_COMPOSE=docker-compose

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-a -installsuffix cgo

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: clean deps fmt vet test build ## Run all checks and build

## Build targets
.PHONY: build
build: ## Build the application binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_PATH)"

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(BUILD_FLAGS) -o $(BINARY_PATH)-linux $(MAIN_PATH)
	@echo "Linux build complete: $(BINARY_PATH)-linux"

.PHONY: build-docker
build-docker: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):latest .
	@echo "Docker image built: $(DOCKER_IMAGE):latest"

## Run targets
.PHONY: run
run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_PATH)

.PHONY: dev
dev: ## Run with auto-reload (requires 'air' or similar)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Install 'air' for live reload: go install github.com/cosmtrek/air@latest"; \
		$(MAKE) run; \
	fi

## Test targets
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	@$(GOTEST) -v -race -coverprofile=coverage.out ./...

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test-short
test-short: ## Run short tests only
	@$(GOTEST) -short ./...

.PHONY: bench
bench: ## Run benchmarks
	@$(GOTEST) -bench=. -benchmem ./...

## Code quality targets
.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	@$(GOFMT) ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@$(GOVET) ./...

.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## Dependency targets
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy

.PHONY: deps-upgrade
deps-upgrade: ## Upgrade dependencies
	@echo "Upgrading dependencies..."
	@$(GOGET) -u ./...
	@$(GOMOD) tidy

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	@$(GOMOD) verify

## Docker targets
.PHONY: docker-up
docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) up -d

.PHONY: docker-down
docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE) down

.PHONY: docker-restart
docker-restart: docker-down docker-up ## Restart Docker containers

.PHONY: docker-logs
docker-logs: ## Show Docker logs
	@$(DOCKER_COMPOSE) logs -f

.PHONY: docker-ps
docker-ps: ## Show running containers
	@$(DOCKER_COMPOSE) ps

.PHONY: docker-rebuild
docker-rebuild: ## Rebuild and restart containers
	@echo "Rebuilding containers..."
	@$(DOCKER_COMPOSE) down
	@$(DOCKER_COMPOSE) up -d --build

## Database targets
.PHONY: db-up
db-up: ## Start only the database
	@echo "Starting database..."
	@$(DOCKER_COMPOSE) up -d db

.PHONY: db-down
db-down: ## Stop the database
	@$(DOCKER_COMPOSE) stop db

.PHONY: db-reset
db-reset: ## Reset database (WARNING: destroys data)
	@echo "Resetting database..."
	@$(DOCKER_COMPOSE) down -v
	@$(DOCKER_COMPOSE) up -d db
	@echo "Database reset complete"

.PHONY: db-shell
db-shell: ## Open database shell
	@docker exec -it vacation-tool-db-1 mysql -u portal -ppassword123 vacation_tool

.PHONY: db-logs
db-logs: ## Show database logs
	@docker logs -f vacation-tool-db-1

## Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

.PHONY: clean-all
clean-all: clean docker-down ## Clean everything including Docker
	@echo "Cleaning Docker volumes..."
	@$(DOCKER_COMPOSE) down -v
	@docker system prune -f

## Utility targets
.PHONY: install
install: build ## Install binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

.PHONY: uninstall
uninstall: ## Uninstall binary from $GOPATH/bin
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)

.PHONY: check
check: fmt vet test ## Run all checks (fmt, vet, test)

.PHONY: ci
ci: deps check build ## Run CI pipeline locally

.PHONY: version
version: ## Show Go version
	@$(GOCMD) version

.PHONY: info
info: ## Show project info
	@echo "Binary: $(BINARY_NAME)"
	@echo "Path: $(BINARY_PATH)"
	@echo "Main: $(MAIN_PATH)"
	@echo "Docker Image: $(DOCKER_IMAGE)"