# go-msx

go-msx is a Go library for microservices and tools interacting with MSX. 

## Versioning

Currently this library and tools are in a pre-alpha state.  They are subject to
backwards-incompatible changes at *any* time.  After reaching the first stable release (v1.0.0),
[SemVer](https://semver.org/) will be used per industry and golang best practices.     

## Requirements

- Go 1.18+

    - Ensure your GOPATH is correctly set and referenced in your PATH.  For example:
        ```bash
        export GOPATH=~/go
        export PATH=$PATH:$GOPATH/bin
        ```

    - Be sure to set your Go proxy settings correctly.  For example:
        ```bash
        go env -w GOPRIVATE=cto-github.cisco.com/NFV-BU
        ```

- Git SSH configuration for `cto-github.cisco.com`

    - Ensure you have a registered SSH key referenced in your `~/.ssh/config`:
    
        ```
        Host cto-github.cisco.com
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

- Skel tool for code generation

    - Check out go-msx into your local workspace:
        
        ```bash
        mkdir -p $HOME/msx && cd $HOME/msx
        git clone git@cto-github.cisco.com:NFV-BU/go-msx.git
        cd go-msx
        go mod download
        ```

    - Install `skel`:
  
        ```bash
        make install-skel
        ```

## Quick Start

- To continue working on an existing go-msx project:

    - Return to the original project README instructions
      and continue.

- To add go-msx to an existing module-enabled go project:

    ```bash
    go get -u cto-github.cisco.com/NFV-BU/go-msx
    ```

- To create a new go-msx microservice skeleton project:
    
    ```bash
    cd $HOME/msx
    skel
    ```
  
## Documentation

Please visit our [internal site](https://cto-github.cisco.com/pages/NFV-BU/go-msx)
or [public site](https://mcrawfo2.github.io/go-msx/).

## License

Copyright © 2019-2022, Cisco Systems Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

