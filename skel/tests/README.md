# Skel Integration Tests

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


## Rebuilding the expected outputs

After modifying skel to change its output, you will need to update the expected outputs:

```bash
make golden
```

## Testing the current outputs

To test whether a new version of `skel` produces the same command output as previously: 

```bash
make test
```
