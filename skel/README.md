# Skel: MSX service skeleton generator

  - [Introduction](#introduction)
  - [Common behaviour when generating](#common-behaviour-when-generating)
  - [Generation Targets](#generation-targets)
    - [MSX Projects](#msx-projects)
      - [Generic MSX App](#generic-msx-app)
      - [Beats](#beats)
      - [Service Pack](#service-pack)
      - [Service Pack UI](#service-pack-ui)
    - [Web Services](#web-services)
      - [Domains](#domains)
      - [OpenAPI](#openapi)
    - [Stream Services](#stream-services)
      - [Channels](#channels)
      - [AsyncAPI](#asyncapi)
      - [more {...}](#more-)
  - [Installation](#installation)
    - [Install via Golang](#install-via-golang)
  - [Running Skel](#running-skel)
  
## Introduction

Skel is a tool for generating MSX service skeletons. It is a part of the [go-msx](https://cto-github.cisco.com/NFV-BU/go-msx) library and tools, and the skeleton projects it generates are compatible with the MSX framework.

## Common behaviour when generating

{explain what approach skel takes in general, which applies to all generate commands, subdirs, overwriting, git ops etc.}

## Generation Targets

Skel can generate:

1. MSX projects: {which are ...}
2. Web services: {which ...}
3. Stream services: {which ...}

### MSX Projects

#### Generic MSX App

- _Contains:_ {wazinnit}  
- _Root dir:_ ./{serviceName}/  
- _Command:_ `generate-app`  
- _Menu:_ Generate Archetype | Generic Microservice  

A generic MSX app skeleton that contains various bony bits ...  

#### Beats

- Contains:
- Root dir:
...

#### Service Pack

#### Service Pack UI

### Web Services

#### Domains

#### OpenAPI

### Stream Services

#### Channels

#### AsyncAPI

#### more {...}

## Installation

You may install Skel either by cloning the git repo and using golang's install command, or by copying the 'skel' binary from the repo's 'bin' directory; however, the former is recommended since you will likely need Go installed and set up anyway.

In either case, you will need to ensure that Git is set up and can communicate with the cto-github.cisco.com server. See the [go-msx README](https://cto-github.cisco.com/NFV-BU/go-msx/blob/master/README.md) for details.

### Install via Golang

1. Go 1.18+ is required

2. Ensure your GOBIN environment variable is correctly set and referenced in your PATH. For example:

    ```bash
    export GOPATH=~/go
    export PATH=$PATH:$GOPATH/bin
    ```

    Recall that GOBIN defaults to `$GOPATH/bin`, or `$HOME/go/bin` if the `GOPATH`
    environment variable is not set.  

3. Be sure to set your Go proxy settings correctly. For example:

    ```bash
    go env -w GOPRIVATE=cto-github.cisco.com/NFV-BU
    ```

4. Check-out go-msx into your local workspace:

    ```bash
    mkdir -p $HOME/msx && cd $HOME/msx
    git clone git@cto-github.cisco.com:NFV-BU/go-msx.git
    cd go-msx
    go mod download
    ```

5. Install `skel`:

    ```bash
    make install-skel
    ```

## Running Skel

Skel may be run using either command-line sub-commands or by using its minimal, but hopefully helpful, menus.

To start the menu version, execute simply `skel`

To get help on skel add the `-h` flag, thus: `skel -h`. For help on a particular command, thus: `skel <command> -h`

In addition to the numerous generating commands, there are the following utility commands:

1. `help` display the help text
2. `version` display the current, and most recent skel build verions {I found this output puzzling so we may need to explain more}
3. `add-go-msx-dependency` ...
4. `completion` ...
