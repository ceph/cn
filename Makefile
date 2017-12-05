.PHONY: cn

VERSION = $(shell git describe --always --long --dirty)
TAG = devel
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

cn: clean
	go build -i -ldflags="-X main.version=$(VERSION) -X main.tag=$(TAG) -X main.branch=$(BRANCH)"
	mv cn cn-$(TAG)-$(VERSION)
	ln -sf "cn-$(TAG)-$(VERSION)" cn

clean:
	rm -f cn &>/dev/null || true

clean-all: clean
	rm -f cn-* &>/dev/null || true
