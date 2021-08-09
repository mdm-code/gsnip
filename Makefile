build: test
	go build -o bin/gsnip main.go

run: build
	./bin/gsnip

test:
	go test ./... -v

all: run

PHONY: build run test
