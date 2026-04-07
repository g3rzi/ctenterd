# ─── Configuration ────────────────────────────────────────────────────────────

BINARY     := ctenterd
MODULE     := $(shell go list -m 2>/dev/null || echo github.com/g3rzi/ctenterd)
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo v0.1.0)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo unknown)

BIN_DIR        := bin
STATIC_DIR     := $(BIN_DIR)/static

LDFLAGS := -X main.version=$(VERSION) \
           -X main.buildTime=$(BUILD_TIME) \
           -X main.gitCommit=$(GIT_COMMIT)

STATIC_LDFLAGS := $(LDFLAGS) -extldflags "-static"

# ─── Targets ──────────────────────────────────────────────────────────────────

.PHONY: all build static clean help

## all: build both dynamic and static binaries
all: build static

## build: compile a dynamically-linked binary into bin/
build: $(BIN_DIR)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) .
	@echo "Built $(BIN_DIR)/$(BINARY)"

## static: compile a fully static binary into bin/static/
static: $(STATIC_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-a \
		-tags netgo \
		-ldflags "$(STATIC_LDFLAGS)" \
		-o $(STATIC_DIR)/$(BINARY) .
	@echo "Built $(STATIC_DIR)/$(BINARY)"

## clean: remove the bin/ directory
clean:
	rm -rf $(BIN_DIR)
	@echo "Cleaned $(BIN_DIR)/"

## help: print this help message
help:
	@echo "Usage: make <target>"
	@echo ""
	@grep -E '^##' Makefile | sed 's/## /  /'

# ─── Directory creation ───────────────────────────────────────────────────────

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(STATIC_DIR):
	mkdir -p $(STATIC_DIR)
