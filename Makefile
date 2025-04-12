BINARY_NAME=pvz_server

MAIN_FILE=./cmd/apiserver/main.go

BUILD_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run:
	@echo "Running pvz_server..."
	@ENV_FILE=$(ENV_FILE) go run $(MAIN_FILE)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"