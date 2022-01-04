.PHONY: build clean lint image imagepush

VERSION    ?= $(shell git rev-parse --short HEAD)
DOCKER_TAG ?= mccutchen/datadog-agent-stub:$(VERSION)

BUILD_ARGS ?= -ldflags="-s -w"

# Built binaries will be placed here
DIST_PATH ?= dist

# Tool dependencies
TOOL_BIN_DIR     ?= $(shell go env GOPATH)/bin
TOOL_GOLINT      := $(TOOL_BIN_DIR)/golint
TOOL_ERRCHECK    := $(TOOL_BIN_DIR)/errcheck
TOOL_STATICCHECK := $(TOOL_BIN_DIR)/staticcheck


# =============================================================================
# build
# =============================================================================
build:
	mkdir -p $(DIST_PATH)
	go build $(BUILD_ARGS) -o $(DIST_PATH)/datadog-agent-stub .

clean:
	rm -rf $(DIST_PATH) $(COVERAGE_PATH)


# =============================================================================
# lint
# =============================================================================
lint: deps
	test -z "$$(gofmt -d -s -e .)" || (echo "Error: gofmt failed"; gofmt -d -s -e . ; exit 1)
	go vet ./...
	$(TOOL_GOLINT) -set_exit_status ./...
	$(TOOL_ERRCHECK) ./...
	$(TOOL_STATICCHECK) ./...


# =============================================================================
# docker images
# =============================================================================
image:
	docker build -t $(DOCKER_TAG) .

imagepush: image
	docker push $(DOCKER_TAG)

# =============================================================================
# dependencies
# =============================================================================
deps: $(TOOL_GOLINT) $(TOOL_ERRCHECK) $(TOOL_STATICCHECK)

$(TOOL_GOLINT):
	go install golang.org/x/lint/golint@latest

$(TOOL_ERRCHECK):
	go install github.com/kisielk/errcheck@latest

$(TOOL_STATICCHECK):
	go install honnef.co/go/tools/cmd/staticcheck@latest
