# Artifacts Build Targets

### `build-assemblies`

The `build-assemblies` target collects folders into tarballs and places their output into the staging
assembly folder (`dist/assembly`).

For each entry in `assembly.custom.*`, the target will create a tar file named `${prefix}-${release}-${build}.tar`.

### `build-debug-executable`

The `build-debug-executable` target compiles the `main` module of the current application and outputs it
to the staging executable folder (`dist/root/usr/bin`).  Unlike `build-executable`, however it outputs a
binary suitable for debugging.  This can be included in a container to remotely debug the application.

### `build-executable`

The `build-executable` target compiles the `main` module of the current application and outputs it to
the staging executable folder (`dist/root/usr/bin`).

Flags passed to `go build` can be customized using the `go.env.*.GOFLAGS` configuration.

### `build-installer-manifest`

The `build-installer-manifest` target will create install manifests for integration servers and production
installation.  The contents of the manifests will be dynamically generated from the `artifactory` and `assemblies`
configuration, along with the docker image generated from the current build configuration.

To deploy the manifest artifact use the `publish-installer-manifest` target.

### `build-package`

The `build-package` target will generate a Service package to be uploaded and deployed by the service pack deployer.
It will include the standard service contents, including assemblies, manifests, images, deployment variables, and
other binaries.

To deploy the package artifact use the `publish-package` target.

### `build-tool`

The `build-tool` target will compile and generate a Tool binary distribution (.tar.gz) to be uploaded and deployed
to Artifactory (or elsewhere).  It will include the binary and any resources defined in the `tool` configuration
section.

NOTE: These binaries are statically compiled and therefore must not be distributed.

### `docker-build`

The `docker-build` target will create a docker image for the current build configuration.  The contents of the
image are stages using `make dist` inside a build container, and then deployed onto an MSX base image to create
the runtime container image.

The base image can be specified using the `docker.repository` and `docker.base-image` configuration [settings](config.md#docker).

The docker image will be named in the format `${docker.repository}/${info.app.name}:${release}-${build}`.

### `docker-build-debug`

The `docker-build-debug` target will create a debugging docker image for the current build configuration.  
The contents of the image are stages using `make dist` inside a build container, and then deployed onto an 
MSX base image to create the runtime container image.

The base image can be specified using the `docker.repository` and `docker.base-image` configuration [settings](config.md#docker).

The docker image will be named in the format `${docker.repository}/${info.app.name}:${release}-${build}`.

### `download-generate-deps`

The `download-generate-deps` target installs cross-project generation dependencies, including:
- github.com/vektra/mockery
- bou.ke/staticfiles

### `download-seccomp-dependencies`

The `download-seccomp-dependencies` target installs the seccomp-profiler for generating seccomp profiles.
See `generate-seccomp-profile`, below.

### `generate-build-info`

The `generate-build-info` target creates the build-specific metadata file `buildinfo.yml`, including version information,
and build timestamps.

The metadata file is generated directly into the staging configuration folder (`dist/root/etc/${app.name}`).  
This file will be parsed on MSX Application startup during the configuration phase, and used to register the
service metadata with Consul.

Default values for the `info.build` fields should be specified in the application `bootstrap.yml` file to enable
local development before generating the metadata file.

### `generate-deployment-variables`

The `generate-deployment-variables` target creates a YAML ansible variables file compatible with the MSX installer.
This file will be published during `publish`.

### `generate-seccomp-profile`

The `generate-seccomp-profile` target creates the configuration file `seccomp.yml`, listing the expected set of linux
syscalls to be allowed during execution.  This prevents a compromised executable from making unauthorized syscalls.

### `install-asyncapi-ui`

The `install-asyncapi-ui` target downloads the AsyncApi Studio package and extracts the
relevant files to the staging web folder (`dist/root/var/lib/${app.name}/www`)

### `install-dependency-configs`

The `install-dependency-configs` target scans maven artifacts for `default-*.properties` files and copies them
into the staging configuration folder (`dist/root/etc/${app.name}`).  At runtime, a go-msx microservice will read these
files ensuring MSX microservices across frameworks have the same default configuration.

### `install-executable-configs`

The `install-executable-configs` target copies configuration files from the `main` module of the application being
built to the staging configuration folder (`dist/root/etc/${app.name}`).

The list of configuration files to be copied is specified in the build configuration at `executable.config-files`:

```yaml
executable:
  configFiles:
    - bootstrap.yml
    - dnaservice.production.yml
```

### `install-extra-configs`

### `install-resources`

The `install-resources` target copies static files from the project tree to the staging resources folder
(`dist/var/lib/${app.name}`).

The list of resources to be copied is specified in the build configuration at `resources.*`.

### `install-swagger-ui`

The `install-swagger-ui` target downloads the Swagger UI webjar and MSX Swagger artifacts and extracts the
relevant files to the staging web folder (`dist/root/var/lib/${app.name}/www`)

