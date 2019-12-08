BUILDER = go run $(BUILDER_FLAGS) cmd/build/build.go --config cmd/build/build.yml

.PHONY: test dist docker publish clean

test:


dist:
	$(BUILDER) generate-build-info
	$(BUILDER) install-executable-configs
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
