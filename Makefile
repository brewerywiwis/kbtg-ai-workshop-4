# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server
BINARY_UNIX=$(BINARY_NAME)_unix

# Main targets
all: test build

.PHONY: build
build:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) ./cmd/server

.PHONY: test
test:
	$(GOTEST) -v ./internal/...

.PHONY: test-unit
test-unit:
	$(GOTEST) -v ./internal/domain ./internal/service

.PHONY: test-coverage
test-coverage:
	$(GOTEST) -cover ./internal/...

.PHONY: test-coverage-html
test-coverage-html:
	$(GOTEST) -coverprofile=coverage.out ./internal/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

.PHONY: run
run: build
	./bin/$(BINARY_NAME)

.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) verify

.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Cross compilation
.PHONY: build-linux
build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./bin/$(BINARY_UNIX) ./cmd/server

.PHONY: build-docker
build-docker:
	docker build -t lbk-points-api .

# Development targets
.PHONY: dev
dev:
	$(GOCMD) run ./cmd/server

.PHONY: watch
watch:
	air -c .air.toml

# Database targets
.PHONY: db-reset
db-reset:
	rm -f users.db
	$(MAKE) run

.PHONY: lint
lint:
	golangci-lint run

.PHONY: format
format:
	$(GOCMD) fmt ./...

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application binary"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  run         - Build and run the application"
	@echo "  dev         - Run the application in development mode"
	@echo "  deps        - Download and verify dependencies"
	@echo "  tidy        - Tidy go modules"
	@echo "  lint        - Run linter"
	@echo "  format      - Format code"
	@echo "  db-reset    - Reset database and start fresh"