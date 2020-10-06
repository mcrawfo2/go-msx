BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml $(BUILD_FLAGS)

.PHONY: all test dist debug assemblies docker docker-publish manifest publish clean precommit tag
.PHONY: update-go-msx update-go-msx-build

all: clean test docker dist

test:
	$(BUILDER) download-test-deps
	$(BUILDER) execute-unit-tests

dist:
	$(BUILDER) generate-build-info
	$(BUILDER) install-executable-configs
	$(BUILDER) install-resources
	$(BUILDER) install-dependency-configs
	$(BUILDER) install-swagger-ui
	$(BUILDER) build-executable

debug:
	$(BUILDER) build-debug-executable

assemblies:
	$(BUILDER) build-assemblies

manifest:
	$(BUILDER) build-installer-manifest

docker:
	go mod vendor
	$(BUILDER) docker-build

docker-publish: docker
	$(BUILDER) docker-push

publish:
	$(BUILDER) docker-push
	$(BUILDER) publish-binaries
	$(BUILDER) publish-installer-manifest

clean:
	rm -Rf dist
	rm -Rf vendor

precommit:
	$(BUILDER) generate
	$(BUILDER) go-fmt

tag:
	$(BUILDER) git-tag

update-go-msx:
	go get cto-github.cisco.com/NFV-BU/go-msx
	go mod tidy

update-go-msx-build:
	go get cto-github.cisco.com/NFV-BU/go-msx-build
	go mod tidy

#if GENERATOR_SP
.PHONY: update-administrationservice update-catalogservice

update-administrationservice:
	go get cto-github.cisco.com/NFV-BU/administrationservice
	go mod tidy

update-catalogservice:
	go get cto-github.cisco.com/NFV-BU/catalogservice
	go mod tidy

#endif GENERATOR_SP
#if UI
.PHONY: ui

ui:
	cd ui && npm run build

#endif UI
