.DEFAULT_GOAL := build

.PHONY: lint
lint: install fmt
	@go vet ./...
	@staticcheck ./...

.PHONY: fmt
fmt: install
	@goimports -l -w .

build:
	@go build -ldflags="-s -w" -trimpath

.PHONY: test
test:
	@go test -v ./...

install:
	@go install golang.org/x/tools/cmd/goimports@v0.12.0
	@go install honnef.co/go/tools/cmd/staticcheck@v0.4.5
