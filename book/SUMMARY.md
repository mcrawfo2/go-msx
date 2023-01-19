# Summary

- [Introduction](README.md)

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

# Application Patterns
- [Operations](types/docs/operations.md)
  - [Middleware](types/docs/middleware.md)
  - [Decorators](types/docs/decorators.md)
  - [Filters](types/docs/filters.md)
  - [Traced Operations](trace/docs/operations.md)
  - [Handlers]()
- [Ports](ops/docs/ports.md)
  - [Input Ports](ops/docs/input-ports.md) 
  - [Output Ports](ops/docs/output-ports.md) 
  - [Validation](ops/docs/validation.md)
- [Services](types/docs/services.md)

## Web Service
- [💀 REST Controller Mk I ](webservice/controller.md) 
- [🎉 REST Controller Mk II ](ops/restops/docs/controller.md) 
  - [REST Input Ports](ops/restops/docs/input-ports.md)
  - [REST Output Ports](ops/restops/docs/output-ports.md)
  - [Middleware](ops/restops/docs/middleware.md)

## Persistence
- [💀 CRUD Repository ](sqldb/repository.md) 
- [🎉Typed Repository](sqldb/docs/new_repositories_typed.md)
- [🎉Goqu Repository](sqldb/docs/new_repositories_goqu.md)
- [🎉SQL Repository](sqldb/docs/new_repositories_sql.md)
- [Migration]()

## Communication
- [Integration]()
  - [🎉 OpenAPI Client ](integration/docs/openapi.md) 
- [Streaming]()
  - [🎉 Stream Operations ](ops/streamops/README.md) 
    - [Ports](ops/streamops/ports.md) 
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

# Code Generation (skel)

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
- [Deployment]()
  - [Kubernetes]()
    - [Skaffold](skel/docs/skaffold.md)
---

# Builds

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

---

# For Developers of `go-msx`

- [Introduction]()
- [Contributing](CONTRIBUTING.md)
- [`skel`: Code Generation]()
  - [Execution Sequence](skel/docs/sequence.md)
  - [Making templates](skel/docs/templating.md)
