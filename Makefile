BINARY_NAME=pvz_server

MAIN_FILE=./cmd/apiserver/main.go

BUILD_DIR=bin

DB_URL=$(DATABASE_URL)


migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	go test ./internal/handlers/ -v -cover

integration_test:
	@echo "Running integration test..."
	bash -c "set -a && source .env.local && set +a && go test ./internal/handlers/integration_test -v -count=1"

run:
	@echo "Running pvz_server..."
	bash -c "set -a && source .env.local && set +a && go run $(MAIN_FILE)"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"