.PHONY: build test lint lint-self clean deps

PROGRAM := namedreturns

# Default target
all: build

# Install dependencies
deps:
	go mod tidy -v

# Build the binary
build: deps
	CGO_ENABLED=0 go build -o $(PROGRAM) .

# Run tests
test: deps
	CGO_ENABLED=1 go test -race -cover ./...

# Run golangci-lint (excluding namedreturns for local dev)
lint:
	golangci-lint run ./... --timeout=30m --config=.golangci.local.yml
	go mod tidy

# Run namedreturns linter on our own codebase
lint-self: build
	./$(PROGRAM) ./...

# Clean build artifacts
clean:
	rm -f $(PROGRAM)

# Show compatibility info
compatibility-info:
	@echo "=== namedreturns Compatibility Info ==="
	@echo "Built with: $(shell go version)"
	@echo "Binary info: $(shell go version -m ./$(PROGRAM) 2>/dev/null | head -1 || echo 'Binary not found - run make build')"
	@echo "Supports analyzing: Go 1.21.0+ codebases"
	@echo "Recommendation: Rebuild with your current Go version for best compatibility"

# Rebuild for current Go version
rebuild: clean build
	@echo "namedreturns rebuilt with $(shell go version)"
	@make compatibility-info

# Update golangci-lint
lint-update:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	golangci-lint --version

# Help target
help:
	@echo "Available targets:"
	@echo "  build              - Build the namedreturns binary"
	@echo "  test               - Run tests with race detection"
	@echo "  lint               - Run golangci-lint"
	@echo "  lint-self          - Run namedreturns on its own codebase"
	@echo "  clean              - Remove built binaries"
	@echo "  deps               - Install/update dependencies"
	@echo "  rebuild            - Clean and rebuild for current Go version"
	@echo "  compatibility-info - Show Go version compatibility information"
	@echo "  lint-update        - Update golangci-lint"
	@echo "  help               - Show this help"

