PROG = bin/app
MODULE = github.com/kevinrobayna/rotabot
GIT_SHA = $(shell git rev-parse --short HEAD)-dev
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_COMMAND = CGO_ENABLED=0 go build -ldflags "-X 'main.Version=$(GIT_SHA)' -X 'main.Date=$(DATE)'"
LINT_COMMAND = golangci-lint run

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
	gotestsum -- -coverprofile=cover.out

.PHONY: dev
dev: clean build
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

.PHONY: install
install: install-deps

.PHONY: install-deps
install-deps:
	go mod download
	go install github.com/cespare/reflex
	go install gotest.tools/gotestsum
