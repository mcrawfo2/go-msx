package build

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"github.com/bmatcuk/doublestar"
	"gopkg.in/pipe.v2"
	"os"
	"path"
	"path/filepath"
)

const permsDir = 0755
const permsFile = 0644

func init() {
	AddTarget("download-test-deps", "Download test dependencies", InstallTestDependencies)
	AddTarget("execute-unit-tests", "Execute unit tests", ExecuteUnitTests)
}

func InstallTestDependencies(args []string) error {
	script := pipe.Script(
		exec.Info("Downloading test dependencies"),
		goGet("github.com/axw/gocov/gocov"),
		goGet("github.com/AlekSi/gocov-xml"),
		goGet("github.com/stretchr/testify/assert"),
		goGet("github.com/stretchr/testify/mock"),
		goGet("github.com/stretchr/testify/http"),
		goGet("github.com/pmezard/go-difflib/difflib"),
		goGet("github.com/jstemmer/go-junit-report"),
		pipe.Write(os.Stdout),
	)
	return pipe.Run(script)
}

func ExecuteUnitTests(args []string) error {
	testFile := func(parts ...string) string {
		return path.Join(append([]string{BuildConfig.TestPath()}, parts...)...)
	}

	logger.Info("Locating testable directories")
	var testableDirectories = locateTestableDirectories()
	if len(testableDirectories) == 0 {
		logger.Warn("No testable directories found.")
		return nil
	}

	var testResults = new(bytes.Buffer)
	var goCoverOutPath = testFile("gocover.out")
	var goCoverHtmlPath = testFile("gocover.html")
	var junitReportXmlPath = testFile("junit-report.xml")
	var coberturaCoverageXmlPath = testFile("cobertura-coverage.xml")

	return exec.ExecutePipes(
		pipe.Line(
			exec.Info("Recreating test result directory"),
			exec.RemoveAll(testFile()),
			pipe.MkDirAll(testFile(), permsDir),
		),
		pipe.Line(
			exec.Info("Executing unit tests"),
			exec.Exec("go", []string{"test", "-coverprofile=" + goCoverOutPath, "-v"}, testableDirectories),
			pipe.Tee(os.Stdout),
			pipe.Write(testResults),
		),
		pipe.Line(
			exec.Info("Generating JUnit XML report"),
			pipe.Read(testResults),
			pipe.Exec("go-junit-report"),
			pipe.WriteFile(junitReportXmlPath, permsFile),
		),
		pipe.Line(
			exec.Info("Generating HTML coverage report"),
			pipe.Exec("go", "tool", "cover", "-html="+goCoverOutPath, "-o", goCoverHtmlPath),
			pipe.Write(os.Stdout),
		),
		pipe.Line(
			exec.Info("Generating Cobertura XML coverage report"),
			exec.Exec("gocov", []string{"test"}, testableDirectories),
			pipe.Exec("gocov-xml"),
			pipe.WriteFile(coberturaCoverageXmlPath, permsFile),
		),
	)
}

func locateTestableDirectories() []string {
	var testDirMap = make(map[string]struct{})
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		isVendor, _ := doublestar.Match("vendor/**/*", path)
		if isVendor {
			return nil
		}
		isTestFile, _ := doublestar.Match("**/*_test.go", path)
		if !isTestFile {
			return nil
		}
		dirName := filepath.Dir(path)
		testDirMap[dirName] = struct{}{}
		return nil
	})

	var results []string
	for testDirName := range testDirMap {
		results = append(results, "./"+testDirName)
	}
	return results
}

func goGet(packageName string) pipe.Pipe {
	return pipe.Exec("go", "get", packageName)
}
