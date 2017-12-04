.PHONY: cn

VERSION = $(shell git describe --always --long --dirty)
TAG = $(shell git for-each-ref refs/tags --sort=-taggerdate --format='%(refname)' --count=1 | cut -d '/' -f 3)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

cn: clean
	go build -i -v -ldflags="-X main.version=$(VERSION) -X main.tag=$(TAG) -X main.branch=$(BRANCH)"

clean:
	rm cn || true
