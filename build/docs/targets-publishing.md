# Publishing Build Targets

### `deploy-package`

The `deploy-package` target will upload your already-built package tar-ball to a specified installer container.
The installer container ssh "host" must be properly configured in your `~/.ssh/config` file, for example:

```ssh-config
Host installer-tme-dmz-01
    HostName rtp-dmz-bbhost.lab.ciscomsx.com
    User root
    Port 23556
    IdentityFile ~/.ssh/installer-tme-dmz.key
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
```

In this example, the installer container is named "installer-tme-dmz-01".  This name should be passed to
using the `deploy.host` configuration, for example:

```bash
go run build/cmd/build.go --config build/cmd/build.yml --deploy.host installer-tme-dmz-01 deploy-package
```

or

```bash
DEPLOY_HOST="installer-tme-dmz-01" make package-deploy
```

### `docker-push`

The `docker-push` target will publish the local docker image generated using `docker-build` to the docker
repository specified in the current build configuration.

The repository can be specified using the `docker.repository` configuration setting.

### `docker-save`

The `docker-save` target will output the local docker image generated using `docker-build` to a tar file
named `${info.app.name}.tar` in the current directory.  The tarred image will include the original repository and
image tag.

### `git-tag`

The `git-tag` target re-creates and overwrites any local and remote tags for the current version `${release}-${build}`.

This is commonly used after publish to tag the source repo with the build.

### `publish-binaries`

The `publish-binaries` target will deploy any assemblies and other installer binaries to artifactory.

The remote repository folder is specified through `artifactory.repository`.  Within the repository folder,
artifacts will be placed underneath `${msx.deploymentGroup}/${release}-${build}/` to isolate files from
each build and deployment group.

Binaries are specified in the `artifactory` configuration.  Assembly publishing can be disabled
setting the `artifactory.assemblies` to `false`.

### `publish-installer-manifest`

The `publish-installer-manifest` target executes Maven to deploy the manifest for the current build configuration.

### `publish-package`

The `publish-package` target will use your local S3 client to upload the service package to S3.  The correct S3 folder
will automatically be calculated.  Your S3 client (`aws s3 ...`) should be properly configured with credentials either
using environment variables or configuration files.

### `publish-tool`

The `publish-tool` target will upload the tool distribution packages (built with `build-tool`) to Artifactory.
Versioned and Latest will be published for easy URL distribution.
