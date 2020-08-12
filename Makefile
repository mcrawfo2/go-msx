BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml

.PHONY: test dist docker debug publish generate clean precommit
.PHONY: skel install-skel

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
	mkdir -p dist/skel
	go build -o dist/skel/skel cmd/skel/skel.go

install-skel:
	go install cto-github.cisco.com/NFV-BU/go-msx/cmd/skel
