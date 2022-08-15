BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml
SKEL_BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build-skel.yml
EXAMPLE_BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build-example.yml
BUILD_NUMBER ?= 0

.PHONY: all clean
.PHONY: test vet vendor generate precommit
.PHONY: license license-check
.PHONY: skel publish-skel
.PHONY: dist debug docker publish
.PHONY: generate-book

# Library
all: clean vet license-check test

test:
	$(BUILDER) download-test-deps
	$(BUILDER) execute-unit-tests

vet:
	$(BUILDER) go-vet

vendor:
	go mod vendor

generate:
	$(BUILDER) download-generate-deps
	$(BUILDER) generate

precommit: generate license
	$(BUILDER) go-fmt

license:
	$(BUILDER) license

license-check:
	$(BUILDER) license --check

skel:
	$(SKEL_BUILDER) build-tool

publish-skel:
	$(SKEL_BUILDER) publish-tool

install-skel:
	go install cto-github.cisco.com/NFV-BU/go-msx/cmd/skel

generate-book:
	$(BUILDER) copy-book-chapters
	mdbook build

# Example Microservice
dist:
	$(EXAMPLE_BUILDER) generate-build-info
	$(EXAMPLE_BUILDER) install-executable-configs
	$(EXAMPLE_BUILDER) install-dependency-configs
	$(EXAMPLE_BUILDER) install-swagger-ui
	$(EXAMPLE_BUILDER) build-executable

debug:
	$(EXAMPLE_BUILDER) build-debug-executable

docker: vendor
	$(EXAMPLE_BUILDER) docker-build

publish:
	$(EXAMPLE_BUILDER) docker-push

clean:
	rm -Rf dist
	rm -Rf vendor
