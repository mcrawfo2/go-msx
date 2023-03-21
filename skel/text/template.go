package text

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"regexp"
	"strings"
)

type TemplateOptions struct {
	Variables  map[string]string // key is var name
	Conditions map[string]bool   // key is condition name
	Strings    map[string]string // key is string name
}

func NewTemplateOptions() TemplateOptions {
	return TemplateOptions{
		Variables:  make(map[string]string),
		Conditions: make(map[string]bool),
		Strings:    make(map[string]string),
	}
}

type TemplateLanguage int

const (
	TemplateLanguageSkel TemplateLanguage = iota
	TemplateLanguageGoText
)

type Template struct {
	Name     string
	Loader   TemplateLoader
	Format   FileFormat
	Language TemplateLanguage
}

func (t Template) Render(options TemplateOptions) (result string, err error) {
	// Load the source
	contents, err := t.Loader()
	if err != nil {
		return
	}

	if len(options.Strings) > 0 {
		contents = t.substituteStrings(contents, options.Strings)
	}

	// Substitute variables
	if len(options.Variables) > 0 {
		contents = t.substituteVariables(contents, options.Variables)
	}

	// Execute conditions
	for condition, value := range options.Conditions {
		contents, err = t.processConditionalBlocks(contents, condition, value)
	}
	if err != nil {
		return
	}

	result = string(contents)
	return
}

func (t Template) substituteStrings(contents []byte, subs map[string]string) []byte {
	// Replace strings
	for sourceString, destString := range subs {
		contents = bytes.ReplaceAll(contents, []byte(sourceString), []byte(destString))
	}
	return contents
}

func (t Template) substituteVariables(source []byte, variableValues map[string]string) []byte {
	rendered := source
	variableInstanceRegex := regexp.MustCompile(`\${([^}]+)}`)
	for _, variableInstance := range variableInstanceRegex.FindAllSubmatch(rendered, -1) {
		variableName := string(variableInstance[1])
		variableValue, ok := variableValues[strings.ToLower(variableName)]
		if ok {
			rendered = bytes.ReplaceAll(rendered,
				[]byte("${"+variableName+"}"),
				[]byte(variableValue))
		}
	}
	return rendered
}

func (t Template) conditionalMarkers() (string, string) {
	prefix, suffix := t.Format.CommentMarkers()
	return prefix + "#", suffix
}

func (t Template) processConditionalBlocks(data []byte, condition string, output bool) (result []byte, err error) {
	type parserState int
	const outside parserState = 0
	const insideIf parserState = 1
	const insideElse parserState = 2

	sb := bytes.Buffer{}
	write := func(out bool, line string) {
		if !out {
			return
		}
		sb.WriteString(line)
		sb.WriteRune('\n')
	}
	insideCondition := outside
	prefix, suffix := t.conditionalMarkers()
	startMarker := prefix + "if " + condition + suffix
	middleMarker := prefix + "else " + condition + suffix
	endMarker := prefix + "endif " + condition + suffix

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		lineTrimmed := strings.TrimSpace(line)
		switch insideCondition {
		case outside:
			switch lineTrimmed {
			case startMarker:
				insideCondition = insideIf
			default:
				write(true, line)
			}

		case insideIf:
			switch lineTrimmed {
			case endMarker:
				insideCondition = outside
			case middleMarker:
				insideCondition = insideElse
			default:
				write(output, line)
			}

		case insideElse:
			switch lineTrimmed {
			case endMarker:
				insideCondition = outside
			default:
				write(!output, line)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "Failed to process conditional blocks")
	}

	return sb.Bytes(), nil
}

func NewTemplate(name string, format FileFormat, language TemplateLanguage, source TemplateLoader) Template {
	var result = Template{
		Name:     name,
		Loader:   source,
		Format:   format,
		Language: language,
	}

	return result
}

type TemplateLoader func() ([]byte, error)

func TemplateJsonOption(value interface{}) TemplateLoader {
	return func() ([]byte, error) {
		return json.Marshal(value)
	}
}

func TemplateJsonPrettyOption(value interface{}, indent string) TemplateLoader {
	return func() ([]byte, error) {
		return json.MarshalIndent(value, "", indent)
	}
}

func TemplateFileOption(fileName string) TemplateLoader {
	return func() ([]byte, error) {
		return os.ReadFile(fileName)
	}
}

func TemplateStringOption(content string) TemplateLoader {
	return func() ([]byte, error) {
		return []byte(content), nil
	}
}

func TemplateBytesOption(source []byte) TemplateLoader {
	return func() ([]byte, error) {
		return source, nil
	}
}
