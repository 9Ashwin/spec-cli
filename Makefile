BINARY  := spec-cli
MODULE  := github.com/9Ashwin/spec-cli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
DATE    := $(shell date +%Y-%m-%d)
LDFLAGS := -s -w -X $(MODULE)/internal/build.Version=$(VERSION) -X $(MODULE)/internal/build.Date=$(DATE)
PREFIX  ?= /usr/local

.PHONY: all build test vet fmt fmt-check install uninstall clean

all: build vet

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test -race -count=1 ./...

vet:
	go vet ./...

fmt:
	gofmt -s -w .

# fmt-check fails when any file would be reformatted by gofmt.
fmt-check:
	@unformatted=$$(gofmt -l . | grep -v '^\.claude/' || true); \
	if [ -n "$$unformatted" ]; then \
		echo "Unformatted Go files:"; \
		echo "$$unformatted"; \
		echo "Run 'make fmt' and commit."; \
		exit 1; \
	fi

install: build
	install -d $(PREFIX)/bin
	install -m755 $(BINARY) $(PREFIX)/bin/$(BINARY)
	@echo "OK: $(PREFIX)/bin/$(BINARY) ($(VERSION))"

uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)

clean:
	rm -f $(BINARY)
