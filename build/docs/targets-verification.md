# Verification Build Targets

### `compare-openapi-spec`



### `download-test-deps`

The `download-test-deps` target installs cross-project test dependencies, including:
- github.com/axw/gocov/gocov
- github.com/AlekSi/gocov-xml
- github.com/stretchr/testify/assert
- github.com/stretchr/testify/mock
- github.com/stretchr/testify/http
- github.com/pmezard/go-difflib/difflib
- github.com/jstemmer/go-junit-report

### `execute-unit-tests`

The `execute-unit-tests` target searches for testable directories (those containing `*_test.go` files),
and invokes their unit tests while collecting line coverage data.  It then generates coverage reports
from the coverage data.

| Format    | Output File                   |
|-----------|-------------------------------|
| HTML      | `test/gocover.html`           |
| JUnit     | `test/junit-report.xml`       |
| Cobertura | `test/cobertura-coverage.xml` |  

### `generate-openapi-spec`

The `generate-openapi-spec` target generates and stores an OpenApi specification to disk.

### `go-vet`

The `go-vet` target executes `go vet` on directories which contain `*.go` files (excluding the `vendor` directory).
Options to pass to `go vet` can be specified in the build configuration under `go.vet.options`.
