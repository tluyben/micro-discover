# Makefile for Discovery Service

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=discovery-service
BINARY_UNIX=$(BINARY_NAME)_unix

# Main package path
MAIN_PACKAGE=.

# Test parameters
TEST_TIMEOUT=10s

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)

test:
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)
	./$(BINARY_NAME)

deps:
	$(GOGET) github.com/gorilla/mux
	$(GOGET) github.com/mattn/go-sqlite3
	$(GOGET) golang.org/x/crypto/bcrypt

tidy:
	$(GOMOD) tidy

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PACKAGE)

docker-build:
	docker build -t discovery-service .

# Database operations
db-create:
	@echo "Creating discovery.db..."
	@touch discovery.db

db-clean:
	@echo "Cleaning up discovery.db..."
	@rm -f discovery.db

# Help target
help:
	@echo "Available targets:"
	@echo "  build        : Build the discovery service"
	@echo "  test         : Run tests"
	@echo "  clean        : Clean up build artifacts"
	@echo "  run          : Build and run the discovery service"
	@echo "  deps         : Get dependencies"
	@echo "  tidy         : Tidy up the go.mod file"
	@echo "  build-linux  : Cross-compile for Linux"
	@echo "  docker-build : Build Docker image"
	@echo "  db-create    : Create empty SQLite database"
	@echo "  db-clean     : Remove SQLite database"

.PHONY: all build test clean run deps tidy build-linux docker-build db-create db-clean help