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

func variables() map[string]string {
	return map[string]string{
		"app.name": skeletonConfig.AppName,
		"server.port": strconv.Itoa(skeletonConfig.ServerPort),
		"server.contextpath": "/" + skeletonConfig.AppName,
		"kubernetes.group": "platformms",
	}
}

func renderTemplate(templateName string) error {
	f, ok := staticFiles[templateName]
	if !ok {
		return errors.Errorf("Template file not found: %s", templateName)
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

	rendered := string(bytes)
	variableValues := variables()
	variableInstanceRegex := regexp.MustCompile(`\${([^}]+)}`)
	for _, variableInstance := range variableInstanceRegex.FindAllStringSubmatch(rendered, -1) {
		variableName := variableInstance[1]
		variableValue, ok := variableValues[strings.ToLower(variableName)]
		if ok {
			rendered = strings.ReplaceAll(rendered, "${" + variableName + "}", variableValue)
		}
	}

	bytes = []byte(rendered)
	targetFileName := path.Join(skeletonConfig.TargetDirectory(), templateName)
	targetDirectory := path.Dir(targetFileName)
	err = os.MkdirAll(targetDirectory, 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(targetFileName, bytes, 0644)
}
