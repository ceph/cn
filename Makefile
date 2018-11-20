.PHONY: build tests

VERSION = $(shell git describe --always --long --dirty)
TAG = devel
TARGET_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)

# Variables to choose cross-compile target
GOOS:=linux
GOARCH:=amd64
CN_EXTENSION:=

build: check clean prepare
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -i -ldflags="-X main.version=$(VERSION) -X main.tag=$(TAG) -X main.branch=$(TARGET_BRANCH)" -o cn-$(TAG)-$(VERSION)-$(GOOS)-$(GOARCH)$(CN_EXTENSION) main.go
	ln -sf "cn-$(TAG)-$(VERSION)-$(GOOS)-$(GOARCH)$(CN_EXTENSION)" cn$(CN_EXTENSION)

check:
ifeq ("$(GOPATH)","")
	@echo "GOPATH variable must be defined"
	@exit 1
endif
ifneq ("$(shell pwd)","$(GOPATH)/src/github.com/ceph/cn")
	@echo "You are in $(shell pwd) !"
	@echo "Please go in $(GOPATH)/src/github.com/ceph/cn to build"
	@exit 1
endif

prepare:
	dep ensure

darwin:
	make GOOS=darwin GOARCH:=amd64

linux-%:
	make GOOS=linux GOARCH:=$*

tests:
	tests/functional-tests.sh

release: darwin linux-amd64 linux-arm64

clean:
	rm -f cn$(CN_EXTENSION) cn &>/dev/null || true

clean-all: clean
	rm -f cn-* &>/dev/null || true
