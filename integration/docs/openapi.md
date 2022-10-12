# OpenApi Client

MSX enables generating OpenApi clients with ease.   

## Client Generation

The following snippets show an example of how to generate an OpenApi client
for Manage Microservice V8 APIs:

- `cmd/build/build.yml`:
  ```yaml
  # Integration Generation
  generate:
    - path: internal/integration/manage
      openapi:
        spec: ../.openapi/manage-service-8.yaml
        config: ../.openapi/manage-service-8-config.json
  ```

- `internal/integration/.openapi/manage-service-8.yaml`:
  Place the openapi contract in this file.

- `internal/integration/.openapi/manage-service-8-config.json`:
  ```json
  {
    "generateInterfaces": true,
    "structPrefix": false,
    "packageName": "manage",
    "enablePostProcessFile": true
  }
  ```

- `internal/integration/manage/.openapi-generator-ignore`
  ```
  .gitignore
  go.mod
  go.sum
  .openapi-generator-ignore
  .travis.yml
  api/**
  docs/**
  git_push.sh
  ```
  
After the above pieces are in place, you can execute the generate build step:

```bash
make generate
```

## Contract Validation

To ensure the upstream contract remains compatible with your local version:


1. Add the following snippet to `build.yml`:

    ```yaml
    # Contract Management
    openapi:
      # Remote (consumer) API contract pairs
      contracts:
        - consumer: internal/integration/.openapi/manage-service-8.yaml
          producer: https://cto-github.cisco.com/raw/NFV-BU/msx-platform-specs/develop/manage-service-8.yaml
    
      # Sources for well-known schemas
      alias:
        - from: https://api.swaggerhub.com/domains/Cisco-Systems46/msx-common-domain/8
          to: https://cto-github.cisco.com/raw/NFV-BU/msx-platform-specs/sdk1.0.10/common-domain-8.yaml
    ```

    Any internal GitHub links will use the GitHub Personal Access Token from your environment
    (`GITHUB_TOKEN`) when retrieving the file.  Ensure you have an up-to-date PAT configured.

2. Add a check to your `build/ci/checks.yml`:

    ```yaml
    checks:
      - name: OpenApi
        commands:
          - make: openapi-compare
        analyzers:
          - builtin: generate-openapi-report
    ```

    This will ensure each commit to your repo checks for backwards-incompatible changes to 
    the contract.
