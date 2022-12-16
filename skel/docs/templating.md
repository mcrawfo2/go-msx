# Using Skel template functions

## General Structure

Skel uses `targets`, which are groups of associated generating actions. Typically, a target will perform substitutions on a set of template files, emit them to appropriate directories in the generated tree and may perform shell functions (e.g. `git` etc.) thereafter to complete the generated application.

Internally, targets are identified by unique strings which allow target lists to be manipulated easily. Each target generally also has a corresponding `skel` cli subcommand which will execute it.

## File types

Skel can template any text file-type, but specifically recognizes the extensions for go, make, json, sql, yaml, groovy, properties, md, go-mod, docker, shell, js, ts and jenkins files.


## Substitution

Skel does substitutions into template files in three phases:

1. It substitutes particular `Strings` verbatim; these are passed into the rednering functions via the RenderOptions struct
2. It then substitutes variable values for the text in the templates matching `${variable}`. e.g. application name would be substituted for `${app.name}`. The possible variables are listed in skel/render.go around line 95.
3. It then evaluates conditional blocks (see below) 

For example, this piece of dockerfile:

```dockerfile
FROM ${BASE_IMAGE}
EXPOSE ${server.port}
EXPOSE ${debug.port}

ENV SERVICE_BIN "/usr/bin/${app.name}"
COPY --from=debug-builder /app/dist/root/ /
COPY --from=debug-builder /go/bin/dlv /usr/bin
```

An easier option than pawing around in the source code: available substitutions and conditionals may be listed by executing a `skel` generation with the debug or trace loglevels: `skel -l=DEBUG` -- they will be printed as part of the render options log lines.


## Conditional Blocks

These are defined by conditional markers and the words `if`, `else` and `endif`. Conditional markers vary by file type:

- make, yaml, properties, docker, bash: `#`
- sql: `--#`
- xml, md: `<--#, -->`
- everything else: `//#`

For example, the following block in a Makefile includes different lists depending on the app archetype:

```makefile
#if GENERATOR_APP
all: clean deps vet test docker assemblies deployment manifest
#endif GENERATOR_APP
#if GENERATOR_BEAT
all: clean deps vet test docker deployment manifest package
#endif GENERATOR_BEAT
#if GENERATOR_SP
all: clean deps vet test docker assemblies deployment manifest package
#endif GENERATOR_SP
```

## File Operations

As a cross-check mechanism, Skel is provided with the ability to insist whether files already exist or not during generation, and to halt if something is unexpected. The options are:

- Add: either add the file or replace it, we care not
- New: must not exist before, halt if it does
- AddNoOverwrite: May add it, or skip it, but don't halt
- Replace: must exist and we replace it
- Delete: must exist (or we halt) then we delete it
- Gone: delete it if it exists, don't halt

## File Names

Each template file has a source filename, which is in the embedded filesystem of templates, and a destination filename, which is relative to the root of the generated project.

Variables may be substituted into filenames, of either type, using the same syntax as within templates. e.g. `"local/${app.name}.remote.yml"`

If a dest file name is not provided, it is assumed to the same as the source.