# Installation

You may install Skel either by cloning the git repo and using golang's install command, or by copying the 'skel' binary
from the repo's 'bin' directory; the latter is recommended since you will need to update it from time to time.

In either case, you will need to ensure that Git is set up and can communicate with the cto-github.cisco.com server. See
the [go-msx README](../../index.md#requirements) for details.

## Install from Artifactory

1. Download the skel tarball:

   **MacOS**:
   ```bash
   curl -L -o - https://engci-maven-master.cisco.com/artifactory/symphony-group/com/cisco/vms/go-msx-skel/latest/go-msx-skel-darwin-latest.tar.gz | tar -xzf -
   ```
   
   **Linux**:
   ```bash
   curl -L -o - https://engci-maven-master.cisco.com/artifactory/symphony-group/com/cisco/vms/go-msx-skel/latest/go-msx-skel-linux-latest.tar.gz | tar -xzf -
   ```

2. Move the skel binary to a location in your path:

   ```bash
   mv skel ~/go/bin
   ```

## Install via Go

Prerequisite: **Go 1.18+**

1. Ensure your GOBIN environment variable is correctly set and referenced in your PATH. For example:

    ```bash
    export GOPATH=~/go
    export PATH=$PATH:$GOPATH/bin
    ```

   Recall that GOBIN defaults to `$GOPATH/bin`, or `$HOME/go/bin` if the `GOPATH`
   environment variable is not set.

2. Be sure to set your Go proxy settings correctly. For example:

    ```bash
    go env -w GOPRIVATE=cto-github.cisco.com/NFV-BU
    ```

3. Check-out go-msx into your local workspace:

    ```bash
    mkdir -p $HOME/msx && cd $HOME/msx
    git clone git@cto-github.cisco.com:NFV-BU/go-msx.git
    cd go-msx
    go mod download
    ```

4. Install `skel`:

    ```bash
    make install-skel
    ```
