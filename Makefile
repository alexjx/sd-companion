.PHONY: all build compile clean pages

BUILDTIME ?= $(shell date +%Y-%m-%d_%I:%M:%S)
GITCOMMIT ?= $(shell git rev-parse HEAD)
ifeq ($(CI_PIPELINE_ID),)
    BUILDNUMER := "private"
else
    BUILDNUMER := $(CI_PIPELINE_ID)
endif

LDFLAGS = -extldflags \
          -static \
          -X "main.BuildTime=$(BUILDTIME)" \
          -X "main.GitCommit=$(GITCOMMIT)" \
          -X "main.BuildNumber=$(BUILDNUMER)"

all: build

clean:
	rm -rf bin

build:
	go build -o bin/image-broswer -ldflags "$(LDFLAGS)" main.go

compile:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/image-broswer-linux-amd64 -ldflags "$(LDFLAGS)" main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/image-broswer-windows-amd64.exe -ldflags "$(LDFLAGS)" main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/image-broswer-darwin-amd64 -ldflags "$(LDFLAGS)" main.go
