# MSX Certificate

The certificate module provides access to static and dynamic TLS/x509
certificate sources, including the following providers:
- File - load certs and keys from the filesystem
- Vault - generate and renew certs and keys from Vault
- Cache - save upstream certs and keys to disk

## Sources

A source identifies the provider and provider parameters required to obtain identity and authority certificates.

Each source is defined in the configuration using a unique name (lowercase alphanumeric only).
Source properties are configured under `certificate.source.<sourcename>`, for example:

```yaml
certificate.source:
  identity:
      provider: ...
      property1: ...
      property2: ...
      property3: ...
```

## Providers

Each source specifies a provider and its parameters.
Individual providers are detailed in the following sections.

### File

To specify the local filesystem as the source for certificates, use the File provider:

```yaml
certificate.source:
  identity:
    provider: file
    ca-cert-file: /etc/ssl/certs/ca-identity.crt
    cert-file: /certs/spokeservice.pem
    key-file: /certs/spokeservice-key.pem
```

When a subsystem requests certificates from the `identity` source, it will:
- Load certificates from the filesystem for each TLS connection

While _active_ renewal is not supported, the file provider does read in changes
to the file for each connection.  The cert/key files may be rotated as convenient.

### Vault

To specify Vault as the source for certificates, use the Vault provider:

```yaml
certificate.source:
  identity:
    provider: vault
    path: pki/vms
    role: "${spring.application.name}"
    cn: "${spring.application.name}"
    alt-names:
      - "${remote.service.hostname}"
      - "${spring.application.name}.svc.kubernetes.cluster.local"
      - "${spring.application.name}.service.consul"
    ip-sans:
      - "${kubernetes.pod.ip}"
      - "${remote.service.ip}"
```

When a subsystem requests certificates from the `identity` source, it will:
- Generate an identity certificate and private key
- Renew the identity certificate half-way through its lifetime.

### Cache

To configure a cache for a remote certificate source, use the Cache provider:

```yaml
certificate.source:

  identity:
    provider: vault
    path: pki/vms
    role: "${spring.application.name}"
    cn: "${spring.application.name}"
    alt-names:
      - "${remote.service.hostname}"
    ip-sans:
      - "${kubernetes.pod.ip}"
      - "${remote.service.ip}"
      
  identitycache:
    provider: cache
    upstream-source: identity
    key-file: "/certs/${spring.application.name}-key.pem"
    cert-file: "/certs/${spring.application.name}.pem"
    ca-cert-file: "/etc/ssl/certs/ca-identity.crt"
```

When a subsystem requests certificates from the `identitycache` source, it will:
- Generate and store an identity certificate and private key under `/certs`
- Retrieve and store the authority certificate under `/etc/ssl/certs`.
- Renew the identity certificate half-way through its lifetime.
