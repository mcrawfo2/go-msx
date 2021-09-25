BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml $(BUILD_FLAGS)

.PHONY: all test dist debug assemblies docker docker-publish manifest publish clean generate precommit tag
.PHONY: update-go-msx update-go-msx-build

all: clean test docker assemblies deployment manifest

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

deployment:
	${BUILDER} generate-deployment-variables

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

generate:
	$(BUILDER) download-generate-deps
	$(BUILDER) generate

precommit: generate
	$(BUILDER) go-fmt

tag:
	$(BUILDER) git-tag

update-go-msx:
	go get cto-github.cisco.com/NFV-BU/go-msx
	go mod tidy

update-go-msx-build:
	go get cto-github.cisco.com/NFV-BU/go-msx-build
	go mod tidy
#if GENERATOR_BEAT

update-go-msx-beats:
	go get cto-github.cisco.com/NFV-BU/go-msx-beats
	go mod tidy
#endif GENERATOR_BEAT

deploy-jenkins-job:
	$(BUILDER) deploy-jenkins-job

deploy-github-repo:
	$(BUILDER) deploy-github-repo

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
	cd ui && npm install && npm run build

#endif UI
