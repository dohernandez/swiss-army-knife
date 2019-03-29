# detecting GOPATH and removing trailing "/" if any
GOPATH = $(realpath $(shell go env GOPATH))
IMPORT_PATH = $(subst $(GOPATH)/src/,,$(realpath $(shell pwd)))

ROOT_PATH ?= $(PWD)
RESOURCES_SCRIPTS_PATH ?= $(ROOT_PATH)/resources/scripts

branch = $(shell git symbolic-ref HEAD 2>/dev/null)
VERSION ?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
revision = $(shell git log -1 --pretty=format:"%H")
build_user = $(USER)
build_date = $(shell date +%FT%T%Z)

VERSION_PKG = $(IMPORT_PATH)/pkg/version
export LDFLAGS = -X $(VERSION_PKG).version=$(VERSION) -X $(VERSION_PKG).branch=$(branch) -X $(VERSION_PKG).revision=$(revision) -X $(VERSION_PKG).buildUser=$(build_user) -X $(VERSION_PKG).buildDate=$(build_date)

BUILD_DIR ?= bin
BINARY_NAME ?= swiss-army-knife

# Filters variables
CFLAGS=-g
export CFLAGS

## Ensure dependencies according to toml file
deps:
	@echo ">> ensuring dependencies"
	@test -s $(GOPATH)/bin/dep || GOBIN=$(GOPATH)/bin go get -u github.com/golang/dep/cmd/dep
	@$(GOPATH)/bin/dep ensure
	@git add ${ROOT_PATH}/Gopkg.lock

## Ensure dependencies according to lock file
deps-vendor:
	@echo ">> ensuring dependencies"
	@test -s $(GOPATH)/bin/dep || GOBIN=$(GOPATH)/bin go get -u github.com/golang/dep/cmd/dep
	@$(GOPATH)/bin/dep ensure --vendor-only
	@git add ${ROOT_PATH}/Gopkg.lock

## Build binary
build:
	@echo ">> building binary"
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) cmd/$(BINARY_NAME)/*

## Run application with CompileDaemon (automatic rebuild on code change)
run-compile-daemon:
	@test -s $(shell go env GOPATH)/bin/CompileDaemon || (echo ">> installing CompileDaemon" && go get -u github.com/githubnemo/CompileDaemon)
	@echo ">> running app with CompileDaemon"
	@$(shell go env GOPATH)/bin/CompileDaemon -exclude-dir=vendor -color=true -build='make build' -command='$(BUILD_DIR)/${BINARY_NAME}' -graceful-kill

## Check with golangci-lint
lint:
	@$(RESOURCES_SCRIPTS_PATH)/check-lint.sh

## Apply goimports and gofmt
fix-lint:
	@$(RESOURCES_SCRIPTS_PATH)/fix-style.sh

## Run unit tests
test:
	@echo ">> unit test"
	@test -s $(GOPATH)/bin/overalls || GOBIN=$(GOPATH)/bin go get -u github.com/go-playground/overalls
	@$(GOPATH)/bin/overalls -project=${IMPORT_PATH} -covermode=atomic -- -race


.PHONY: deps deps-vendor build run-compile-daemon lint fix-lint test help

.DEFAULT_GOAL := help
HELP_SECTION_WIDTH="      "
HELP_DESC_WIDTH="                       "
help:
	@printf "$(BINARY_NAME) routine operations\n\n";
	@awk '{ \
			if ($$0 ~ /^.PHONY: [a-zA-Z\-\_0-9]+$$/) { \
				helpCommand = substr($$0, index($$0, ":") + 2); \
				if (helpMessage) { \
					printf "  \033[32m%-20s\033[0m %s\n", \
						helpCommand, helpMessage; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^[a-zA-Z\-\_0-9.]+:/) { \
				helpCommand = substr($$0, 0, index($$0, ":")); \
				if (helpMessage) { \
					printf "  \033[32m%-20s\033[0m %s\n", \
						helpCommand, helpMessage; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^##/) { \
				if (helpMessage) { \
					helpMessage = helpMessage"\n"${HELP_DESC_WIDTH}substr($$0, 3); \
				} else { \
					helpMessage = substr($$0, 3); \
				} \
			} else { \
				if (helpMessage) { \
					print "\n"${HELP_SECTION_WIDTH}helpMessage"\n" \
				} \
				helpMessage = ""; \
			} \
		}' \
		$(MAKEFILE_LIST)
	@printf "\nUsage\n";
	@printf "  make <flags> [options]\n";
