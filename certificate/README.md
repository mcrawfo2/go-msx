# MSX Certificate

The certificate module provides access to static and dynamic TLS/x509
certificate sources, including the following providers:
- File - load certs and keys from the filesystem
- Vault - generate and renew certs and keys from Vault
- Cache - save upstream certs and keys to disk

The certificate module also provides a common configuration parser for
TLS configuration.

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
    ca-cert-file: /etc/pki/tls/certs/ca-identity.crt
    cert-file: /etc/pki/tls/certs/spokeservice.pem
    key-file: /etc/pki/tls/private/spokeservice-key.pem
```

When a subsystem requests certificates from the `identity` source, it will:
- Load certificates from the filesystem for each TLS connection

While _active_ renewal is not supported, the file provider does read in changes
to the file for each connection.  The cert/key files may be rotated as convenient.

#### Configuration Properties

| Key            | Default | Required | Description                           |
|----------------|---------|----------|---------------------------------------|
| `ca-cert-file` | -       | Required | CA Certificate File, PEM format       |
| `cert-file`    | -       | Required | Identity Certificate File, PEM format |
| `key-file`     | -       | Required | Identity Private Key File, PEM format | 

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
      - "${network.hostname}"
      - "${spring.application.name}.svc.kubernetes.cluster.local"
      - "${spring.application.name}.service.consul"
    ip-sans:
      - "${network.outbound.address}"
```

When a subsystem requests certificates from the `identity` source, it will:
- Generate an identity certificate and private key
- Renew the identity certificate half-way through its lifetime.

#### Configuration Properties

| Key         | Default | Required | Description                                               |
|-------------|---------|----------|-----------------------------------------------------------|
| `path`      | -       | Required | Vault PKI mount point                                     |
| `role`      | -       | Required | Vault PKI issuer role name                                |
| `cn`        | -       | Required | Desired identity certificate CN                           | 
| `alt-names` | -       | Optional | Desired identity certificate subject alternative names    | 
| `ip-sans`   | -       | Optional | Desired identity certificate IP subject alternative names | 

Note: `alt-names` and `ip-sans` will be stripped of empty entries so may include
undefined variables with empty defaults:

    - ${some.undefined.variable:}

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

#### Configuration Properties

| Key            | Default | Required | Description                           |
|----------------|---------|----------|---------------------------------------|
| `ca-cert-file` | -       | Required | CA Certificate File, PEM format       |
| `cert-file`    | -       | Required | Identity Certificate File, PEM format |
| `key-file`     | -       | Required | Identity Private Key File, PEM format | 

## TLS Configuration

TLS connection configuration is used in many places in go-msx, including:
- Kafka client
- Web server

For ease of use, these configurations have been unified into a single format.

#### Configuration Properties

| Key                    | Default      | Required | Description                                                         |
|------------------------|--------------|----------|---------------------------------------------------------------------|
| `enabled`              | `false`      | Optional | TLS should be enabled for this client/server                        |
| `insecure-skip-verify` | `false`      | Optional | Verify the server's certificate chain and host name                 |
| `min-version`          | `tls12`      | Optional | Minimum TLS version to support.  One of: tls10, tls11, tls12, tls13 |
| `certificate-source`   | -            | Optional | Server or Client certificate source.  Required for server.          |
| `cipher-suites`        | <sup>1</sup> | Optional | Cipher suites to enable.                                            |
| `server-name`          | -            | Optional | Server name to check in certificate when connecting from client.    |

<sup>1</sup> Current default cipher suites are:
- TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305
- TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
- TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
- TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
- TLS_RSA_WITH_AES_256_GCM_SHA384
- TLS_RSA_WITH_AES_256_CBC_SHA

