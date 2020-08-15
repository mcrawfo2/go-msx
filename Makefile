BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml
SKEL_BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build-skel.yml
BUILD_NUMBER ?= 0

.PHONY: test dist docker debug publish generate clean precommit
.PHONY: skel publish-skel

test:
	$(BUILDER) download-test-deps
	$(BUILDER) execute-unit-tests

dist:
	$(BUILDER) generate-build-info
	$(BUILDER) install-executable-configs
	$(BUILDER) install-dependency-configs
	$(BUILDER) install-swagger-ui
	$(BUILDER) build-executable

debug:
	$(BUILDER) build-debug-executable

docker:
	go mod vendor
	$(BUILDER) docker-build

publish:
	$(BUILDER) docker-push

clean:
	rm -Rf dist
	rm -Rf vendor

generate:
	$(BUILDER) download-generate-deps
	$(BUILDER) generate

precommit:
	$(BUILDER) generate
	$(BUILDER) go-fmt

skel:
	$(SKEL_BUILDER) build-tool
	cp cmd/skel/README.md dist/tools/go-msx-skel/linux
	cp cmd/skel/README.md dist/tools/go-msx-skel/darwin

publish-skel:
	BUILD_NUMBER=$(BUILD_NUMBER) $(SKEL_BUILDER) publish-tool

install-skel:
	go install cto-github.cisco.com/NFV-BU/go-msx/cmd/skel
