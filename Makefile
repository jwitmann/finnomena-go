# Finnomena-Go Makefile
# Build, test, and development tasks

.PHONY: all test test-race test-integration lint format format-check complexity duplicates vulncheck check clean

# Default target
all: check

# =============================================================================
# Testing
# =============================================================================

test:
	go test -short ./...

test-race:
	go test -race -short ./...

test-integration:
	go test -v -run Integration ./...

# =============================================================================
# Code Quality
# =============================================================================

lint: lint-golangci lint-staticcheck lint-vet complexity duplicates

vulncheck:
	govulncheck ./...

lint-golangci:
	golangci-lint run ./...

lint-staticcheck:
	staticcheck ./...

lint-vet:
	go vet ./...

format:
	gofumpt -w .

format-check:
	@test -z "$$(gofumpt -l .)" || (echo "Formatting issues found. Run 'make format' to fix."; exit 1)

complexity:
	gocyclo -over 15 *.go models/*.go

duplicates:
	dupl -t 100 *.go models/*.go 2>/dev/null || true

check: format-check lint test
	@echo "✅ All checks passed"

# =============================================================================
# Clean
# =============================================================================

clean:
	go clean -cache
	rm -rf bin/
