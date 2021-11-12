GO=go
GOFLAGS=-race
DEV_BIN=bin

all: test build clean

mkbin:
	mkdir -p bin

.PHONY: build
build:
	go build $(GOFLAGS) -o $(DEV_BIN)/gsnip cmd/gsnip/main.go
	go build $(GOFLAGS) -o $(DEV_BIN)/gsnipd cmd/gsnipd/main.go

.PHONY: run
run: build
	./$(DEV_BIN)/gsnipd &>/dev/null &
	echo pprog | ./$(DEV_BIN)/gsnip

.PHONY: test
test:
	$(GO) clean -testcache
	$(GO) test ./... -v

.PHONY: install
install: test
	go install ./...

.PHONY: cover
cover:
	go test -coverprofile=cp.out ./...
	go tool cover -html=cp.out

.PHONY: clean
clean:
	$(GO) clean -testcache
	rm -f cp.out
