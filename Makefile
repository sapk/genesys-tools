APP_NAME=genesys-tools
APP_VERSION=$(shell git describe --tags --abbrev=0)
APP_PACKAGE=github.com/sapk/go-genesys

GIT_HASH=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
DATE := $(shell date -u '+%Y-%m-%d-%H%M-UTC' | base64)

LDFLAGS = \
  -s -w \
-X $(APP_PACKAGE)/cmd.Version=$(APP_VERSION) -X $(APP_PACKAGE)/cmd.Branch=$(GIT_BRANCH) -X $(APP_PACKAGE)/cmd.Commit=$(GIT_HASH) -X $(APP_PACKAGE)/cmd.BuildTime=$(DATE)

ERROR_COLOR=\033[31;01m
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
WARN_COLOR=\033[33;01m

all: build compress done

build: clean format compile
#build: deps clean format compile

release: clean format
	gox -ldflags "$(LDFLAGS)" -output="build/$(APP_NAME)-{{.OS}}-{{.Arch}}"
	@upx build/$(APP_NAME)-* || true

deps:
	@echo -e "$(OK_COLOR)==> Installing dependencies ...$(NO_COLOR)"
	@go get -u -v github.com/golang/dep/cmd/dep #Vendoring
	@go get -u -v github.com/mitchellh/gox #Build tool
	@$(GOPATH)/bin/dep ensure

clean:
	@if [ -x $(APP_NAME) ]; then rm $(APP_NAME); fi
	@if [ -d build ]; then rm -R build; fi
	@go clean ./...

format:
	@echo -e "$(OK_COLOR)==> Formatting...$(NO_COLOR)"
	go fmt ./...

compile:
	@echo -e "$(OK_COLOR)==> Building...$(NO_COLOR)"
	go build -v -ldflags "$(LDFLAGS)" -o $(APP_NAME)

compress:
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute $(APP_NAME) || upx-ucl --brute $(APP_NAME) || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

done:
	@echo -e "$(OK_COLOR)==> Done.$(NO_COLOR)"
