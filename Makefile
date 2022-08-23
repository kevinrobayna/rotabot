PROG = bin/app
MODULE = github.com/kevinrobayna/rotabot
GIT_SHA = $(shell git rev-parse --short HEAD)-dev
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_COMMAND = CGO_ENABLED=0 go build -ldflags "-X 'main.Version=$(GIT_SHA)' -X 'main.Date=$(DATE)'"
LINT_COMMAND = golangci-lint run
UNAME_S := $(shell uname -s)

.PHONY: clean
clean:
	rm -rvf $(PROG) $(PROG:%=%.linux_amd64) $(PROG:%=%.darwin_amd64)

.PHONY: build
build: $(PROG)

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
	go run gotest.tools/gotestsum \
		--format short-verbose \
		--packages="./..." \
		--rerun-fails=3 \
		-- -coverprofile=cover.out

.PHONY: run
run:
	$(PROG)

.PHONY: check
check:
	go mod verify
	go mod tidy
	go vet ./...

.PHONY: lint
lint:
	$(LINT_COMMAND)

.PHONY: lint-fix
lint-fix:
	$(LINT_COMMAND) --fix
