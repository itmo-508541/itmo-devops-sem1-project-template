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

.PHONY: compose
compose:
	@mkdir --parent ./build; \
	docker compose -f docker-compose.deploy.yml config > ./build/deploy.yml

.PHONY: migrate
migrate:
	@go run ./cmd/main.go migrate

.PHONY: start-server
start-server:
	@go run ./cmd/main.go start-server
