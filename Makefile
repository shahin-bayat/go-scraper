BUILD_DIR = ./bin
SCRAPE_APP_NAME = go_scraper
PHOTOGEN_APP_NAME = go_photogen
SCRAPE_FILE = ./cmd/scrape/scrape.go
PHOTOGEN_FILE = ./cmd/photogen/photogen.go

build-scrape:
	@go build -o $(BUILD_DIR)/$(SCRAPE_APP_NAME) $(SCRAPE_FILE)

run-scraper: build-scrape
	@$(BUILD_DIR)/$(APP_NAME)

build-photogen:
	@go build -o $(BUILD_DIR)/$(PHOTOGEN_APP_NAME) $(PHOTOGEN_FILE)

run-photogen: build-photogen
	@$(BUILD_DIR)/$(PHOTOGEN_APP_NAME)

test:
	@go test -v ./...
