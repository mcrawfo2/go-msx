# Build Usage

`build` may be run directly using command-line targets.

* To list the targets and options for the build command, add the `-h` flag:
    ```bash
    go run cmd/build/build.go --config cmd/build/build.yml -h
    ```

* To get help for a particular target:
    ```bash
    go run cmd/build/build.go --config cmd/build/build.yml <target> -h
    ```

* To pass a custom build configuration, use the `--config` option:
    ```bash
    go run cmd/build/build.go --config cmd/build/build-custom.yml <target>
    ```

In addition to the numerous build targets (below), there are the following utility targets:

- `help`: display the help text
- `version`: display the current, and most recent `skel` build versions
- `completion`: generate the BASH completion script for `skel`


## Targets

```
Available Commands:
build-assemblies              Builds Assemblies
build-debug-executable        Build the binary debug executable
build-executable              Build the binary executable
build-installer-manifest      Generate the installer manifests
build-package                 Build the service deployment package
build-tool                    Build the binary tool
compare-openapi-spec          Compares the current openapi spec with the stored version
completion                    Generate the autocompletion script for the specified shell
deploy-github-repo            Deploy Github repository
deploy-jenkins-job            Deploy Jenkins job
deploy-package                Deploy the service to an MSX instance
docker-build                  Build the target release docker image
docker-build-debug            Build the target debug docker image
docker-push                   Push the target docker image to the upstream repository
docker-save                   Save the target docker image to the specified file
download-generate-deps        Download generate dependencies
download-seccomp-dependencies Download seccomp dependencies
download-test-deps            Download test dependencies
execute-unit-tests            Execute unit tests
generate                      Generate code
generate-build-info           Create a build metadata file
generate-deployment-variables Stage variables file with build version
generate-openapi-spec         Stores the current openapi spec into a file
generate-seccomp-profile      Create a seccomp profile
git-tag                       Tag the current commit
go-fmt                        Format all go source files
go-vet                        Vet all go source files
help                          Help about any command
install-asyncapi-ui           Installs AsyncAPI/Studio package
install-dependency-configs    Download dependency config files to distribution config directory
install-entrypoint            Copy custom entrypoint to distribution root directory
install-executable-configs    Copy configured files to distribution config directory
install-extra-configs         Copy custom files to distribution config directory
install-resources             Installs Resources
install-swagger-ui            Installs Swagger-UI package
license                       License all go source files
publish-binaries              Publishes Binaries
publish-installer-manifest    Deploy the installer manifests
publish-package               Publish the service deployment package
publish-tool                  Publish the binary tool
```