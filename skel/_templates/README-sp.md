# ${app.shortname}

## Quick Start 

### Prerequisites

1. [Go-MSX prerequisites](https://cto-github.cisco.com/NFV-BU/go-msx#requirements)
    
2. [Mini-VMS](https://cto-github.cisco.com/NFV-BU/mini-vms)

   1. All Cisco developers can obtain a pre-installed [Virtual Machine](https://jenkins.infra.ciscomsx.com/job/eng-generic/job/development/job/dev-ops-vm/).

   2. To enable routing from the UI, add the `${app.shortname}` route to Consul.  
      From your Mini-VMS folder, execute:
    
      ```bash
      ./vms infra consul route ${app.shortname} ${server.contextpath}
      docker-compose restart router
      ```

   3. To connect your local microservice to a Remote Mini-VMS instance, configure
      the location of your remote environment.

      Edit the `./local/${app.name}.remote.yml` and set your remote
      environment's IP address in `remote.service.address`. 
      
      If you are connected to VPN, set your local tunnel address and uncomment the 
      `discovery.ipaddress` line inside `spring.cloud.consul`.

      ```yaml
      remote.service.address: 10.81.85.174
      local.tunnel.address: 10.11.12.13

      spring.cloud:
        consul:
          host: ${remote.service.address}
          discovery.instanceId: local
          #discovery.ipaddress: ${local.tunnel.address}
        vault:
          host: ${remote.service.address}
      ```

### Goland

#### Generate IDE run configuration

From the project root directory, execute: 

```bash
skel generate-goland
```

To install `skel`, see [prerequisites](#prerequisites), above.

After generating, open the project folder in GoLand.
   
#### Local MSX Development Instance

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

3. Initialize the database:
 
    - From the Goland toolbar, execute the `${app.name} migrate (remote)` run configuration

4. Populate platform microservices:
 
    - From the Goland toolbar, execute the `${app.name} populate (remote)` run configuration
 
5. Start the main application entrypoint:

    - From the Goland toolbar, execute the `${app.name} (remote)` run configuration

### Visual Studio Code

#### Generate IDE run configuration

From the project root directory, execute:

```bash
skel generate-vscode
```

To install `skel`, see [prerequisites](#prerequisites), above.

After generating, open the project folder in Visual Studio Code.
    
#### Local MSX Development Instance

1. Stage artifacts into the `dist` folder:

    - From the VS Code Run Tasks, execute the `make dist` task

2. Initialize the database:

    - From the VS Code Run pane, execute the `${app.name} migrate (local)` run target

3. Populate platform microservices:

    - From the VS Code Run pane, execute the `${app.name} populate (local)` run target

4. Start the main application entrypoint:

    - From the VS Code Run pane, execute the `${app.name} (local)` run target

#### Remote MSX Development Instance

1. Configure the `local/${app.name}.remote.yaml` file per [prerequisites](#prerequisites), above.

2. Stage artifacts into the `dist` folder:

    - From the VS Code Run Tasks, execute the `make dist` task
    
3. Initialize the database:

    - From the VS Code Run pane, execute the `${app.name} migrate (local)` run target
    
4. Populate platform microservices:

    - From the VS Code Run pane, execute the `${app.name} populate (local)` run target

5. Start the main application entrypoint:

    - From the VS Code Run pane, execute the `${app.name} (local)` run target

### CLI

#### Local MSX Development Instance

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

#### Remote MSX Development Instance

1. Configure the `local/${app.name}.remote.yaml` file per [prerequisites](#prerequisites), above.

2. Stage artifacts into the `dist` directory:
    ```bash
    make dist
    ``` 

   All the required resources and artifacts will be staged.

3. Switch to the main entrypoint directory:
    ```bash
    cd cmd/app
    ```

   From this directory, you can execute the entrypoint and its various commands.

4. Migrate your database:
    ```bash
    go run main.go migrate --profile remote
    ```

5. Populate platform microservices:
    ```bash
    go run main.go populate --profile remote
    ```

6. Start the main application entrypoint:
    ```bash
    go run main.go --profile remote
    ```

## Source Code

The following directories are provided for your code and configuration:

* `build`
    - CI/CD job configuration
    - Container packaging

* `cmd`
    - Executable entry-points  
    - Application entrypoint is `cmd/app/main.go`
    - Build entrypoint is `cmd/build/build.go`
    
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
