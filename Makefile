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

# Update golangci-lint
lint-update:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	golangci-lint --version

