# Skel Tests


*Copyright © 2023, Cisco Systems Inc.
Use of this source code is governed by an MIT-style license that can be
found in the LICENSE file or at https://opensource.org/licenses/MIT.*

This directory contains resources to build and execute tests of
the `skel` generator program.

## Components

- `txtartools`
    - `txtarwrap` is a program for making source archives from trees of files [txtar format](https://pkg.go.dev/golang.org/x/tools/txtar)
    - `txtargen` is for making [testscript](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript) scripts that test trees of files or compare them with other trees  
    - `clienv` a simple and primitive lib to allow fallback from cli flag to env var without having to resort to Viper etc.
- `fixtures` 
    - `build` contains script fragments and scripts that combine them with assets from `golden` to produce final test scripts in `final`
      - `make-normal-tests` contains a golang program that builds the bulk of the tests
      - `testlist` is a list of tests, intended to be used by both make-normal-tests and skel_test.go
    - `golden` contains "golden" known-good file trees
    - `before` contains file sets to be provided to the program under test before it runs
    - `final` contains the final self-contained test scripts & archives
      - `root` contains the root test txtar and a script to run it after manually running `skel`
    - `difficulties` contains scripts that need more work :(

All following sections presume that you have txtartools and testscript installed and in your path, and that you have set your FIXT env var:  
```sh
    make install-txtartools  # in the tests dir
    go install github.com/rogpeppe/go-internal/cmd/testscript@latest
    export FIXT= # full path to your fixtures directory, likely in a git repo
```
## Building & Updating

After changes to txtarwrap or txtargen you will want to `make install-txtartools`
After changes to skel you will need to `make install-skel` in the root of go-msx then perform the steps below as applicable.

### Rebuilding the root assets and tests after changes to `skel`

We can't automatically run the full root command because of its menus, so we must run `skel` root **manually**, supplying the default values for each prompt. Then we can build tests for execution later.  

1. cd to an empty directory somewhere safe
2. Run: `skel` with no parameters and enter defaults at all prompts by hitting enter at each
3. This will create the root set of files for the `someservice` example service
4. In the same dir, run:
    ```sh
    $FIXT/build/make-roots.sh
    ```
5. This will update the root test script/archive in `$FIXT/final/root/root-test.txtar`

### Rebuilding special tests

Some tests have slightly wierd needs, they may be rebuilt using these scripts:

1. Completions test: run `make-completions-test.sh`
2. Version test: is currently version agnostic and thus never needs updating (it is `final/version-test.txtar`)
3. Help tests: Oddly, the tests of the help system are presently run by `skel/skel_test.go`, they simply test for non-error running and thus should not need updating as long as `testlist` is maintained.

### Rebuilding tests for most other targets

The program `go run build/make-normal-tests/make-normal-tests.go` builds all working tests that it can.

To rebuild a particular test use:
    ```shell
    $FIXT/build/make-any.sh #test-name
    ```
This builds an intermediate txtar that loads a plain root set of files and runs skel against them to produce a new golden set, which it outputs to `golden/`, while putting a new self-contained test txtar in `final/`

## Testing

### Root Test

To test whether a new version of `skel` produces the same root command output as previously: 

1. cd to an empty directory somewhere safe
2. Run: `export FIXT=<full path to your fixtures directory, likely in a git repo>`
3. Run: `skel` with no parameters and enter defaults at all prompts by hitting enter at each
4. This will create the root set of files for the `someservice` example service
5. Run the test script against the file tree you just made:
    ```sh
    $FIXT/final/root/root-test.sh
    ``` 
It will report errors on any files it made that are different from the golden files collected when the test was last built. If there are unexpected changes you will need to address them.

### All other Tests

All the txtars found in final, except in those in final/root may be run standalone by feeding them to testscript:
```shell
    cd final
    testscript *.txtar
```