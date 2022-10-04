# Build Targets

go-msx Build has a collection of default build targets encompassing standard build steps that can be reused
in a project's Makefile.  The following chapters describe each of the standard build targets.

Users can also define custom build targets for project-specific needs.

## Custom Build Targets

Build targets can be added using the `build.AddTarget` function.  This will register a
CLI handler function for a new build target.  Build configuration can be accessed
vi the pkg.BuildConfiguration global variable from the handler function.  Ensure your
module containing the custom build target is initialized at startup by either:
- Defining your build target in the build `main` package of your project; or
- Importing the module containing your custom build target from the build `main` package
  of your project.

Example:

```go
package main

import build "cto-github.cisco.com/NFV-BU/go-msx-build/pkg"

var myCustomTargetFlag bool

func init() {
	cmd := build.AddTarget("custom-target", "A custom build target", customTarget)
	cmd.Flags().BoolVarP(&myCustomTargetFlag, "enabled", "e", false, "Custom target option")
}

func customTarget(args []string) error {
	// custom target steps go here
	if myCustomTargetFlag {
		// ...
    }
	return nil
}
```
