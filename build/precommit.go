package build

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"github.com/bmatcuk/doublestar"
	"os"
	"path/filepath"
)

func init() {
	AddTarget("go-fmt", "Format all go source files", GoFmt)
}

func GoFmt(args []string) error {
	return exec.ExecutePipes(
		exec.Exec("go", []string{"fmt"}, findGoFiles()))
}

func findGoFiles() []string {
	var dirMap = make(map[string]struct{})
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		isVendor, _ := doublestar.Match("vendor/**/*", path)
		if isVendor {
			return nil
		}
		isGoFile, _ := doublestar.Match("**/*.go", path)
		if !isGoFile {
			return nil
		}
		dirName := filepath.Dir(path)
		dirMap[dirName] = struct{}{}
		return nil
	})

	var results []string
	for dirName := range dirMap {
		results = append(results, "./"+dirName)
	}
	return results
}
