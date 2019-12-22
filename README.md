# go-msx

go-msx is a Go library for microservices and tools interacting with MSX. 

## Requirements

- Go 1.12

## Usage

- To add go-msx to an existing module-enabled go project:

    ```bash
    go get -u cto-github.cisco.com/NFV-BU/go-msx
    ```

- To create a new go-msx microservice skeleton project:
    - Install the `skel` generator:
        ```bash
        go install cto-github.cisco.com/NFV-BU/cmd/skel
        ```
    - Run the `skel` generator:
        ```bash
        skel
        ```

## License

Copyright (C) 2019 Cisco Systems, Inc.  All rights reserved.