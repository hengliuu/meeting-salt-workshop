.PHONY: build run dev test clean docker-build docker-run help

# Variables
BINARY_NAME=meeting-salt-api
MAIN_PATH=./cmd/server
BUILD_DIR=./bin

# Default target
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
dev: ## Run the application in development mode
	@echo "Starting development server..."
	@air -c .air.toml

run: ## Run the application
	@echo "Running application..."
	@go run $(MAIN_PATH)/main.go

build: ## Build the application
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

##@ Testing
test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

##@ Database
migrate-up: ## Run database migrations up
	@echo "Running database migrations..."
	@go run $(MAIN_PATH)/main.go --migrate

migrate-down: ## Run database migrations down
	@echo "Rolling back database migrations..."
	@echo "Migration rollback not implemented yet"

##@ Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

##@ Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

docker-compose-up: ## Start services with docker-compose
	@echo "Starting services with docker-compose..."
	@docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	@echo "Stopping services with docker-compose..."
	@docker-compose down

##@ Code Quality
lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

tidy: ## Tidy and vendor dependencies
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod vendor

##@ Production
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

build-windows: ## Build for Windows  
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)

build-all: build build-linux build-windows ## Build for all platforms