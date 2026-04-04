.PHONY: all build build-cli build-server clean install help

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)
BIN_DIR := $(ROOT_DIR)/bin
SERVER_DIR := $(ROOT_DIR)/apps/server

all: build

build: build-cli

build-cli:
	@echo "Building Mist CLI..."
	@mkdir -p $(BIN_DIR)
	cd $(SERVER_DIR) && go build -o $(BIN_DIR)/mist-cli ./cmd/cli

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)

install: build
	@echo "Installing Mist CLI..."
	sudo cp $(BIN_DIR)/mist-cli /usr/local/bin/mist-cli
	@echo "Installation complete!"

# Help target
help:
	@echo "Mist Makefile targets:"
	@echo "  make build        - Build CLI (default)"
