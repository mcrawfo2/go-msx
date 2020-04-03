# go-msx

go-msx is a Go library for microservices and tools interacting with MSX. 

## Versioning

Currently this library and tools are in a pre-alpha state.  They are subject to
backwards-incompatible changes at *any* time.  After reaching the first stable release (v1.0.0),
[SemVer](https://semver.org/) will be used per industry and golang best practices.     

## Requirements

- Go 1.13+

    - Ensure your GOPATH is correctly set and referenced in your PATH.  For example:
      ```bash
      GOPATH=$HOME/go
      PATH=$PATH:$GOPATH/bin
      ``` 
        
      Execute the following commands:
        ```bash
        export GOPATH=~/go
        export PATH=$PATH:$(go env GOPATH)/bin
        ```

    - Be sure to set your Go proxy settings correctly.  For example:
      ```bash
      GOPRIVATE="cto-github.cisco.com/NFV-BU"
      GOPROXY="https://proxy.golang.org,direct"
      ```
      
      An other way to set the above is to run the command below.
      ```
      go env -w GOPRIVATE=cto-github.cisco.com/NFV-BU
      ```

- Git SSH configuration for `cto-github.cisco.com`

    - Ensure you have a registered SSH key referenced in your `~/.ssh/config`:
    
        ```
        Host github
              HostName cto-github.cisco.com
              User git
              IdentityFile ~/.ssh/github.key
        ```
      
      Note that this key must be registered via the [Github UI](https://cto-github.cisco.com/settings/keys).

    - Ensure you have SSH protocol override for git HTTPS urls to our github in your `~/.gitconfig`:
    
      ```
      [url "ssh://git@cto-github.cisco.com/"]
              insteadOf = https://cto-github.cisco.com/
      ```

## Quick Start

- To add go-msx to an existing module-enabled go project:

    ```bash
    go get -u cto-github.cisco.com/NFV-BU/go-msx
    ```

- To create a new go-msx microservice skeleton project:
    - Install the `skel` generator by running the command below from the `go-msx` directory:
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
