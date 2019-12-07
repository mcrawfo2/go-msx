BUILDER = go run cmd/build/build.go --config cmd/build/build.yml

.PHONY: test dist docker publish clean

test:


dist:
	$(BUILDER) create-dist-dir
	$(BUILDER) generate-build-info
	$(BUILDER) install-executable-configs
	$(BUILDER) build-executable

docker:
	go mod vendor
	$(BUILDER) docker-build

publish:
	$(BUILDER) docker-push

clean:
	$(BUILDER) delete-dist-dir
	rm -Rf vendor
