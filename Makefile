GO=go
GOFLAGS=-race
DEV_BIN=bin

mkbin:
	mkdir -p bin

build:
	go build $(GOFLAGS) -o $(DEV_BIN)/gsnip cmd/gsnip/main.go
	go build $(GOFLAGS) -o $(DEV_BIN)/gsnipd cmd/gsnipd/main.go

run: build
	./$(DEV_BIN)/gsnip -snippets snips

test:
	$(GO) clean -testcache
	$(GO) test ./... -v

install: test
	go install ./...

all: test build

PHONY: build run test install
