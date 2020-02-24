# Skel

go-msx skeleton project generator

## Quick Start

Prerequisites:

- Go 1.13

0. Configure your GO proxy:
   ```bash
   export GOPRIVATE="cto-github.cisco.com/NFV-BU"
   export GOPROXY="https://proxy.golang.org,direct"
   ```

1. Install the skeleton generator:

    ```bash
    go install cto-github.cisco.com/NFV-BU/go-msx/cmd/skel
    ```
    
2. Invoke the skeleton generator and provide details of your new project:

    ```bash
    skel
    ? Project Parent Directory: /Users/mcrawfo2/msx
    ? Version: 3.9.0
    ? App name: administrationservice
    ? App display name: Administration Microservice
    ? App description: Administration Microservice
    ? Web server port: 9210
    ? Context path: /administration
    ```
   
   Your skeleton project will be created in the target directory.
   
3. Open the new project in GoLand or IDE of your choice.
