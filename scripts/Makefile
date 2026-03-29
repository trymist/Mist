.PHONY: all build build-cli build-server clean install help

all: build

build: build-cli

build-cli:
	@echo "Building Mist CLI..."
	cd cli && go build -o ../bin/mist-cli .

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

install: build
	@echo "Installing Mist CLI..."
	sudo cp bin/mist-cli /usr/local/bin/mist-cli
	@echo "Installation complete!"

# Help target
help:
	@echo "Mist Makefile targets:"
	@echo "  make build        - Build CLI (default)"
