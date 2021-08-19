package build

import (
	"bufio"
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

const thirdPartyLicenseHeaderFile = "HEADER"

func init() {
	AddTarget("go-fmt", "Format all go source files", GoFmt)
	AddTarget("license", "License all go source files", LicenseHeaders)
}

func GoFmt(args []string) error {
	return exec.ExecutePipes(
		exec.Exec("go", []string{"fmt"}, findGoFileDirectories()))
}

func findGoFileDirectories() []string {
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
	sort.Strings(results)
	return results
}

func LicenseHeaders(args []string) error {
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
			err = applyThirdPartyLicenseToDir(goFileDirectory)
			if err != nil {
				return err
			}
			continue
		}

		// Internal license header
		err = applyCiscoLicenseToDir(goFileDirectory)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyCiscoLicenseToDir(dir string) error {
	fileNames, err := doublestar.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		err = applyCiscoLicenseToFile(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyCiscoLicenseToFile(fileName string) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if hasCiscoLicense(b) || isGenerated(b) {
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

func applyThirdPartyLicenseToDir(dir string) error {
	licenseHeader, err := thirdPartyLicenseHeader(dir)
	if err != nil {
		return err
	}

	goFileNames, err := doublestar.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return err
	}

	for _, fileName := range goFileNames {
		err = applyThirdPartyLicenseToFile(fileName, licenseHeader)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyThirdPartyLicenseToFile(fileName string, headerBytes []byte) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if hasThirdPartyLicense(b, headerBytes) || isGenerated(b) {
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
