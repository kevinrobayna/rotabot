PROG = bin/app
MODULE = github.com/kevinrobayna/rotabot
GIT_SHA = $(shell git rev-parse --short HEAD)-dev
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_COMMAND = CGO_ENABLED=0 go build -ldflags "-X 'main.Version=$(GIT_SHA)' -X 'main.Date=$(DATE)'"
LINT_COMMAND = golangci-lint run

LICENSED_VERSION = 3.7.2
UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Linux)
		LICENSED_URL = https://github.com/github/licensed/releases/download/$(LICENSED_VERSION)/licensed-$(LICENSED_VERSION)-linux-x64.tar.gz
endif
ifeq ($(UNAME_S),Darwin)
		LICENSED_URL = https://github.com/github/licensed/releases/download/$(LICENSED_VERSION)/licensed-$(LICENSED_VERSION)-darwin-x64.tar.gz
endif

.PHONY: clean
clean:
	rm -rvf $(PROG) $(PROG:%=%.linux_amd64) $(PROG:%=%.darwin_amd64)

.PHONY: build
build: clean $(PROG)

.PHONY: all darwin linux
all: darwin linux
darwin: $(PROG:=.darwin_amd64)
linux: $(PROG:=.linux_amd64)

bin/%.linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD_COMMAND) -a -o $@ cmd/$*/main.go

bin/%.darwin_amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(BUILD_COMMAND) -a -o $@ cmd/$*/main.go

bin/%:
	$(BUILD_COMMAND) -o $@ cmd/$*/main.go

.PHONY: test
test:
	gotestsum --packages="./..." -- -coverprofile=cover.out

.PHONY: dev
dev: build
	reflex --sequential --decoration=fancy --config=reflex.conf

.PHONY: run
run:
	$(PROG)

.PHONY: lint
lint:
	$(LINT_COMMAND)

.PHONY: lint-fix
lint-fix:
	$(LINT_COMMAND) --fix

.PHONE: check-licenses
check-licenses: cache-licenses
	./tools/licensed status
	./scripts/check-changes.sh

.PHONE: cache-licenses
cache-licenses:
	./tools/licensed cache

.PHONY: install
install: install-deps install-licensed

.PHONY: install-deps
install-deps:
	go mod download
	go install github.com/cespare/reflex
	go install gotest.tools/gotestsum
	go install github.com/golangci/golangci-lint/cmd/golangci-lint


.PHONY: install-licensed
install-licensed:
ifndef LICENSED_URL
	$(error Could not resolve LICENSED_URL for the current OS ($(UNAME_S)))
endif
	mkdir -p tools
	curl -sSfL $(LICENSED_URL) > tools/licensed-$(LICENSED_VERSION).tar.gz
	tar xzf tools/licensed-$(LICENSED_VERSION).tar.gz --directory tools