.PHONY: install-txtartools
.PHONY: run-root-test
.PHONY: make-root-goldens
.PHONY: make-tests
.PHONY: run-tests

# Library
all: install-txtartools make-root-goldens run-root-test

install-txtartools:
	go install cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtargen
	go install cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtarwrap
	go install cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtarunwrap
	go install cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtarcmp

make-root-goldens:
	@echo "Target: make root goldens"
	${FIXT}/build/make-roots.sh

run-root-test:
	@echo "Target: run root test"
	${FIXT}/final/root/root-test.sh

make-tests:
	@echo "Target: make tests"
	go run ${FIXT}/build/make-tests/make-tests.go

run-tests:
	@echo "Target: run tests"
	go run ${FIXT}/final/run-tests/run-tests.go

