BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml
SKEL_BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build-skel.yml
EXAMPLE_BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build-example.yml
BUILD_NUMBER ?= 0

.PHONY: test dist docker debug publish generate clean precommit
.PHONY: skel publish-skel vet

# Library

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

precommit: generate
	$(BUILDER) go-fmt

skel:
	$(SKEL_BUILDER) build-tool
	cp cmd/skel/README.md dist/tools/go-msx-skel/linux
	cp cmd/skel/README.md dist/tools/go-msx-skel/darwin

publish-skel:
	BUILD_NUMBER=$(BUILD_NUMBER) $(SKEL_BUILDER) publish-tool

install-skel:
	go install cto-github.cisco.com/NFV-BU/go-msx/cmd/skel

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
