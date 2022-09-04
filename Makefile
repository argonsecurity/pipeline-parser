VERSION := $(shell git describe --tags --always)
LDFLAGS=-ldflags "-s -w -X=main.version=$(VERSION)"

# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

.PHONY: build
build:
	go build $(LDFLAGS) ./cmd/pipeline-parser

.PHONY: run
run:
	go run $(LDFLAGS) ./cmd/pipeline-parser $(RUN_ARGS)

.PHONY: tag
tag:
	@LEVEL=$(LEVEL) ./scripts/tag.sh

.PHONY: test
test:
	go clean -testcache
	go test ./...

.PHONY: test-coverage
test-coverage:
	go clean -testcache
	go test -coverprofile=coverage.out -covermode=atomic -v ./...