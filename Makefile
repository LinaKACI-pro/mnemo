# Makefile — mnemo
SHELL := /usr/bin/env bash

PKG       := ./...
BIN_API   := mnemo-api
BIN_CLI   := mnemo
CMD_API   := ./cmd/mnemo-api
CMD_CLI   := ./cmd/mnemo
BUILD_DIR := ./dist

# Versions pinned
GOLANGCI_LINT_VERSION := v1.60.3
GOFUMPT_VERSION       := v0.6.0
LATEST   := latest

GOFLAGS   := -trimpath
LDFLAGS   := -s -w

.DEFAULT_GOAL := help
.PHONY: help tools tidy fmt lint test cover gen build run run-api run-cli clean

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

tools:
	@echo "==> Installing tools"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(LATEST)
	@go install mvdan.cc/gofumpt@$(LATEST)
	@go install golang.org/x/vuln/cmd/govulncheck@$(LATEST)

tidy:
	go mod tidy

fmt:
	@gofumpt -l -w .
	@echo "✅ formatted"

lint:
	@golangci-lint run ./...
	@govulncheck ./...
	@echo "✅ lint & vuln checks passed"

test:
	@go test -count=1 $(PKG)
	@echo "✅ tests passed"

cover:
	@go test -race -covermode=atomic -coverprofile=coverage.out $(PKG)
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ coverage report: coverage.html"

gen:
	@go generate ./...
	@echo "✅ code generated"

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(BIN_CLI) $(CMD_CLI)
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(BIN_API) $(CMD_API)
	@echo "Binaries in $(BUILD_DIR)/"

run: run-api

run-api:
	go run $(CMD_API)

run-cli:
	go run $(CMD_CLI)

clean:
	@rm -rf $(BUILD_DIR) ./coverage.out ./coverage.html
