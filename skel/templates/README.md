# ${app.name}

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
    - Publish docker image
* `clean`
    - Remove test and distribution directories
* `precommit`
    - Prepare code for check-in

To execute a target from the shell, run `make <target>`.

## Goland

This project can be opened in goland, and has several default run configurations:

* `${app.name} (local)`
    - Run the microservice against a local installation of MSX
* `${app.name} (remote)`
    - Run the microservice against a remote installation of MSX.  The remote
      connection configuration is locate in the [local config](./local) subdirectory
      of the project.
* `make clean` 
    - Clean the project
* `make dist`
    - Build the project distribution folder
* `make docker`
    - Build the docker image 
