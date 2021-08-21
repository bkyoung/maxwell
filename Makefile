SHELL := /bin/bash

export MODULE := github.com/bkyoung/maxwell
GOVER ?= $(go version | awk '{print $3}' | sed -E 's/go([[:digit:]]\.[[:digit:]]{2}).*/\1/')
BUILD ?= "0"
COMMIT ?= `git rev-parse --short HEAD`
DATE ?= `date +%s`
VERSION ?= dev
EXECUTABLE_NAME ?= `basename $${MODULE}`

all: clean linux

clean:
	rm -f ./$(EXECUTABLE_NAME)

darwin:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
	go build -a -tags netgo \
	-ldflags "-w -extldflags '-static' \
	-X '$(MODULE)/cmd.executableName=$(EXECUTABLE_NAME)' \
	-X '$(MODULE)/cmd.version=$(VERSION)'" \
	-o ./$(EXECUTABLE_NAME)

linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -a -tags netgo \
	-ldflags "-w -extldflags '-static' \
	-X '$(MODULE)/cmd.executableName=$(EXECUTABLE_NAME)' \
	-X '$(MODULE)/cmd.version=$(VERSION)'" \
	-o ./$(EXECUTABLE_NAME)

release:
	rm -rf dist
	goreleaser --rm-dist

test:
	go test -v ./...

tidy:
	go mod tidy
