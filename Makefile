SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.PHONY: deps
deps:
	@go mod tidy; \
	go mod download; \
	go mod verify

.PHONY: fmt
fmt:
	@golangci-lint fmt -c .golangci.yaml ./...

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@go test ./internal/...

.PHONY: build
build:
	@mkdir --parent ./build; \
	go build -o ./build/app ./cmd/main.go

.PHONY: dist
dist:
	@mkdir --parent ./build; \
	go build -tags dist -o ./build/app ./cmd/main.go

.PHONY: migrate
migrate:
	@go run ./cmd/main.go migrate

.PHONY: start-server
start-server:
	@go run ./cmd/main.go start-server
