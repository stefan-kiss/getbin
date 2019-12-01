VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVEND=$(GOCMD) mod vendor
GO111MODULE=on

export GO111MODULE

all: clean get vend build

build: test build-darwin build-linux build-windows

clean:
	$(GOCLEAN) ./...
	rm -rf bin/*

get:
	$(GOGET) ./...

vend:
	$(GOVEND)

test:
	$(GOTEST) ./...

build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o "bin/getbin-linux-amd64" cmd/getbin/main.go

build-windows:
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -v -o "bin/getbin-windows-amd64.exe" cmd/getbin/main.go

build-darwin:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -v -o "bin/getbin-darwin-amd64" cmd/getbin/main.go
