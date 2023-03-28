# Skel Tests

*Copyright © 2023, Cisco Systems Inc.
Use of this source code is governed by an MIT-style license that can be
found in the LICENSE file or at https://opensource.org/licenses/MIT.*

This directory contains resources to build and execute tests of
the `skel` generator program.

## Unresolved Issues

1. Needs to be connected into CI/CD pipeline

## Components

  - `txtartools`
    - `txtarwrap` is a program for making source archives from trees of files. [Txtar format](https://pkg.go.dev/golang.org/x/tools/txtar).
    - `txtarunwrap` does the inverse of `txtarwrap`
    - `txtargen` is for making [testscript](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript) scripts that test trees of files or compares them with other trees  
    - `clienv` a simple and primitive lib to allow fallback from cli flag to env var without having to resort to Viper etc.
    - `txtarcmp` compares the files found in a txtar with their equivalents in the local filesystem
- `fixtures` 
    - `before` contains file sets to be provided to the program under test before it runs
    - `build` contains things that build files for use in testing
      - `make-tests` contains a golang program that builds the tests (except root)
      - `testlist` is a list of tests, and utilities that need to be common to both building and running
    - `final` contains the final self-contained test scripts & archives
      - `root` contains the root test txtar and a script to run it after manually running `skel`
      - `run-tests` contains a go program that runs all the tests (except root)
    - `golden` contains "golden" known-good file trees

All following sections presume that you have txtartools and testscript installed and in your path, and that you have set your FIXT env var:

```bash
    make install-txtartools  # in the tests dir    
    go install github.com/rogpeppe/go-internal/cmd/testscript@latest
    export FIXT= # full path to your fixtures directory, likely in a git repo
```

## Building & Updating

After changes to any of the txtartools you will want to `make install-txtartools`
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

### Rebuilding tests for everything except root

If the program behaviour changes the tests may fail. Once you have confirmed that the change result is correct, or have corrected it, rebuild the corresponding test(s).

`make make-tests` will run `go run build/make-tests/make-tests.go` to build all working tests. 

To rebuild a particular test(s) use:

```sh
    go run build/make-tests/make-tests.go test-name1 test-name2 ...
```

Test names may be found in $FIXT/build/testlist/testlist.go. At time of writing, all tests are 'happy path'.

### Adding new tests

If the test runs a `skel` command and outputs files, simply add a line for it in $FIXT/build/testlist/testlist.go. This file also contains examples of handling a few unusual cases.

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

`make run-tests` will run all the other tests. To run particular tests, use: 

```sh
    go run $FIXT/final/run-tests/run-tests.go test-name1 test-name2 ...
```

Test names may be found in $FIXT/build/testlist/testlist.go
