package build

import (
	"bufio"
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"encoding/json"
	"github.com/bmatcuk/doublestar"
	"gopkg.in/pipe.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func init() {
	AddTarget("go-fmt", "Format all go source files", GoFmt)
	AddTarget("go-vet", "Vet all go source files", GoVet)
}

func GoFmt(_ []string) error {
	return exec.ExecutePipes(
		exec.Exec("go", []string{"fmt"}, findGoFiles()))
}

func GoVet(_ []string) (err error) {
	vetOptions := []string{"vet"}
	vetOptions = append(vetOptions, BuildConfig.Go.Vet.Options...)
	vetOptions = append(vetOptions, "-json")

	sourceDirectories := findGoFiles()

	var vetResults = new(bytes.Buffer)
	s := pipe.NewState(ioutil.Discard, vetResults)
	p := exec.Exec("go", vetOptions, sourceDirectories)

	if err = p(s); err == nil {
		err = s.RunTasks()
	}
	if err != nil {
		return err
	}

	return outputGoVetResults(vetResults.Bytes())
}

type VetResult struct {
	Package  string `json:"package"`
	Tool     string `json:"tool"`
	Position string `json:"posn"`
	Message  string `json:"message"`
}

func outputGoVetResults(results []byte) (err error) {
	if _, err = os.Stdout.Write(results); err != nil {
		return err
	}

	var vetResultsJsonPath = path.Join(BuildConfig.DistPath(), "vet.json")
	var resultsReader = bytes.NewReader(results)
	var scanner = bufio.NewScanner(resultsReader)

	var vetResults []VetResult

	// Read comment line
	if !scanner.Scan() {
		return
	}

	for err == nil {
		vetResultBuffer := new(bytes.Buffer)
		for scanner.Scan() {
			lineBytes := scanner.Bytes()
			if len(lineBytes) > 0 && lineBytes[0] == '#' {
				break
			}
			vetResultBuffer.Write(lineBytes)
		}

		results = vetResultBuffer.Bytes()
		if len(results) > 0 {
			var singleVetResults []VetResult
			singleVetResults, err = parseGoVetResults(results)
			if err != nil {
				return
			}
			vetResults = append(vetResults, singleVetResults...)
		} else {
			break
		}
	}

	vetResultBytes, err := json.Marshal(vetResults)
	if err != nil {
		return err
	}

	err = os.MkdirAll(BuildConfig.DistPath(), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(vetResultsJsonPath, vetResultBytes, 0644)
}

func parseGoVetResults(results []byte) (vetResults []VetResult, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	// Read json body
	var m map[string]map[string][]VetResult
	if err = json.Unmarshal(results, &m); err != nil {
		return
	}

	// Flatten
	for pkg, pkgEntry := range m {
		for tool, toolEntries := range pkgEntry {
			for _, entry := range toolEntries {
				entry.Package = pkg
				entry.Tool = tool
				entry.Position = strings.TrimPrefix(entry.Position, wd + "/")
				vetResults = append(vetResults, entry)
			}
		}
	}

	return
}

func findGoFiles() []string {
	var dirMap = make(map[string]struct{})
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if isVendor, _ := doublestar.Match("vendor/**/*", path); isVendor {
			return nil
		}

		if isSkelTemplate, _ := doublestar.Match("skel/_templates/**/*", path); isSkelTemplate {
			return nil
		}

		if isGoFile, _ := doublestar.Match("**/*.go", path); !isGoFile {
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
