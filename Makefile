BINARY  := spec-cli
MODULE  := github.com/9Ashwin/spec-cli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
DATE    := $(shell date +%Y-%m-%d)
LDFLAGS := -s -w -X $(MODULE)/internal/build.Version=$(VERSION) -X $(MODULE)/internal/build.Date=$(DATE)
PREFIX  ?= /usr/local

.PHONY: all build test vet fmt fmt-check install uninstall clean release

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

# Release builds for npm distribution.
RELEASE_DIR := dist
PLATFORMS   := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

release: vet fmt-check
	@rm -rf $(RELEASE_DIR)
	@mkdir -p $(RELEASE_DIR)
	@for p in $(PLATFORMS); do \
		os=$$(echo $$p | cut -d/ -f1); \
		arch=$$(echo $$p | cut -d/ -f2); \
		out=$(BINARY)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then out="$$out.exe"; fi; \
		echo "  GOOS=$$os GOARCH=$$arch go build -o $(RELEASE_DIR)/$$out"; \
		GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags "$(LDFLAGS)" -o $(RELEASE_DIR)/$$out .; \
		if [ "$$os" != "windows" ]; then \
			tar -czf $(RELEASE_DIR)/$(BINARY)-$(VERSION)-$$os-$$arch.tar.gz -C $(RELEASE_DIR) $$out; \
		else \
			zip -j $(RELEASE_DIR)/$(BINARY)-$(VERSION)-$$os-$$arch.zip $(RELEASE_DIR)/$$out; \
		fi; \
	done
	@cd $(RELEASE_DIR) && shasum -a 256 *.tar.gz *.zip > checksums.txt
	@echo ""
	@echo "Release artifacts in $(RELEASE_DIR)/:"
	@ls -la $(RELEASE_DIR)/*.tar.gz $(RELEASE_DIR)/*.zip $(RELEASE_DIR)/checksums.txt 2>/dev/null
	@echo ""
	@echo "Next: cp $(RELEASE_DIR)/checksums.txt . && npm publish"
