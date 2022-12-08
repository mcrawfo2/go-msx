# Summary

[Introduction](README.md)
[Contributing](CONTRIBUTING.md)

---

# Cross-Cutting Concerns
- [Logging](log/README.md)
- [Errors](types/docs/errors.md)
- [Configuration](config/README.md)
  - [Consul Configuration Provider](config/consulprovider/README.md) 
- [Lifecycle](app/README.md)
- [Dependencies](app/context.md)
- [Stats](stats/README.md)
- [Tracing](trace/README.md)

---

# Application Components

## Web Service 
- [Controller](webservice/controller.md)
- [Filter]()

## Persistence
- [Repository](sqldb/repository.md)
- [Migration]()

## Communication
- [Integration]()
  - [OpenAPI Client 🎉](integration/docs/openapi.md) 
- [Streaming]()
  - [Stream Operations 🎉](ops/streamops/README.md)
    - [Ports](ops/streamops/ports.md)
    - [Validation](ops/streamops/validation.md)
    - [Publishers](ops/streamops/publishers.md)
    - [Subscribers](ops/streamops/subscribers.md)
    - [AsyncApi](schema/asyncapi/README.md)
  - [Stream Providers]()
    - [Kafka]()
    - [SQL]()
    - [GoChannel]()
    - [Redis]()

---

# Utilities

- [Audit Events]()
- [Auditable Models]()
- [Cache](cache/lru/README.md)
- [Certificates and TLS](certificate/README.md)
- [Executing Commands]()
- [Health Checks]()
- [Http Client]()
- [Leader Election]()
- [Pagination]()
- [Resources](resource/README.md)
- [Retry](retry/README.md)
- [Sanitization](sanitize/README.md)
- [Scheduled Tasks](scheduled/README.md)
- [Transit Encryption](transit/README.md)
- [Validation]()

---

# Code Generation

- [Introduction](skel/README.md)
- [Installation](skel/docs/installation.md)
- [Usage](skel/docs/usage.md)
- [Projects]()
  - [Generic Microservice](skel/docs/projects-generic.md)
  - [Probes (Beats)](skel/docs/projects-beats.md)
  - [Service Pack Microservice]()
  - [Service Pack UI]()
- [Continuous Integration]()
- [Web Services]()
  - [Domains]()
  - [OpenAPI]()
- [Stream Services]()
  - [Channels](skel/asyncapi/channels.md)
  - [AsyncAPI](skel/asyncapi/spec.md)

---

# Builds 🎉

- [Introduction]()
- [Makefile Usage](build/docs/usage-make.md)
- [Build Usage](build/docs/usage-build.md)
- [Configuration](build/docs/config.md)
- [Build Targets](build/docs/targets.md)
  - [Project Maintenance](build/docs/targets-project.md)
  - [Development](build/docs/targets-development.md)
  - [Artifacts](build/docs/targets-artifacts.md)
  - [Verification](build/docs/targets-verification.md)
  - [Publishing](build/docs/targets-publishing.md)

---

# Continuous Integration

- [Checks 🎉](checks/README.md)
- [Jenkins]()