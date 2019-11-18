BUILDER = go run build.go

.PHONY: test dist docker clean publish

test:


dist:
	$(BUILDER) create-dist-dir
	$(BUILDER) generate-build-info
	$(BUILDER) install-executable-configs
	$(BUILDER) generate-dockerfile

docker: dist
	go mod vendor
	$(BUILDER) docker-build
	$(BUILDER) docker-login

publish: dist
	$(BUILDER) docker-push

clean:
	$(BUILDER) delete-dist-dir
	rm -Rf vendor
