# Consul Configuration Provider

The Consul config provider reads settings from the KV version 1 consul plugin.
It currently supports two separate read paths: default and service-specific.
These read paths are expected to exist directly under the KV mount point.

The provider will, by default, load KV settings from the following locations:

- `userviceconfiguration/defaultapplication` - default settings
- `userviceconfiguration/${info.app.name}` - service-specific settings

## Provider Configuration

| Key                                        | Default                                                                                | Required | Description                                                 |
|--------------------------------------------|----------------------------------------------------------------------------------------|----------|-------------------------------------------------------------|
| spring.cloud.consul.config.enabled         | false                                                                                  | Optional | Enable loading configuration from consul KV                 |
| spring.cloud.consul.config.disconnected    | false                                                                                  | Optional | Activate "disconnected" mode for CLI commands               |
| spring.cloud.consul.config.prefix          | userviceconfiguration                                                                  | Optional | Consul KV mount point                                       |
| spring.cloud.consul.config.default-context | defaultapplication                                                                     | Optional | KV folder path under mount point containing global settings |
| spring.cloud.consul.config.pool            | false                                                                                  | Optional | Pool consul connections                                     |
| spring.cloud.consul.config.delay           | 3s                                                                                     | Optional | Retry delay after KV setting retrieval failure              |
| spring.cloud.consul.config.required        | `- ${spring.cloud.consul.config.prefix}/${spring.cloud.consul.config.default-context}` | Optional | KV settings paths that _must_ return KV values              |
