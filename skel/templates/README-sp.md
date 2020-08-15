# ${app.name}

## Quick Start 

### Prerequisites

1. [Go-MSX prerequisites](https://cto-github.cisco.com/NFV-BU/go-msx#requirements)

2. [Mini-VMS](https://cto-github.cisco.com/NFV-BU/mini-vms)

    - A pre-installed virtual machine can be allocated from our Ops team 
      [here](https://jenkins.infra.ciscomsx.com/job/eng-generic/job/development/job/dev-ops-vm/)

### Goland

#### Local MSX Development  Instance

1. Stage artifacts into the `dist` folder:

    - From the Goland toolbar, execute the `make dist` run configuration
    
2. Initialize the database:
 
    - From the Goland toolbar, execute the `${app.name} migrate (local)` run configuration

3. Populate platform microservices:
 
    - From the Goland toolbar, execute the `${app.name} populate (local)` run configuration
 
4. Start the main application entrypoint:
 
    - From the Goland toolbar, execute the `${app.name} (local)` run configuration
    
#### Remote MSX Development Instance

1. Stage artifacts into the `dist` folder:

    - From the Goland toolbar, execute the `make dist` run configuration
    
2. Configure the location of your remote environment:
   
    - Edit the `./local/${app.name}.remote.yml` and set your remote
      environment's IP address:
      
        ```yaml
        remote.service.address: 10.81.85.174
        
        spring.cloud:
          consul:
            host: ${remote.service.address}
            discovery.instanceId: local
          vault:
            host: ${remote.service.address}
        ```
    
3. Initialize the database:
 
    - From the Goland toolbar, execute the `${app.name} migrate (remote)` run configuration

4. Populate platform microservices:
 
    - From the Goland toolbar, execute the `${app.name} populate (remote)` run configuration
 
5. Start the main application entrypoint:

    - From the Goland toolbar, execute the `${app.name} (remote)` run configuration

### CLI
1. Stage artifacts into the `dist` directory:
    ```bash
    make dist
    ``` 
   
   All the required resources and artifacts will be staged.
   
2. Switch to the main entrypoint directory:
    ```bash
    cd cmd/app
    ```
   
   From this directory, you can execute the entrypoint and its various commands. 
   
3. Migrate your database:
    ```bash
    go run main.go migrate
    ```

4. Populate platform microservices:
    ```bash
    go run main.go populate
    ```
   
5. Start the main application entrypoint:
    ```bash
    go run main.go
    ```

## Source Code

The following directories are provided for your code and configuration:

* `cmd`
    - Executable entrypoints.  
    - Application entrypoint is `cmd/app/main.go`.
    - Build entrypoint is `cmd/build/build.go`.
    
* `docker`, `k8s`
    - Containerization artifacts.
    
* `local`
    - Local configuration files.  Excluded from repository.

For further details on go project structure, see 
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

## Build

${app.name} uses GNU Make for the build system.  The following targets
are available:

* `test` 
    - Execute tests for all modules and generate reports ([`test`](./test))
* `dist` 
    - Build distribution directory and executable ([`dist`](./dist)) 
* `debug` 
    - Build debug executable
* `docker` 
    - Build docker image
* `publish` 
    - Publish docker image and installer manifest
* `clean`
    - Remove test and distribution directories
* `precommit`
    - Prepare code for check-in

To execute a target from the shell, run `make <target>`.

## Goland

This project can be opened in goland, and has several default run configurations:

### Local Mini-MSX

* `${app.name} (local)`
    - Run the microservice against a local installation of MSX
* `${app.name} migrate (local)`
    - Migrate the microservice
* `${app.name} populate (local)`
    - Populate platform microservices

### Remote Mini-MSX
* `${app.name} (remote)`
    - Run the microservice against a remote installation of MSX.  The remote
      connection configuration is locate in the [local config](./local) subdirectory
      of the project.
* `${app.name} migrate (remote)`
    - Migrate the microservice
* `${app.name} populate (remote)`
    - Populate platform microservices

### Build Targets
* `make clean` 
    - Clean the project
* `make dist`
    - Build the project distribution folder
* `make docker`
    - Build the docker image 
* `make publish`
    - Publish a build of the service pack
