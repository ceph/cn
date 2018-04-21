.PHONY: build tests

VERSION = $(shell git describe --always --long --dirty)
TAG = devel
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

# Variables to choose cross-compile target
GOOS:=linux
GOARCH:=amd64
CN_EXTENSION:=

build: check clean prepare
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -i -ldflags="-X main.version=$(VERSION) -X main.tag=$(TAG) -X main.branch=$(BRANCH)"
	mv cn$(CN_EXTENSION) cn-$(TAG)-$(VERSION)-$(GOOS)-$(GOARCH)$(CN_EXTENSION)
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
	go get github.com/docker/docker/api
	go get github.com/docker/docker/client
	go get github.com/inconshreveable/mousetrap
	go get github.com/spf13/cobra
	go get github.com/jmoiron/jsonq
	go get github.com/apcera/termtables
	go get golang.org/x/sys/unix

darwin:
	make GOOS=darwin GOARCH:=amd64

windows:
	make GOOS=windows CN_EXTENSION=".exe" GOARCH:=amd64

linux-%:
	make GOOS=linux GOARCH:=$*

tests:
	tests/functional-tests.sh

release: darwin windows linux-amd64 linux-arm64

clean:
	rm -f cn &>/dev/null || true
	rm -f cn.exe &>/dev/null || true

clean-all: clean
	rm -f cn-* &>/dev/null || true
