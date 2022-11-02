# Skel Execution Sequence

## Program flow

1. `cmd/skel/skel.go` is skel's main function. It calls the run method from the 
skel package (in `skel/`), passing in a build number, which, in turn, is passed
in during build via `ldflags` from `build-tool`'s `BuildConfig` (`build-tool` is in `build/tool.go`)  
2. the skel package init loads skel templates into a map from the embedded filesystem
3. it also sets up pre-run project and generation configuration loading routines that will be called before each command is run
4. config is loaded into a package-level variable called skeletonConfig
5. if skel does not find a pre-existing project in the current dir (by looking for `.skel.json`) it diverts to the interactive menus
6. the menus are navigated in `skel/configure.go` which simply fills in the values into skeletonConfig which is declared therein 
7. skel then routes on via the `cli/` package in order to enable cobra command processing (`github.com/spf13/cobra`)
8. (aside: the actual skel generation targets are defined and registered in `skel/skeleton.go` which comprises the heart of skel)
9. each target/command is called using its string key after the machinery in skeleton.go adds common pre and post generator keys to the list of those to be run
10. The sandwich filling between the pre and post slices of generators is derived from the archetype (see `skel/archetype.go`)
11. in addition to pre and post generators, there are some that must be run explicitly by using cli commands (see [groupings](#generator-groupings) )
12. once the complete list of generators to be run has been assembled, they are executed in order
13. each target/command/generator will typically:
    a. call out to OS level routines using the `gopkg.in/pipe.v2` lib
    b. fill out substitutions and apply templates using the routines in `render.go` which provides many options
14. *fin*

In general, `skel` simply creates new files, overwriting any that might exist before them; however, most of the targets do not overlap, and may be freely run in any order, or, in the case of Domain and AsyncAPI may be run mutiple times to build up additional variants.


## Domain generation

## AsyncAPI generation

## Generator groupings

(sp = service pack)

### Pre generators
`generate-skel-json`            Create the skel configuration file  
`generate-build`                Create makefile, build.go, build.yml  
`generate-app`                  Readme, go.mod, main.go, 2 yml
`generate-test`                 internal/empty_test.go

### Archetype specific
`generate-migrate`              (beat,sp) Create the migrate package  
`generate-domain-beats`         (beat) Generate beats domain implementation  
`generate-service-pack`         (sp) Generate service pack implementation  
`generate-kubernetes`           (all) Create production kubernetes manifest templates  

### Post
`generate-deployment-variables` deployment_variables.yml    
`add-go-msx-dependency `        Adds go modules appropriate to the archetype    
`generate-local        `        local&remote address, consul/vault    
`generate-manifest     `        installer manifest (maven)  
`generate-dockerfile   `        dockerfile: build and distn  
`generate-goland       `        Create a Goland project for the application  
`generate-vscode       `        Create a VSCode project for the application  
`generate-jenkins      `        Create Jenkins CI templates  
`generate-github       `        Create github configuration files  
`generate-git          `        Create git repository  

### Called explicitly
`completion`                    autocompletion script for the specified shell  
`generate-certificate`          Generate an X.509 server certificate and private key  
  
`generate-channel`              Create async channel  
`generate-channel-asyncapi`     Create stream from AsyncApi 2.4 specification  
`generate-channel-publisher`    Create async channel publisher  
`generate-channel-subscriber`   Create async channel subscriber  
  
`generate-domain-openapi`       Create domains from OpenAPI 3.0 manifest  
`generate-domain-system `       Generate system domain implementation  
`generate-domain-tenant`        Generate tenant domain implementation  
  
`generate-timer`                Generate timer implementation  
`generate-topic-publisher`      Generate publisher topic implementation  
`generate-topic-subscriber`     Generate subscriber topic implementation  
`generate-webservices`          Create web services from swagger manifest  
