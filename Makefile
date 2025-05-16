
CONFIG_PATH ?= config/config.yml

DB_DSN ?= postgres://user:pass@localhost:5432/limits_db?sslmode=disable

BINARY     := load_balancer
CMD        := cmd/load_balancer/main.go

BIN_DIR    := bin

.PHONY: all build run migrate-up migrate-down test docker-up docker-down

all: build

build:
	@echo "→ Building $(BINARY)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY) $(CMD)

run: build
	@echo "→ Running $(BINARY)..."
	@CONFIG_PATH=$(CONFIG_PATH) \
	DB_DSN=$(DB_DSN) \
	$(BIN_DIR)/$(BINARY)

migrate-up:
	@echo "→ Applying migrations (up)..."
	@goose -dir migrations postgres "$(DB_DSN)" up

migrate-down:
	@echo "→ Rolling back last migration..."
	@goose -dir migrations postgres "$(DB_DSN)" down

test:
	@echo "→ Running tests..."
	@go test ./... -race

docker-up:
	@echo "→ Bringing up Docker services..."
	@docker-compose up --build

docker-down:
	@echo "→ Tearing down Docker services..."
	@docker-compose down
