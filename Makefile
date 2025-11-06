# Semantic Search Bot Makefile

# Variables
BINARY_NAME=semantic-search-bot
MAIN_PATH=./main.go
BUILD_DIR=./bin

# Load database path from .env file (fallback to default)
DB_PATH := $(shell grep '^DATABASE_PATH=' .env 2>/dev/null | cut -d'=' -f2 | tr -d '"' || echo "./messages.db")

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run         - Run the bot (checks Ollama automatically)"
	@echo "  make build       - Build the binary"
	@echo "  make setup       - Initial setup (deps + .env file)"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts and database"
	@echo "  make clean-db    - Clean only the database file"
	@echo "  make fmt         - Format code"

# Main commands
.PHONY: run
run: deps
	@echo "🔍 Checking Ollama status..."
	@curl -s http://localhost:11434/api/tags > /dev/null || (echo "❌ Ollama not running. Start with: ollama serve" && exit 1)
	@ollama list 2>/dev/null | grep -q "all-minilm" || (echo "⚠️  Installing all-minilm model..." && ollama pull all-minilm)
	@echo "✅ Ollama ready"
	@echo "🤖 Starting bot..."
	@echo "Make sure TELEGRAM_TOKEN is set in .env file"
	go run $(MAIN_PATH)

.PHONY: build
build: deps fmt
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary created: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: setup
setup: deps
	@echo "Setting up .env file..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from template"; \
		echo ""; \
		echo "⚠️  IMPORTANT: Edit .env file and set your TELEGRAM_TOKEN"; \
		echo "   1. Get token from @BotFather on Telegram"; \
		echo "   2. Edit .env file: TELEGRAM_TOKEN=your_actual_token"; \
		echo ""; \
	else \
		echo ".env file already exists"; \
	fi
	@echo "Installing Ollama model..."
	@which ollama > /dev/null || (echo "❌ Ollama not found. Install from: https://ollama.ai/download" && exit 1)
	@curl -s http://localhost:11434/api/tags > /dev/null || (echo "❌ Ollama not running. Start with: ollama serve" && exit 1)
	@ollama pull all-minilm
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Edit .env file and set your TELEGRAM_TOKEN"
	@echo "2. Run the bot: make run"

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	@if [ -f "$(DB_PATH)" ]; then \
		echo "Removing database: $(DB_PATH)"; \
		rm -f "$(DB_PATH)"; \
	else \
		echo "Database file not found: $(DB_PATH)"; \
	fi
	go clean
	@echo "Cleaned build artifacts and database"

.PHONY: clean-db
clean-db:
	@if [ -f "$(DB_PATH)" ]; then \
		echo "Removing database: $(DB_PATH)"; \
		rm -f "$(DB_PATH)"; \
		echo "Database cleaned"; \
	else \
		echo "Database file not found: $(DB_PATH)"; \
	fi