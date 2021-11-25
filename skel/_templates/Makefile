BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml $(BUILD_FLAGS)

.PHONY: all vet test dist debug assemblies manifest package package-publish package-deploy
.PHONY: docker docker-publish publish clean generate precommit tag
.PHONY: update-go-msx update-go-msx-build update-go-msx-beats

#if GENERATOR_APP
all: clean vet test docker assemblies deployment manifest
#endif GENERATOR_APP
#if GENERATOR_BEAT
all: clean vet test docker deployment manifest
#endif GENERATOR_BEAT
#if GENERATOR_SP
all: clean vet test docker assemblies deployment manifest package
#endif GENERATOR_SP

vet:
	$(BUILDER) go-vet

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
	$(BUILDER) generate-deployment-variables
#if GENERATOR_BEAT
#else GENERATOR_BEAT

assemblies:
	$(BUILDER) build-assemblies
#endif GENERATOR_BEAT

manifest:
	$(BUILDER) build-installer-manifest
#if GENERATOR_SP

package:
	$(BUILDER) build-package

package-publish: all
	$(BUILDER) publish-package

package-deploy:
	$(BUILDER) deploy-package
#endif GENERATOR_SP

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
