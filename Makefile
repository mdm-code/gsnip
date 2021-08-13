build:
	go build -o bin/gsnip main.go

run: build
	./bin/gsnip

test:
	go test ./... -v

all: test build

PHONY: build run test
