build:
	go build -o bin/gsnip main.go

run: build
	./bin/gsnip

all: run

PHONY: build run
