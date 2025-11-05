.PHONY: build test install hello

build:
	go build -o alex-runner ./cmd/alex-runner

test:
	go test ./internal/... -v

install:
	go install ./cmd/alex-runner

hello:
	@echo "👋 Hello from Makefile!"
