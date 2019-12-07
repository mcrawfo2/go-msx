BUILDER = go run -mod=vendor cmd/build/build.go --config cmd/build/build.yml

.PHONY: test dist docker publish clean

test:


dist:
	go mod vendor
	$(BUILDER) create-dist-dir
	$(BUILDER) generate-build-info
	$(BUILDER) install-executable-configs
	$(BUILDER) build-executable

docker:
	go mod vendor
	$(BUILDER) docker-build

publish:
	go mod vendor
	$(BUILDER) docker-push

clean:
	rm -Rf dist
	rm -Rf vendor
