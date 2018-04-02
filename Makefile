.PHONY: build tests

VERSION = $(shell git describe --always --long --dirty)
TAG = devel
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

# Variables to choose cross-compile target
GOOS:=linux
GOARCH:=amd64
CN_EXTENSION:=

build: clean prepare
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -i -ldflags="-X main.version=$(VERSION) -X main.tag=$(TAG) -X main.branch=$(BRANCH)"
	mv cn$(CN_EXTENSION) cn-$(TAG)-$(VERSION)-$(GOOS)-$(GOARCH)$(CN_EXTENSION)
	ln -sf "cn-$(TAG)-$(VERSION)-$(GOOS)-$(GOARCH)$(CN_EXTENSION)" cn$(CN_EXTENSION)

prepare:
	go get github.com/docker/docker/api
	go get github.com/docker/docker/client
	go get github.com/inconshreveable/mousetrap
	go get github.com/spf13/cobra
	go get github.com/jmoiron/jsonq

darwin:
	make GOOS=darwin

linux:
	make GOOS=linux

windows:
	make GOOS=windows CN_EXTENSION=".exe"

tests:
	tests/functional-tests.sh

release: darwin linux windows

clean:
	rm -f cn &>/dev/null || true
	rm -f cn.exe &>/dev/null || true

clean-all: clean
	rm -f cn-* &>/dev/null || true
