BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml $(BUILD_FLAGS)

.PHONY: all deps vet test dist debug dist-debug assemblies manifest package package-publish package-deploy
.PHONY: docker docker-debug docker-push docker-publish publish clean generate precommit tag
.PHONY: update-go-msx update-go-msx-build update-go-msx-beats
.PHONY: verify-contract openapi-generate openapi-compare

#if GENERATOR_APP
all: clean deps vet test docker assemblies deployment manifest
#endif GENERATOR_APP
#if GENERATOR_BEAT
all: clean deps vet test docker deployment manifest package
#endif GENERATOR_BEAT
#if GENERATOR_SP
all: clean deps vet test docker assemblies deployment manifest package
#endif GENERATOR_SP

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1\3/p' \
	| sed -n 's/^\([^ ]*\) \(.*\)$$/\1,\2/p' \
	| column -t  -s ',' \
	| sort

deps: ## Install dependencies
	go mod download

vet:  ## Use go vet to validate sources
	$(BUILDER) go-vet

test: ## Execute unit tests
	$(BUILDER) download-test-deps
	$(BUILDER) execute-unit-tests

dist: ## Build all outputs required for a container image
	$(BUILDER) generate-build-info
	$(BUILDER) install-executable-configs
	$(BUILDER) install-resources
#if EXTERNAL
#else EXTERNAL
	$(BUILDER) install-dependency-configs
	$(BUILDER) install-swagger-ui
	$(BUILDER) install-asyncapi-ui
#endif EXTERNAL
	$(BUILDER) build-executable
	$(BUILDER) download-seccomp-dependencies
	$(BUILDER) generate-seccomp-profile

debug: dist ## Build a debug executable
	$(BUILDER) build-debug-executable

dist-debug: ## Build all outputs required for a debug container image
	$(BUILDER) generate-build-info
	$(BUILDER) install-executable-configs
	$(BUILDER) install-resources
#if EXTERNAL
#else EXTERNAL
	$(BUILDER) install-dependency-configs
#endif EXTERNAL
	$(BUILDER) build-debug-executable

docker-push:
	$(BUILDER) docker-push

deployment: ## Generate the installer deployment variables
	$(BUILDER) generate-deployment-variables
#if GENERATOR_BEAT
#else GENERATOR_BEAT

assemblies: ## Generate supplemental artifacts
	$(BUILDER) build-assemblies
#endif GENERATOR_BEAT

manifest: ## Generate the installer manifest
	$(BUILDER) build-installer-manifest

package: ## Generate an SLM package
	$(BUILDER) build-package

package-publish: all ## Publish this service as an SLM package to S3
	$(BUILDER) publish-package

package-deploy: ## Deploy this service using SLM to an MSX instance
	$(BUILDER) deploy-package

docker: ## Generate a docker image for this service
	go mod vendor
	$(BUILDER) docker-build

docker-debug: ## Generate a debugging docker image for this service
	go mod vendor
	$(BUILDER) docker-build-debug

docker-publish: docker ## Publish a docker image for this service
	$(BUILDER) docker-push

publish: ## Publish all artifacts required for the installer
	$(BUILDER) docker-push
	$(BUILDER) publish-binaries
	$(BUILDER) publish-installer-manifest

clean: ## Remove any temporary build outputs
	rm -Rf dist
	rm -Rf vendor

generate: ## Execute code generation
	$(BUILDER) download-generate-deps
	$(BUILDER) generate

precommit: generate ## Ensure the code is ready for committing to version control
	$(BUILDER) go-fmt

tag: ## Tag the repository with a new PATCH version number
	$(BUILDER) git-tag

update-go-msx: ## Update the go-msx library dependency to the latest version
	go get cto-github.cisco.com/NFV-BU/go-msx
	go mod tidy

update-go-msx-build: ## Update the go-msx-build library dependency to the latest version
	go get cto-github.cisco.com/NFV-BU/go-msx-build
	go mod tidy

update-go-msx-populator: ## Update the go-msx-populator library dependency to the latest version
	go get cto-github.cisco.com/NFV-BU/go-msx-populator
	go mod tidy
#if GENERATOR_BEAT

update-go-msx-beats: ## Update the go-msx-beats library dependency to the latest version
	go get cto-github.cisco.com/NFV-BU/go-msx-beats
	go mod tidy
#endif GENERATOR_BEAT

deploy-jenkins-job: ## Upload a standard Jenkins build job to MSX Jenkins
	$(BUILDER) deploy-jenkins-job

deploy-github-repo: ## Configure a standard github repository
	$(BUILDER) deploy-github-repo

#if GENERATOR_SP
.PHONY: update-administrationservice update-catalogservice

update-administrationservice: ## Update the administration dependency to the latest version
	go get cto-github.cisco.com/NFV-BU/administrationservice
	go mod tidy

update-catalogservice: ## Update the catalog dependency to the latest version
	go get cto-github.cisco.com/NFV-BU/catalogservice
	go mod tidy

#endif GENERATOR_SP
#if UI
.PHONY: ui

ui:
	cd ui && npm install && npm run build

#endif UI

verify-contract: ## Ensure the openapi contract matches the generated code
	ls ccv.yml && ccv check || echo "WARNING! No ccv.yml found"

openapi-generate: ## Store the openapi contract for the microservice
	$(BUILDER) generate-openapi-spec

openapi-compare: ## Compare the openapi contracts for the microservice
	$(BUILDER) download-generate-deps
	$(BUILDER) compare-openapi-specs
