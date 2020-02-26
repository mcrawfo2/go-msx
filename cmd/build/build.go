package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/build"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
)

func init() {
	build.AddTarget("update-skel-templates", "Update skeleton templates", UpdateSkelTemplates)
}

func UpdateSkelTemplates(args []string) error {
	return exec.ExecutePipes(
		exec.WithDir("skel",
			exec.Exec("go", []string{"generate"})))
}

func main() {
	build.Run()
}
