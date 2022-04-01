// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import (
	"bufio"
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"encoding/json"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/spf13/cobra"
	"gopkg.in/pipe.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func init() {
	AddTarget("go-fmt", "Format all go source files", GoFmt)
	AddTarget("go-vet", "Vet all go source files", GoVet)

	licenseCmd := AddTarget("license", "License all go source files", LicenseHeaders)
	licenseCmd.Args = cobra.NoArgs
	licenseCmd.Flags().Bool("check", false, "Verify license header has been applied")

}

func GoFmt(_ []string) error {
	return exec.ExecutePipes(
		exec.Exec("go", []string{"fmt"}, findGoFileDirectories()))
}

func GoVet(_ []string) (err error) {
	vetOptions := []string{"vet"}
	vetOptions = append(vetOptions, BuildConfig.Go.Vet.Options...)
	vetOptions = append(vetOptions, "-json")

	sourceDirectories := findGoFileDirectories()

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
				entry.Position = strings.TrimPrefix(entry.Position, wd+"/")
				vetResults = append(vetResults, entry)
			}
		}
	}

	return
}

func findGoFileDirectories() []string {
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

const thirdPartyLicenseHeaderFile = "HEADER"

var licenseHeadersMissing bool = false

func LicenseHeaders(_ []string) error {
	check, _ := BuildConfig.Cfg.BoolOr("cli.flag.check", false)
	if check {
		logger.Infof("Validating license headers")
	} else {
		logger.Infof("Applying license headers")
	}

	for _, goFileDirectory := range findGoFileDirectories() {
		skippedDir := isSkippedDir(goFileDirectory)
		if skippedDir {
			// do nothing
			continue
		}

		isThirdPartyLicense, err := isThirdPartyDir(goFileDirectory)
		if err != nil {
			return err
		}
		if isThirdPartyLicense {
			err = applyThirdPartyLicenseToDir(goFileDirectory, check)
			if err != nil {
				return err
			}
			continue
		}

		// Internal license header
		err = applyCiscoLicenseToDir(goFileDirectory, check)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyCiscoLicenseToDir(dir string, check bool) error {
	fileNames, err := doublestar.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		err = applyCiscoLicenseToFile(fileName, check)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyCiscoLicenseToFile(fileName string, check bool) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if hasCiscoLicense(b) || isGenerated(b) {
		return nil
	} else if check {
		logger.Errorf("%q missing Cisco license header", fileName)
		licenseHeadersMissing = true
		return nil
	}

	headerBytes := ciscoLicenseHeader()
	fileBytes := append(headerBytes, b...)
	return ioutil.WriteFile(fileName, fileBytes, 0644)
}

func ciscoLicenseHeader() []byte {
	year := time.Now().Year()
	lines := []string{
		"// " + fmt.Sprintf(`Copyright © %d, Cisco Systems Inc.`, year),
		"// Use of this source code is governed by an MIT-style license that can be",
		"// found in the LICENSE file or at https://opensource.org/licenses/MIT.",
		"",
	}

	buf := bytes.Buffer{}
	for _, line := range lines {
		buf.WriteString(line)
		buf.Write([]byte("\n"))
	}

	return buf.Bytes()
}

func hasCiscoLicense(data []byte) bool {
	bufReader := bufio.NewReader(bytes.NewReader(data))
	firstLine, _, _ := bufReader.ReadLine()
	firstLine = bytes.ToLower(firstLine)
	return bytes.Contains(firstLine, []byte("copyright")) && bytes.Contains(firstLine, []byte("cisco"))
}

var generatedRegExp = regexp.MustCompile(`(?m)^.{1,2} Code generated .* DO NOT EDIT\.$`)

func isGenerated(data []byte) bool {
	return generatedRegExp.Match(data)
}

func isSkippedDir(p string) bool {
	for _, excludeGlob := range BuildConfig.License.Excludes {
		if matches, err := doublestar.Match(excludeGlob, p); err != nil {
			return false
		} else if matches {
			return true
		}
	}

	return false
}

func applyThirdPartyLicenseToDir(dir string, check bool) error {
	licenseHeader, err := thirdPartyLicenseHeader(dir)
	if err != nil {
		return err
	}

	goFileNames, err := doublestar.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return err
	}

	for _, fileName := range goFileNames {
		err = applyThirdPartyLicenseToFile(fileName, licenseHeader, check)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyThirdPartyLicenseToFile(fileName string, headerBytes []byte, check bool) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if hasThirdPartyLicense(b, headerBytes) || isGenerated(b) {
		return nil
	} else if check {
		logger.Errorf("%q missing third-party license header", fileName)
		licenseHeadersMissing = true
		return nil
	}

	fileBytes := append(headerBytes, b...)
	return ioutil.WriteFile(fileName, fileBytes, 0644)
}

func hasThirdPartyLicense(fileBytes, headerBytes []byte) bool {
	if len(headerBytes) > len(fileBytes) {
		return false
	}

	return bytes.Equal(fileBytes[:len(headerBytes)], headerBytes)
}

func thirdPartyLicenseHeader(dir string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(dir, "HEADER"))
}

var thirdPartyLicensesFound = make(map[string]bool)

func isThirdPartyDir(dir string) (bool, error) {
	dir = filepath.Clean(dir)

	if found, ok := thirdPartyLicensesFound[dir]; ok {
		return found, nil
	}

	fullPath := filepath.Join(dir, thirdPartyLicenseHeaderFile)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		thirdPartyLicensesFound[dir] = false
		return false, nil
	} else if err != nil {
		return false, err
	}

	thirdPartyLicensesFound[dir] = true
	return true, nil
}
