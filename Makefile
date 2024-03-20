APP_NAME = go_scraper
BUILD_DIR = ./bin
MAIN_FILE = ./cmd/main.go

build:
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

run: build
	@$(BUILD_DIR)/$(APP_NAME)

test:
	@go test -v ./...
