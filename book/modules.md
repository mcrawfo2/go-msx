# Modules Overview

go-msx is composed of a number of layers and modules.

![Modules Diagram](go-msx-modules.svg)

## Application

- `app`: Application lifecycle
- `background`: Application background errors

## Platform

- `restops`: HTTP REST endpoints
- `streamops`: Stream channels
- `scheduled`: Scheduled tasks
- `audit`: Updating model audit fields and logging auditable events
- `exec`: Subprocess execution
- `httpclient`: HTTP client
- `rbac`: Role-based Access control
- `security`: Attribute-based Access control
- `retry`: Reliability
- `sanitize`: Input/Output sanitization
- `transit`: Transit encryption
- `validate`: Data validation
- `migrate`: Database migration
  - `/sqldb/migrate`: SQL database migration
- `populate`: API population

# Integration
- `discovery`: Register and Locate microservices
  - `consulprovider`: Consul discovery provider 
- `stream`: Communicate using streams
- `webservice`: REST web server
  - `adminprovider`: Admin actuator
  - `aliveprovider`: Liveness actuator
  - `apilistprovider`: API list documentation
  - `asyncapiprovider`: AsyncApi documentation
  - `authprovider`: Authentication
  - `debugprovider`: Debug profiling
  - `envprovider`: Configuration actuator
  - `healthprovider`: Health actuator
  - `idempotency`: Idempotency-Key filter
  - `infoprovider`: Info actuator
  - `loggersprovider`: Logging actuator
  - `maintenanceprovider`: Maintenance actuator
  - `metricsprovider`: Metrics actuator
  - `prometheusprovider`: Prometheus stats
  - `swaggerprovider`: Swagger documentation
- `cli`: Command line interaction
- `health`: Health checks
  - `consulcheck`: Consul health check
  - `kafkacheck`: Kafka health check
  - `redischeck`: Redis health check
  - `sqldbcheck`: SQL health check
  - `vaultcheck`: Vault health check
- `integration`: REST API client
- `cache`: Caching
  - `lru`: In-Memory cache provider
  - `/redis/cache`: Redis cache provider
- `operations`: Operations support
- `schema`: Schema documentation
  - `asyncapi`: AsyncApi schema documentation
  - `js`: JSON schema documentation
  - `openapi`: OpenApi schema documentation
  - `swagger`: Swagger schema documentation
- `leader`: Leader election
  - `consulprovider`: Consul leader provider
- `certificate`: Certificate management

# Infrastructure
- `consul`: Consul driver
- `vault`: Vault driver
- `redis`: Redis driver
- `sqldb`: SQL database driver
- `kafka`: Kafka driver
- `trace/datadog`: Datadog tracing
- `trace/jaeger`: Jaeger tracing

# Core
- `config`: Configuration
- `log`: Logging
- `trace`: Tracing
- `stats`: Statistics
- `fs`: Filesystems
- `resources`: Resources
- `types`: Reusable data types

# Continuous Integration
- `build`: Build execution
