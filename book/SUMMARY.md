# Summary

[Introduction](README.md)

---

# Cross-Cutting Concerns
- [Logging](log/README.md)
- [Configuration](config/README.md)
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
- [Retry]()
- [Sanitization](sanitize/README.md)
- [Scheduled Tasks](scheduled/README.md)
- [Transit Encryption](transit/README.md)
- [Validation]()
