# Project Maintenance Build Targets

### `deploy-github-repo`

The `deploy-github-repo` target will create and/or reconfigure a GitHub repository with the appropriate users
and webhooks.

Target server and repository is configured via the `github.*` build [settings](config.md#github).

### `deploy-jenkins-job`

The `deploy-jenkins-job` target will upload your jenkins job as defined in the build/ci folder from `config.xml`.
This file will normally be auto-generated.

Target server and job is configured via the `jenkins.*` build [settings](config.md#jenkins).

### `update-go-msx`

### `update-go-msx-build`

### `update-go-msx-populator`

