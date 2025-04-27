.DEFAULT_GOAL := build

.PHONY: build
build: generate
	@go build -ldflags="-s -w" -trimpath ./cmd/api-server/main.go
	@go build -ldflags="-s -w" -trimpath ./cmd/delete/main.go
	@go build -ldflags="-s -w" -trimpath ./cmd/standard-writer/main.go

.PHONY: test
test: generate
	@go test -v ./...

.PHONY: dry-lint
dry-lint:
	@golangci-lint run

.PHONY: lint
lint:
	@golangci-lint run --fix

.PHONY: coverage
coverage: generate
	@go test -cover ./... -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html

generate:
	@go generate ./...
