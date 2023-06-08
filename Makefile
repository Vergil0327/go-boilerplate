.PHONY: start dev build analysis

AIR = ~/go/bin/air
WIRE = ~/go/bin/wire

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = 0.0.1

APP 			= boilerplate
SERVER_BIN  	= ./cmd/${APP}/${APP}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

TARGET = ./cmd/${APP}

all: start

build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN) ./cmd/${APP}

analysis:
	@go build -gcflags="-m" $(TARGET)

# Live Reload
dev:
	${AIR} -c .air.conf

start:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" ./main.go --config ./configs/config.local.toml --version $(RELEASE_VERSION)

# Dependency Injection
wire:
	${WIRE} gen ./internal/app