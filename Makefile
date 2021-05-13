#!/bin/make

GOROOT:=$(shell PATH="/pkg/main/dev-lang.go/bin:$$PATH" go env GOROOT)
GO_TAG:=$(shell /bin/sh -c 'eval `$(GOROOT)/bin/go tool dist env`; echo "$${GOOS}_$${GOARCH}"')
GIT_TAG:=$(shell git rev-parse --short HEAD)
GOPATH:=$(shell $(GOROOT)/bin/go env GOPATH)
GOOS:=$(shell $(GOROOT)/bin/go env GOOS)
SOURCES:=$(shell find . -name '*.go')
ifeq ($(DATE_TAG),)
DATE_TAG:=$(shell date '+%Y%m%d%H%M%S')
endif
ifeq ($(GOOS),windows)
SUFFIX=.exe
else
SUFFIX=
endif
export DATE_TAG
export GO111MODULE=on
export CGO_ENABLED=0

.PHONY: all deps update fmt test check doc gen cov

all: shells-cli$(SUFFIX)

shells-cli$(SUFFIX): $(SOURCES)
	$(GOPATH)/bin/goimports -w -l .
	$(GOROOT)/bin/go build -v -gcflags="-N -l" -ldflags=all="-X main.GIT_TAG=$(GIT_TAG) -X main.DATE_TAG=$(DATE_TAG)"

clean:
	$(GOROOT)/bin/go clean

deps:
	$(GOROOT)/bin/go get -v .

update:
	$(GOROOT)/bin/go get -u .

fmt:
	$(GOROOT)/bin/go fmt ./...
	$(GOPATH)/bin/goimports -w -l .

test:
	$(GOROOT)/bin/go test ./...

gen:
	$(GOROOT)/bin/go generate

cov:
	$(GOROOT)/bin/go test -coverprofile=coverage.out ./...
	$(GOROOT)/bin/go tool cover -html=coverage.out -o coverage.html

check:
	@if [ ! -f $(GOPATH)/bin/gometalinter ]; then $(GOROOT)/bin/go get github.com/alecthomas/gometalinter; fi
	$(GOPATH)/bin/gometalinter ./...

doc:
	@if [ ! -f $(GOPATH)/bin/godoc ]; then $(GOROOT)/bin/go get golang.org/x/tools/cmd/godoc; fi
	$(GOPATH)/bin/godoc -v -http=:6060 -index -play

