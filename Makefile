.PHONY: build test install hello

build:
	go build -o alex-runner ./cmd/alex-runner

test:
	go test ./internal/... -v

install:
	go install ./cmd/alex-runner

demo-generate:
	vhs cassette.tape

demo-record:
	vhs record > cassette.tape.tmp
	@echo "Output ./public/images/demo.webp" > cassette.tape
	@echo "Set TypingSpeed 0.1s" >> cassette.tape
	@echo "" >> cassette.tape
	@cat cassette.tape.tmp >> cassette.tape
	@rm cassette.tape.tmp

hello:
	@echo "👋 Hello from Makefile!"
