# ${app.name}

## Placeholder README

## Derived from a Generated Angular 9 Tenant Centric Service Pack Sample
This is an example Tenant Centric UI.  It includes everything needed to build a running UI and deploy metadata to MSX through the Component Management function of the MSX Product so that it appears as a purchaseable service.  This example is based on Angular 9.  It does expect that you have basic HTML, CSS, and JavaScript knowledge, as well as understanding of Angular 9.

## Development Environment

The preferred development environment is either macOS or Linux.  These provide all the general shell commands and scripting support to make startup and development simpler.  Other environments will work, but they are not covered in this document.

## Dependencies

The MSX UI has a few dependencies that must be installed before you can start or develop code for the UI.


### macOS

To develop on macOS, it is highly recommended that you have Homebrew installed.
Homebrew allows easy install of many common unix utilities on macOS that make
development easier.  It is also recommended that you have the command line tools
for Xcode installed.

1. [Homebrew](https://brew.sh/) - A utility for installing all sorts of
   UNIX/Linux utilities on macOS. It is highly recommended that you install
   this, as all the other dependencies can be installed through Homebrew.


### Linux

Nothing special is required to be installed on Linux specifically.


### Common

1. [NPM](https://www.npmjs.com/) - The Node Package Manager.  Node.js is used to execute the build system, manage dependencies.  This comes automatically with Node.js, can be downloaded from the web and installed, or installed through your preferred Linux package manager when on Linux.  It's best to use the version that comes with nodejs from nodejs.org in the package they provide.  Version 6.14.15 is known to work well, and comes with node v14.18.1.
2. [Node.js](https://nodejs.org/en/) - A commandline JavaScript runtime.  It is used for running the builds and local webserver.  This can be installed through Homebrew (`brew install node@10`), downloaded from the web and installed, or installed through your preferred Linux package manager when on Linux.  Any version equal to or greater than 8, but less than 15 should be okay.  
3. [Webpack](https://webpack.js.org/) - Webpack is used for generating the final set of files for the service and the metadata to deploy.  It also handles processing Javascript to be compatible across browsers with Babel, as well as allowing usage of type script.  This can be installed through NPM (`npm install -g webpack webpack-cli`), downloaded from the web and installed, or installed through your preferred Linux package manager when on Linux
4. [Docker](https://www.docker.com/) - You need docker for generating the SLM installable container.
5.  Curl - A useful utility for making web requests from the command line.  Used for debugging microservice calls.

## The Top Level Build Project
This level contains the basic build files used to construct a TC-UI build tree.  What it does is produce a UI build tree from a set of templated files to make it easy to quickly stub out an MSX UI installable UI build tree for development in.  The template generator consists of a Webpack configuration file that is designed copy and insert the configuration information you provide to create the new sub product UI for MSX.  

## Generating the new UI Project
Building the new UI project is simple. You just execute a shell script and pass it a few arguments.

> **Gotcha**
>
> The arguments `-image-file` and `-output-dir` must be absolute paths.

Here is an example of calling the script as user `johndoe`.

```shell
./createTemplate.sh -project-name="fakecoSomeNewService" \
-project-description="My Awesome UI for SomeNewService" \
-image-file="/Users/johndoe/Images/image.png" \
-output-dir="/Users/johndoe/Projects/SomeNewDirectory"
```



This will execute the project creator and generate a new UI build tree in the provided output directory.


# Details

Generator is webpack based. Files from `template` folder are copied into new folder using
`webpack` and its configuration file `webpack.config.js`.
Parameters passed in command line are used for placeholder replacement in the template files.  
Resulting folder contains MSX deployable Service Pack with the name passed in `project-name` parameter. Service Pack names must be unique within one MSX installation.

Deployment result is new Service is visible in MSX Service Catalog.

Taking example above 
```shell
./createTemplate.sh -project-name="fakecoSomeNewService" \
-project-description="My Awesome UI for SomeNewService" \
-image-file="/Users/johndoe/Images/image.png" \
-output-dir="/Users/johndoe/Projects/SomeNewDirectory"
```
this will generate code for new Service Pack called `fakecoSomeNewService`.

Once code is compiled, packaged and deployed MSX Service Catalog will show new Service 
`My Awesome UI for SomeNewService`.

See *README.md* file in the generated code folder for further details on using and customizing the Service Pack UI.
