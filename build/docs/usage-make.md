# Makefile Usage

go-msx uses GNU Make to present abstract build targets for developers and Continuous Integration
systems.  This allows for consistent builds across a variety of environments, and development
of Continuous Integration without a hosted job runner.

`make` may be run directly to execute targets.

* To list the targets in the Makefile, execute the `help` target:
    ```bash
    make help
    ```
  Sample output is shown below.


* To pass flags to the `go` command when executing `build`:
    ```bash
    export BUILDER_FLAGS='-exec xprog'
    make vet
    ```

* To pass flags to the `build` command when executing `build`:
    ```bash
    export BUILD_FLAGS='--artifactory.password="cisco123"'
    make publish
    ```

In addition to the numerous build targets (below), there are the following utility targets:

- `help`: display the help text

## Targets

```
assemblies               Generate supplemental artifacts
clean                    Remove any temporary build outputs
debug                    Build a debug executable
deploy-github-repo       Configure a standard github repository
deploy-jenkins-job       Upload a standard Jenkins build job to MSX Jenkins
deployment               Generate the installer deployment variables
deps                     Install dependencies
dist                     Build all outputs required for a container image
docker                   Generate a docker image for this service
docker-debug             Generate a debugging docker image for this service
docker-publish           Publish a docker image for this service
generate                 Execute code generation
help                     Show this help
manifest                 Generate the installer manifest
openapi-compare          Compare the openapi contracts for the microservice
openapi-generate         Store the openapi contract for the microservice
package                  Generate an SLM package
package-deploy           Deploy this service using SLM to an MSX instance
package-publish          Publish this service as an SLM package to S3
precommit                Ensure the code is ready for committing to version control
publish                  Publish all artifacts required for the installer
tag                      Tag the repository with a new PATCH version number
test                     Execute unit tests
update-go-msx            Update the go-msx library dependency to the latest version
update-go-msx-build      Update the go-msx-build library dependency to the latest version
update-go-msx-populator  Update the go-msx-populator library dependency to the latest version
verify-contract          Ensure the openapi contract matches the generated code
vet                      Use go vet to validate sources
```