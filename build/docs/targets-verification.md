# Verification Build Targets

### `compare-openapi-consumer-spec`

The `compare-openapi-producer-spec` target will obtain the latest version of a
_consumer_ OpenApi contract specification, as identified in the build configuration.
After obtaining the latest contract, it will compare it with the stored version,
and generate a report of the differences.

Consumer local and remote specification are configured via the `openapi.contracts[*]` build [setting](config.md#openapi).

### `compare-openapi-producer-spec`

The `compare-openapi-producer-spec` target will obtain the latest version of the
microservice _producer_ OpenApi contract specification and compare it with the stored
version. After comparison, it will generate a report of the differences.

Producer specification file is configured via the `openapi.spec` build [setting](config.md#openapi).

### `compare-openapi-specs`

The `compare-openapi-specs` target will execute the `compare-openapi-producer-spec`
target, and then `compare-openapi-consumer-spec` target for each registered contract.
A summary report will be generated.

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

### `go-vet`

The `go-vet` target executes `go vet` on directories which contain `*.go` files (excluding the `vendor` directory).
Options to pass to `go vet` can be specified in the build configuration under `go.vet.options`.
