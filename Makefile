SHELL := /bin/bash

export APP = github.com/bkyoung/maxwell
export GOVER ?= "1.16"
export BUILD ?= "0"
export COMMIT ?= `git rev-parse --short HEAD`
export DATE ?= `date +%s`
export VERSION ?= "dev"
export EXECUTABLE_NAME ?= `basename ${APP}`

all: clean darwin

clean:
	rm -f ${EXECUTABLE_NAME}

darwin:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
	go build -a -tags netgo \
	-ldflags "-w -extldflags '-static' \
	-X '${APP}/cmd.executableName=${EXECUTABLE_NAME}' \
	-X '${APP}/cmd.version=${VERSION}'" \
	-o ./${EXECUTABLE_NAME}

linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -a -tags netgo \
	-ldflags "-w -extldflags '-static' \
	-X '${APP}/cmd.executableName=${EXECUTABLE_NAME}' \
	-X '${APP}/cmd.version=${VERSION}'" \
	-o ./${EXECUTABLE_NAME}

deps-upgrade:
	go get -u -t -d -v ./...
	go mod vendor
	go mod tidy

release:
	rm -rf dist
	goreleaser --rm-dist

test:
	go test -v ./...

tidy:
	go mod tidy
	go mod vendor
