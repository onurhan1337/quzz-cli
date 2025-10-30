.PHONY: build install clean test help run

BINARY_NAME=quzz
INSTALL_PATH=/usr/local/bin

help:
	@echo "Quzz CLI - Build targets:"
	@echo "  make build      - Build the binary"
	@echo "  make install    - Build and install to $(INSTALL_PATH)"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make run        - Build and run with example"
	@echo "  make help       - Show this help message"

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) -ldflags="-s -w" .
	@echo "Build complete: ./$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installation complete. Run 'quzz --help' to get started."

clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f quzz.config.js quzz.config.ts
	@echo "Clean complete."

test:
	@echo "Running tests..."
	@go test -v ./...

run: build
	@echo "Running $(BINARY_NAME) with example traces..."
	@./$(BINARY_NAME) visualize examples/traces.json --limit 5
