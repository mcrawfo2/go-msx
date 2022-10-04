# Build Configuration

go-msx Build uses YAML build configuration files to define the build to be executed.
A build configuration describes the build metadata for one of:
- Microservice
- Command Line tool
- Service Pack
- Library

A build configuration file can include these artifacts:
- Binary artifacts
- Assemblies (tarballs)
- Resources
- Runtime Configuration Files
- Docker Images

## Configuration Sources

Like all go-msx applications, go-msx-build can retrieve configuration from a variety of sources:
- Environment
- Command-Line Options
- Build Configuration Files
- Application Configuration Files
- Defaults

To specify the primary build configuration file, pass the `--config` option to build:

```bash
go run cmd/build/build.go --config cmd/build/build.yml
```

This will normally be handled by the Makefile.

Configuration passed in by either Environment Variables or Command-Line Options will override
values also specified in Files or Defaults.

### Environment Variables

Some settings are intended to be injected from environment variables.  These include:
- `docker.username` (`DOCKER_USERNAME`)
- `docker.password` (`DOCKER_PASSWORD`)
- `artifactory.username` (`ARTIFACTORY_USERNAME`)
- `artifactory.password` (`ARTIFACTORY_PASSWORD`)
- `build.number` (`BUILD_NUMBER`)
- `manifest.folder` (`MANIFEST_FOLDER`)
- `jenkins.username` (`JENKINS_USERNAME`)
- `jenkins.password` (`JENKINS_PASSWORD`)
- `github.token` (`GITHUB_TOKEN`)

It is considered unsafe or inflexible to store them directly in the configuration file.
The default generated Jenkinsfile will automatically inject these environment
variables as required by the relevant steps.

### Application Configuration

Some settings below are intended to be read from the application configuration files.  
These include:
- `info.app.*` - `bootstrap.yml`
- `server.*` - `bootstrap.yml`

To ensure these are being read from the correct source, ensure the `executable.config-files` list
contains the base application configuration files (eg `bootstrap.yml`).

Example:

```yaml
executable:
  configFiles:
    - bootstrap.yml
```

## Configuration Sections

### `executable`

The `executable` configuration specifies the entrypoint and primary configuration file(s) of this build.

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `executable.cmd` | `app` | Optional | The `cmd` sub-folder containing the application `main` module. |
| `executable.config-files` | - | Required | A list of configuration files within the main module to include in the build. |

Example:

```yaml
executable:
  configFiles:
    - bootstrap.yml
    - dnaservice.production.yml
```

### `msx`

The `msx` configuration specifies details of the MSX release to interface with.

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `msx.release` | - | Required | The MSX release of this microservice (output). |
| `msx.deployment-group` | - | Required | The deployment group of the build. |
| `msx.platform.parent-artifacts` | - | Required | Maven artifact roots to scan for default properties. |
| `msx.platform.version` | - | Required | The platform version to use for locating maven artifacts.  Accepts `EDGE` and `STABLE` builds. |
| `msx.platform.include-groups` | - | Required | Maven artifact groupIds to include in artifact scanning. |
| `msx.platform.swagger-artifact` | `com.cisco.nfv:nfv-swagger` | Optional | MSX artifact groupId and artifactId for nfv-swagger. |
| `msx.platform.swagger-webjar` | `org.webjars:swagger-ui:3.23.11` | Optional | Maven artifact triple for swagger web jar. |

Example:

```yaml
msx:
  release: 3.10.0
  deploymentGroup: dna
  platform:
    parentArtifacts:
      - com.cisco.vms:vms-service-starter
      - com.cisco.vms:vms-service-starter-core
      - com.cisco.vms:vms-service-starter-kafka
      - com.cisco.nfv:nfv-integration-consul-leader
    version: 3.10.0-EDGE
    includeGroups: "com.cisco.**"
```

### `docker`

The `docker` configuration controls interactions with the docker daemon, global repository,
images, and `Dockerfile` scripts.

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `docker.dockerfile` | `docker/Dockerfile` | Optional | The `Dockerfile` used for this build. |
| `docker.baseimage` | `msx-base-buster:3.9.0-70` | Optional | The base image within the repository. |
| `docker.repository` | `dockerhub.cisco.com/vms-platform-dev-docker` | Optional | The repository source and destination. |
| `docker.username` | - | Optional | User name to authenticate to repository. |
| `docker.password` | - | Optional | Password to authenticate to repository. |
| `docker.buildkit` | - | Optional | `true` to use docker buildkit when building the docker image. |
| `docker.base.dynamic.enabled` | `true` | Optional | `true` to use manifests to dynamically locate the base docker image. |
| `docker.base.dynamic.stream` | `EI-Stable` | Optional | Manifest stream to search within for manifests |
| `docker.base.dynamic.version` | `${msx.release}` | Optional | MSX release to search within for manifests |
| `docker.base.dynamic.manifest` | `msxbase-bullseye-manifest` | Optional | MSX manifest to search within for builds |
| `docker.base.dynamic.image` | `msx-base-image` | Optional | Manifest key identifying the image to use |

Example:

```yaml
docker:
    dockerfile: build/package/Dockerfile
```

### `kubernetes`

The `kubernetes` configuration provides defaults for generating kubernetes manifests.

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `kubernetes.group` | `platformms` | Optional | The kubernetes group used for pods in production. |

### `manifest`

The `manifest` configuration specifies how to build and publish installer manifests.

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `manifest.folder` | `Build-Stable` | Optional | The maven output folder to publish the manifest to. |

### `resources`

The `resources` section identifies the files to be included as part of the docker image.

Each entry has the following properties:

| Key        | Default   | Required | Description |
|------------|-----------|----------|-------------|
| `resources.includes` | - | Optional | List of globs of files to include. Processed first. |
| `resources.excludes` | - | Optional | List of globs of files to exclude. Processed second. |

Example:

```yaml
resources:
  includes:
    - "/internal/migrate/**/*.sql"
    - "/internal/populate/**/*"
  excludes:
    - "/internal/populate/**/*.go"
```

### `assemblies`

The `assemblies` configuration specifies `.tar` file generation.  The `.tar` files will be included in generated
manifests and published (unless disabled).

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `assemblies.root` | `platform-common` | Optional | The folder from which assemblies are created by default. All sub-folders with a 'templates' folder or 'manifest.json' are included. |
| `assemblies.custom` | - | Optional | List of custom assemblies to include.  See [below](#assembliescustom) |

Example:

```yaml
assemblies:
  root: platform-common
```

#### `assemblies.custom`

The `assemblies.custom` setting contains a list of custom assemblies to generate.  These
will be uploaded to artifactory and recorded as binaries in the manifest, unless disabled
with `artifactory.assemblies`.

Each entry in this list has the following properties:

| Key                 | Default | Required | Description |
|---------------------|---------|----------|-------------|
| `path` | - | Required | The root path of the assembly files. |
| `path-prefix` | - | Optional | A folder inside the assembly to prefix the files during the build. |
| `manifest-prefix` | - | Required | The prefix of the file name in the manifest. |
| `manifest-key` | - | Required | The location of the entry in the JSON manifest. |
| `includes` | `/**/*` | Optional | Glob of files to include.  Processed first. |
| `excludes` | - | Optional | Glob of files to exclude.  Processed second. |

Example:

To create an assembly file called "skyfallui-files-${release}-${build}.tar":

```yaml
assemblies:
  custom:
    - path: ui/build
      pathPrefix: services
      manifestPrefix: skyfallui-files
      manifestKey: ${msx.deploymentGroup}-ui
```

- Each file from the `ui/build` subtree will be prefixed with the `services` folder in the output tar.
  e.g. 'ui/build/dna/index.js' will be relocated to `services/dna/index.js`.
- The assembly will be added to the generated artifact manifests at e.g. `dna-ui`.

### `artifactory`

The `artifactory` configuration specifies artifactory connectivity, folders, binaries, and images.

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `artifactory.assemblies` | `true` | Optional | Include [assemblies](#assemblies) in publishing and manifests |
| `artifactory.installer` | `deployments/kubernetes` | Optional | The folder in which installer binaries can be found.  eg pod, rc, meta templates. |
| `artifactory.repository` | `https://.../vms-3.0-binaries` | Optional | The base url for storing published artifacts |
| `artifactory.installer-folder` | `binaries/vms-3.0-binaries` | Optional | The folder prefix of binaries to record in the manifest |
| `artifactory.username` | - | Optional | The user name with which to authenticate to Artifactory. |
| `artifactory.password` | - | Optional | The password with which to authenticate to Artifactory. |
| `artifactory.custom` | - | Optional | List of custom binaries to include.  See [below](#artifactorycustom) |
| `artifactory.images` | - | Optional | List of docker images to include. |

Example:

```yaml
artifactory:
  installer: deployments/production
  images:
    - nacservice
```

#### `artifactory.custom`

The `artifactory.custom` setting contains a list of custom binaries to include.  These
will be uploaded to artifactory and recorded in the manifest.

Each entry in this list has the following properties:

| Key                             | Required | Description |
|---------------------------------|----------|-------------|
| `path` | Required | The source path of the file to include. |
| `output-name` | Required | The destination name of the file. |
| `manifest-prefix` | Required | The prefix of the file name in the manifest. |
| `manifest-key` | Required | The location of the entry in the JSON manifest. |

Example:

```yaml
artifactory:
  custom:
    - path: deploymentvariables/nac_deployment_variables.yml
      outputName: nac_deployment_variables.yml
      manifestPrefix: deployment-variables
      manifestKey: deployment_variables
    - path: deploymentvariables/nac_variables.yml
      outputName: nac_variables.yml
      manifestPrefix: variables
      manifestKey: variables    
```

### `go`

The `go` configuration specifies environment variables and options to be passed to Go tools during the build.

| Key | Description |
|-----|-------------|
| `go.env.all.*` | Environment variables for all platforms |
| `go.env.linux.*` | Environment variables for linux platform |
| `go.env.darwin.*` | Environment variables for darwin (MacOS) platform |
| `go.vet.options[*]` | List of command line options to pass to `go vet` |

### `build`

The `build` configuration specifies information about the build used to generate `buildinfo.yml`.

| Key | Default | Required | Description |
|-----|---------|----------|-------------|
| `build.number` | `SNAPSHOT` | Required | The build number of this build. |
| `build.group` | `com.cisco.msx` | Optional | The build group. 

### `info.app`

The `info.app` configuration specifies details about the application used across all parts of the build.

| Key | Default | Required | Description |
|-----|---------|----------|-------------|
| `info.app.name` | - | Required | The name of the application being built. |
| `info.app.attributes.display-name` | - | Required | The display name of the application being built. |

Example:

```yaml
info.app:
    name: dnaservice
    attributes:
      displayName: DNA Microservice
```

### `server`

The `server` configuration specifies details about the web server used across all parts of the build.

| Key | Default | Required | Description |
|-----|---------|----------|-------------|
| `server.port` | - | Required | The web server port of the application being built. |

Example:

```yaml
server:
    port: 9393
```

### `jenkins`

The `jenkins` configuration specifies details about the Jenkins CI server used by the project.

| Key | Default | Required | Description |
|-----|---------|----------|-------------|
| `jenkins.job` | - | Optional | The simplified job path to the Jenkins Job on the server. |
| `jenkins.server` | `https://jenkins.infra.ciscomsx.com` | Optional | The base url of the Jenkins CI server. |
| `jenkins.username`| - | Optional | User name to authenticate to Jenkins. |
| `jenkins.password` | - | Optional | API Token to authenticate to Jenkins. Can be created on the User Configure page in Jenkins UI. |

Example:

```yaml
jenkins.job: eng-sp-umbrella/builds/umbrellaservice
```

### `github`

The `github` configuration specifies details about the GitHub Source Control server used by the project.

| Key | Default | Required | Description |
|-----|---------|----------|-------------|
| `github.repository` | '${spring.application.name}' | Optional | The name of the repository on the server. |
| `github.organization` | 'NFV-BU' | Optional | The owner of the repository on the server. |
| `github.server` | `https://cto-github.cisco.com` | Optional | The base url of the GitHub server. |
| `github.token` | - | Optional | API Token to authenticate to GitHub. Can be created on the `User Settings > Developer Settings > Personal Access Tokens` page in the GitHub UI. |
| `github.hook.push` | `${jenkins.server}/github-webhook/` | Optional | Github Push Webhook to configure on the repository. |
| `github.hook.pull-request` | `${jenkins.server}/ghprbhook/` | Optional | Github PR Webhook to configure on the repository. |
| `github.teams.jenkins` | `Jenkins-generic-users` | Optional | GitHub CI Team to assign write access to the repository. |
| `github.teams.eng` | - | Optional | GitHub Engineering Team to assign write access to the repository. |

Example:

```yaml
github.organization: xiaoydu
```


### `aws`

The `aws` configuration specifies credentials and target details for AWS.

| Key | Default | Required | Description |
|-----|---------|----------|-------------|
| `aws.access-key-id`| `${aws.access.key.id}` | Optional | Access Key Id to authenticate to AWS. |
| `aws.secret-access-key` | `${aws.secret.access.key}` | Optional | Secret Access Key to authenticate to AWS. |

These values default to the standard environment variables (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`) so
no extra configuration should be required if using them.


### `deploy`

The `deploy` configuration specifies the target for package deployment.

| Key | Default | Required | Description |
|-----|---------|----------|-------------|
| `deploy.host`| - | Required | SSH config host name to target for deployment.  Must point to an installer container. |
