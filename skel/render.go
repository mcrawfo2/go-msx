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

func variables() map[string]string {
	return map[string]string{
		"app.name":           skeletonConfig.AppName,
		"app.description":    skeletonConfig.AppDescription,
		"app.displayname":    skeletonConfig.AppDisplayName,
		"server.port":        strconv.Itoa(skeletonConfig.ServerPort),
		"server.contextpath": "/" + skeletonConfig.AppName,
		"kubernetes.group":   "platformms",
		"target.dir":         skeletonConfig.TargetDirectory(),
	}
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

func renderTemplate(template Template) error {
	sourceFile := template.SourceFile
	f, ok := staticFiles[sourceFile]
	if !ok {
		return errors.Errorf("Template file not found: %s", sourceFile)
	}

	var reader io.Reader
	if f.size != 0 {
		var err error
		reader, err = gzip.NewReader(strings.NewReader(f.data))
		if err != nil {
			return err
		}
	} else {
		reader = strings.NewReader(f.data)
	}

	bytes, err := ioutil.ReadAll(reader)
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
