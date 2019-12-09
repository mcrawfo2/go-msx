# Skel

go-msx skeleton project generator

## Quick Start

Prerequisites:

- Go 1.12

1. Install the skeleton generator:

    ```bash
    go install cto-github.cisco.com/NFV-BU/go-msx/cmd/skel
    ```
    
2. Invoke the skeleton generator and provide details of your new project:

    ```bash
    skel
    ? Project Parent Directory: /Users/mcrawfo2/vms-3.1/demos
    ? App name: fancyservice
    ? App display name: Fancy Microservice
    ? App description: Fancy Microservice
    ? Web server port: 9998
    ```
   
   Your skeleton project will be created in the target directory.
   
3. Open the new project in GoLand or IDE of your choice.
