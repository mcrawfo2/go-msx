//go:generate staticfiles -o templates.go templates/

package skel

import (
	"compress/gzip"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Template struct {
	SourceFile string
	DestFile   string
}

type StaticFile struct {
	Data     []byte
	DestFile string
}

func variables() map[string]string {
	return map[string]string{
		"app.name":                     skeletonConfig.AppName,
		"app.description":              skeletonConfig.AppDescription,
		"app.displayname":              skeletonConfig.AppDisplayName,
		"app.version":                  skeletonConfig.AppVersion,
		"app.migrateversion":           skeletonConfig.AppMigrateVersion(),
		"app.packageurl":               skeletonConfig.AppPackageUrl(),
		"server.port":                  strconv.Itoa(skeletonConfig.ServerPort),
		"server.contextpath":           path.Clean("/" + skeletonConfig.ServerContextPath),
		"kubernetes.group":             "platformms",
		"target.dir":                   skeletonConfig.TargetDirectory(),
		"repository.cassandra.enabled": strconv.FormatBool(skeletonConfig.Repository == "cassandra"),
		"repository.cockroach.enabled": strconv.FormatBool(skeletonConfig.Repository == "cockroach"),
	}
}

func writeStaticFiles(files map[string]StaticFile) error {
	for message, static := range files {
		logger.Infof("- %s (%s)", message, static.DestFile)
		err := writeStatic(static)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeStatic(static StaticFile) (err error) {
	variableValues := variables()

	destFile := static.DestFile
	if destFile == "" {
		return errors.New("Static file missing destination filename")
	}
	destFile = substituteVariables(destFile, variableValues)

	targetFileName := path.Join(skeletonConfig.TargetDirectory(), destFile)
	targetDirectory := path.Dir(targetFileName)
	err = os.MkdirAll(targetDirectory, 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(targetFileName, static.Data, 0644)
}

func renderTemplates(templates map[string]Template) error {
	for message, template := range templates {
		logger.Infof("- %s (%s)", message, template.SourceFile)
		err := renderTemplate(template)
		if err != nil {
			return err
		}
	}
	return nil
}

func readTemplate(sourceFile string) ([]byte, error) {
	f, ok := staticFiles[sourceFile]
	if !ok {
		return nil, errors.Errorf("Template file not found: %s", sourceFile)
	}

	var reader io.Reader
	if f.size != 0 {
		var err error
		reader, err = gzip.NewReader(strings.NewReader(f.data))
		if err != nil {
			return nil, err
		}
	} else {
		reader = strings.NewReader(f.data)
	}

	return ioutil.ReadAll(reader)
}

func renderTemplate(template Template) error {
	sourceFile := template.SourceFile
	bytes, err := readTemplate(template.SourceFile)
	if err != nil {
		return err
	}

	variableValues := variables()

	destFile := template.DestFile
	if destFile == "" {
		destFile = sourceFile
	} else {
		destFile = substituteVariables(destFile, variableValues)
	}

	rendered := substituteVariables(string(bytes), variableValues)
	bytes = []byte(rendered)
	targetFileName := path.Join(skeletonConfig.TargetDirectory(), destFile)
	targetDirectory := path.Dir(targetFileName)
	err = os.MkdirAll(targetDirectory, 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(targetFileName, bytes, 0644)
}

func substituteVariables(source string, variableValues map[string]string) string {
	rendered := source
	variableInstanceRegex := regexp.MustCompile(`\${([^}]+)}`)
	for _, variableInstance := range variableInstanceRegex.FindAllStringSubmatch(rendered, -1) {
		variableName := variableInstance[1]
		variableValue, ok := variableValues[strings.ToLower(variableName)]
		if ok {
			rendered = strings.ReplaceAll(rendered, "${"+variableName+"}", variableValue)
		}
	}
	return rendered
}
