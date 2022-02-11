GO=go
GOFLAGS=-race
DEV_BIN=bin
COV_PROFILE=cp.out

.DEFAULT_GOAL := build

.PHONY: fmt
	$(GO) fmt ./...

.PHONY: vet
vet: fmt
	$(GO) vet ./...

.PHONY: lint
lint: vet
	golint -set_exit_status=1 ./...

.PHONY: test
test: lint
	$(GO) clean -testcache
	$(GO) test ./... -v

.PHONY: install
install: test
	$(GO) install ./...

mkbin:
	mkdir -p bin

.PHONY: build
build: test mkdir
	$(GO) build $(GOFLAGS) -o $(DEV_BIN)/gsnip cmd/gsnip/main.go
	$(GO) build $(GOFLAGS) -o $(DEV_BIN)/gsnipd cmd/gsnipd/main.go

.PHONY: cover
cover:
	$(GO) test -coverprofile=$(COV_PROFILE) -covermode=atomic ./...
	$(GO) tool cover -html=$(COV_PROFILE)

.PHONY: clean
clean:
	$(GO) clean github.com/mdm-code/gsnip/...
	$(GO) clean -testcache
	rm -f $(COV_PROFILE)
