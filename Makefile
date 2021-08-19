SHELL := /bin/bash

export MODULE = `go list`
export GOVER ?= $(go version | awk '{print $3}' | sed -E 's/go([[:digit:]]\.[[:digit:]]{2}).*/\1/')
export BUILD ?= "0"
export COMMIT ?= `git rev-parse --short HEAD`
export DATE ?= `date +%s`
export VERSION ?= "dev"
export EXECUTABLE_NAME ?= `basename $(MODULE)`

all: clean darwin

clean:
	rm -f $(EXECUTABLE_NAME)

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
