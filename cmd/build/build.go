package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/build"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"gopkg.in/pipe.v2"
)

func init() {
	build.AddTarget("update-skel-templates", "Update skeleton templates", UpdateSkelTemplates)
}

func UpdateSkelTemplates(args []string) error {
	return exec.ExecutePipesIn("skel", pipe.Exec("go", "generate"))
}

func main() {
	build.Run()
}
