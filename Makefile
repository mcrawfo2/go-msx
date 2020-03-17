BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml

.PHONY: test dist docker publish clean precommit

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

precommit:
	$(BUILDER) generate
	$(BUILDER) go-fmt
