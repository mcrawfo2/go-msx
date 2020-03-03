# go-msx

go-msx is a Go library for microservices and tools interacting with MSX. 

## Requirements

- All Versions - Ensure your GOPATH is correctly set and referenced in your PATH.  For example:
  ```bash
  GOPATH=$HOME/go
  PATH=$PATH:$GOPATH/bin
  ```   

- Go 1.12 -  No extra configuration required.
- Go 1.13 -  Be sure to set your Go proxy settings correctly.  For example:
  ```bash
  GOPRIVATE="cto-github.cisco.com/NFV-BU"
  GOPROXY="https://proxy.golang.org,direct"
  ```
  
## Quick Start

- To add go-msx to an existing module-enabled go project:

    ```bash
    go get -u cto-github.cisco.com/NFV-BU/go-msx
    ```

- To create a new go-msx microservice skeleton project:
    - Install the `skel` generator:
        ```bash
        go install cto-github.cisco.com/NFV-BU/go-msx/cmd/skel
        ```
    - Run the `skel` generator:
        ```bash
        skel
        ```

## Documentation

### Cross-Cutting Concerns
* [Logging](log/README.md)
* [Configuration](config/README.md)
* [Lifecycle](app/README.md)
* [Dependencies](app/context.md)
* [Stats](stats/README.md)
* [Tracing](trace/README.md)

### Application Components
* Web Service
    * [Controller](#)
    * [Filter](#)

* Persistence
    * [Repository](#)
    * [Migration](#)

* Communication
    * [Integration](#)
    * [Streaming](#)



## License

Copyright (C) 2020 Cisco Systems, Inc.  All rights reserved.
