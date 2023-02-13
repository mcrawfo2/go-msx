.PHONY: install-txtartools
.PHONY: completion-test
.PHONY: generate-app-test
.PHONY: run-root-test
.PHNOY: most

# Library
all: install-txtartools completion-test generate-app-test

install-txtartools:
	go install cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtargen
	go install cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtarwrap

completion-test:
	@echo "Target: completion test script"
	@./fixtures/build/make-completions-test.sh

generate-app-test:
	@echo "Target: make generate-app test script"
	@./fixtures/build/make-generate-app-test.sh


run-root-test:
	@$FIXT/final/root/root-test.sh

run-generate-app-test:
	echo "Target: run generate-app test"
	testscript -e FIXT=$FIXT $FIXT/final/generate-app-test.txtar

most:
	fixtures/build/make-any.sh generate-app
	fixtures/build/make-any.sh generate-build
	fixtures/build/make-any.sh generate-certificate
	fixtures/build/make-any.sh generate-channel
	fixtures/build/make-any.sh generate-channel-asyncapi
	fixtures/build/make-any.sh generate-channel-publisher
	fixtures/build/make-any.sh generate-channel-subscriber
	fixtures/build/make-any.sh generate-deployment-variables
	fixtures/build/make-any.sh generate-dockerfile
	fixtures/build/make-any.sh generate-domain
	fixtures/build/make-any.sh generate-domain-beats
	fixtures/build/make-any.sh generate-domain-openapi
	fixtures/build/make-any.sh generate-domain-system
	fixtures/build/make-any.sh generate-domain-tenant
	fixtures/build/make-any.sh generate-git
	fixtures/build/make-any.sh generate-github
	fixtures/build/make-any.sh generate-goland
	fixtures/build/make-any.sh generate-harness
	fixtures/build/make-any.sh generate-jenkins
	fixtures/build/make-any.sh generate-kubernetes
	fixtures/build/make-any.sh generate-local
	fixtures/build/make-any.sh generate-manifest
	fixtures/build/make-any.sh generate-migrate
	fixtures/build/make-any.sh generate-service-pack
	fixtures/build/make-any.sh generate-skel-json
	fixtures/build/make-any.sh generate-spui
	fixtures/build/make-any.sh generate-test
	fixtures/build/make-any.sh generate-timer
	fixtures/build/make-any.sh generate-topic-publisher
	fixtures/build/make-any.sh generate-topic-subscriber
	fixtures/build/make-any.sh generate-vscode
	fixtures/build/make-any.sh generate-webservices