// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	gohtml "html/template"
	"os"
	"regexp"
	"strings"
	gotext "text/template"
)

type TemplateOptions struct {
	Variables  map[string]string // key is var name
	Conditions map[string]bool   // key is condition name
	Strings    map[string]string // key is string name
}

func (r TemplateOptions) AddString(source, dest string) {
	r.Strings[source] = dest
}

func (r TemplateOptions) AddStrings(strings map[string]string) {
	for k, v := range strings {
		r.Strings[k] = v
	}
}

func (r TemplateOptions) AddVariable(source, dest string) {
	r.Variables[strings.ToLower(source)] = dest
}

func (r TemplateOptions) AddVariables(variables map[string]string) {
	for k, v := range variables {
		r.Variables[strings.ToLower(k)] = v
	}
}

func (r TemplateOptions) AddCondition(condition string, value bool) {
	r.Conditions[condition] = value
}

func (r TemplateOptions) AddConditions(conditions map[string]bool) {
	for k, v := range conditions {
		r.Conditions[k] = v
	}
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
	TemplateLanguageGoHtml
)

type Template struct {
	Name         string
	Loader       TemplateLoader
	Format       FileFormat
	Language     TemplateLanguage
	Transformers Transformers
}

func (t Template) Render(options TemplateOptions) (result string, err error) {
	// Load the source
	contents, err := t.Loader()
	if err != nil {
		return
	}

	var resultBytes []byte

	switch t.Language {
	case TemplateLanguageSkel:
		resultBytes, err = t.renderPreprocessorTemplate(contents, options)
	case TemplateLanguageGoText:
		resultBytes, err = t.renderGoTextTemplate(contents, options)
	case TemplateLanguageGoHtml:
		resultBytes, err = t.renderGoHtmlTemplate(contents, options)
	default:
		err = errors.Errorf("Unknown template language %q", t.Language)
	}

	if err != nil {
		return "", err
	}

	result = string(resultBytes)
	for _, transformer := range t.Transformers {
		result = transformer(result)
	}

	return result, nil
}

func (t Template) renderPreprocessorTemplate(contents []byte, options TemplateOptions) (result []byte, err error) {
	// Substitute strings
	if len(options.Strings) > 0 {
		contents = t.substituteStrings(contents, options.Strings)
	}

	// Substitute inline variables
	if len(options.Variables) > 0 {
		contents = t.substituteInlineVariables(contents, options.Variables)

		// Substitute block variables
		contents, err = t.processBlockVariables(contents, options.Variables)
		if err != nil {
			return nil, err
		}
	}

	// Execute conditions
	for condition, value := range options.Conditions {
		contents, err = t.processConditionalBlocks(contents, condition, value)
		if err != nil {
			return
		}
	}

	// Substitute inline identifiers
	contents, err = t.processUserIdentifiers(contents)
	if err != nil {
		return
	}

	// Remove ignored sections
	contents, err = t.processIgnoredBlocks(contents)

	// Combine joined sections
	contents, err = t.processJoinedBlocks(contents)

	return contents, nil
}

func (t Template) renderGoTextTemplate(contents []byte, options TemplateOptions) ([]byte, error) {
	goTemplate := gotext.New(t.Name)
	_, err := goTemplate.Parse(string(contents))
	if err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)
	err = goTemplate.Execute(out, options)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (t Template) renderGoHtmlTemplate(contents []byte, options TemplateOptions) ([]byte, error) {
	goTemplate := gohtml.New(t.Name)
	_, err := goTemplate.Parse(string(contents))
	if err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)
	err = goTemplate.Execute(out, options)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (t Template) substituteStrings(contents []byte, subs map[string]string) []byte {
	// Replace strings
	for sourceString, destString := range subs {
		contents = bytes.ReplaceAll(contents, []byte(sourceString), []byte(destString))
	}
	return contents
}

func (t Template) substituteInlineVariables(source []byte, variableValues map[string]string) []byte {
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

func (t Template) directiveMarkers() Markers {
	markers := t.Format.CommentMarkers()
	markers.Prefix += "#"
	return markers
}

func (t Template) processBlockVariables(data []byte, variables map[string]string) (result []byte, err error) {
	return t.processVisitor(data, &VariableLineVisitor{
		Variables: variables,
		Directive: t.directiveMarkers(),
	})
}

// processConditionalBlocks handles `#if`/`#else`/`#endif` blocks, conditionally outputting lines of code.
// Two passes occur, one for the positive condition, one for the negated condition.
func (t Template) processConditionalBlocks(data []byte, condition string, output bool) (result []byte, err error) {
	// Process direct condition
	data, err = t.processVisitor(data, &ConditionalLineVisitor{
		Name:      condition,
		Value:     output,
		Directive: t.directiveMarkers(),
	})
	if err != nil {
		return nil, err
	}

	// Process negated condition
	return t.processVisitor(data, &ConditionalLineVisitor{
		Name:      "!" + condition,
		Value:     !output,
		Directive: t.directiveMarkers(),
	})
}

// processUserIdentifiers handles `#id` directives, replacing names with template-defined values
func (t Template) processUserIdentifiers(data []byte) (result []byte, err error) {
	return t.processVisitor(data, &IdentifierLineVisitor{
		Directive: t.directiveMarkers(),
	})
}

func (t Template) processIgnoredBlocks(contents []byte) ([]byte, error) {
	return t.processVisitor(contents, &IgnoreLineVisitor{
		Directive: t.directiveMarkers(),
	})
}

func (t Template) processJoinedBlocks(contents []byte) ([]byte, error) {
	return t.processVisitor(contents, &JoinLineVisitor{
		Directive: t.directiveMarkers(),
	})
}

func (t Template) processVisitor(data []byte, visitor TemplateLineVisitor) (result []byte, err error) {
	sb := bytes.Buffer{}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var line string
		var output bool

		line, output, err = visitor.VisitLine(scanner.Text())
		if err != nil {
			return nil, err
		}

		if output {
			sb.WriteString(line)
			sb.WriteRune('\n')
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "Failed to process conditional blocks")
	}

	return sb.Bytes(), nil
}

func NewTemplate(name string, format FileFormat, language TemplateLanguage, source TemplateLoader, transformers ...Transformer) Template {
	var result = Template{
		Name:         name,
		Loader:       source,
		Format:       format,
		Language:     language,
		Transformers: transformers,
	}

	return result
}

type TemplateLineVisitor interface {
	VisitLine(line string) (string, bool, error)
}

type ConditionalLineVisitor struct {
	Name      string  // name of condition
	Value     bool    // value of condition
	Directive Markers // preprocessor marker prefix and suffix

	state   int      // insideIf, insideElse, outside
	markers []string // #if, #else, #endif
}

func (v *ConditionalLineVisitor) VisitLine(line string) (string, bool, error) {
	const outside = 0
	const insideIf = 1
	const insideElse = 2

	if len(v.markers) == 0 {
		v.markers = []string{
			v.Directive.Wrap("if " + v.Name),
			v.Directive.Wrap("else " + v.Name),
			v.Directive.Wrap("endif " + v.Name),
		}
	}

	startMarker, middleMarker, endMarker := v.markers[0], v.markers[1], v.markers[2]

	lineTrimmed := strings.TrimSpace(line)
	switch v.state {
	case outside:
		switch lineTrimmed {
		case startMarker:
			v.state = insideIf
		default:
			return line, true, nil
		}

	case insideIf:
		switch lineTrimmed {
		case endMarker:
			v.state = outside
		case middleMarker:
			v.state = insideElse
		default:
			return line, v.Value, nil
		}

	case insideElse:
		switch lineTrimmed {
		case endMarker:
			v.state = outside
		default:
			return line, !v.Value, nil
		}
	}

	return line, false, nil
}

type IdentifierLineVisitor struct {
	Definitions map[string]string // map of user-defined strings
	Directive   Markers           // preprocessor marker prefix and suffix
	marker      string            // #str
	identifiers map[string]*regexp.Regexp
}

var identifierRegexp = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]+)\s+(.+)\s*$`)

func (v *IdentifierLineVisitor) VisitLine(line string) (string, bool, error) {
	if v.marker == "" {
		v.marker = v.Directive.Prefixed("id ")
	}

	lineTrimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(lineTrimmed, v.marker) {
		// Replace identifiers
		for sourceIdentifier, destString := range v.Definitions {
			re := v.identifiers[sourceIdentifier]
			line = re.ReplaceAllLiteralString(line, destString)
		}
		return line, true, nil
	}

	lineTrimmed = strings.TrimPrefix(lineTrimmed, v.marker)
	lineTrimmed = strings.TrimSuffix(lineTrimmed, v.Directive.Suffix)

	groups := identifierRegexp.FindStringSubmatch(lineTrimmed)
	if len(groups) == 0 {
		return line, false, errors.Errorf("Invalid identifier definition: %s", lineTrimmed)
	}

	name, value := groups[1], groups[2]

	if v.Definitions == nil {
		v.Definitions = make(map[string]string)
	}
	v.Definitions[name] = value

	if v.identifiers == nil {
		v.identifiers = make(map[string]*regexp.Regexp)
	}
	v.identifiers[name] = regexp.MustCompile(`\b` + name + `\b`)

	return line, false, nil
}

type VariableLineVisitor struct {
	Variables map[string]string
	Directive Markers
	marker    string
}

type variableCommand struct {
	Name    string
	Options map[string]string
}

func (v *VariableLineVisitor) VisitLine(line string) (string, bool, error) {
	if v.marker == "" {
		v.marker = v.Directive.Prefixed("var ")
	}

	lineTrimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(lineTrimmed, v.marker) {
		return line, true, nil
	}

	suffix := lineTrimmed[len(v.marker):]
	varCommand, err := v.parseCommand(suffix)
	if err != nil {
		return line, false, err
	}

	varValue, ok := v.Variables[varCommand.Name]
	if ok {
		lineSuffix, ok := varCommand.Options["suffix"]
		if ok {
			// Insert the suffix before any trailing whitespace
			varValueTrimmed := strings.TrimRight(varValue, " \n\t\r")
			varValueTrimmedSuffix := varValue[len(varValueTrimmed):]
			varValue = varValueTrimmed + lineSuffix + varValueTrimmedSuffix
		}

		return varValue, true, nil
	}

	return line, true, nil
}

var varNameRegexp = regexp.MustCompile(`^\s*([[:alpha:]_][\w\\.]*\w)\b\s*`)
var optRegexp = regexp.MustCompile(`^\s*(\w+)=(\S+)[ \t]*`)

func (v *VariableLineVisitor) parseCommand(suffix string) (c variableCommand, err error) {
	nameGroups := varNameRegexp.FindStringSubmatch(suffix)
	if nameGroups == nil {
		err = errors.Errorf("Malformed #var preprocessor command: %s", suffix)
		return
	}

	c.Name = strings.ToLower(nameGroups[1])
	suffix = strings.TrimSpace(suffix[len(nameGroups[0]):])

	for len(suffix) > 0 {
		optGroups := optRegexp.FindStringSubmatch(suffix)
		if optGroups == nil {
			err = errors.Errorf("Malformed #var preprocessor command: %s", suffix)
			return
		}

		optName, optValue := optGroups[1], optGroups[2]
		if c.Options == nil {
			c.Options = make(map[string]string)
		}
		c.Options[optName] = optValue

		suffix = strings.TrimSpace(suffix[len(optGroups[0]):])
	}

	return
}

type JoinLineVisitor struct {
	Directive Markers
	markers   []string
	lines     []string
	state     int
}

func (v *JoinLineVisitor) VisitLine(line string) (string, bool, error) {
	const outside = 0
	const inside = 1

	if len(v.markers) == 0 {
		v.markers = []string{
			v.Directive.Wrap("join"),
			v.Directive.Wrap("endjoin"),
		}
	}

	startMarker, endMarker := v.markers[0], v.markers[1]

	lineTrimmed := strings.TrimSpace(line)
	switch v.state {
	case outside:
		if lineTrimmed == startMarker {
			v.lines = nil
			v.state = inside
			return line, false, nil
		} else {
			return line, true, nil
		}
	case inside:
		if lineTrimmed == endMarker {
			v.state = outside
			return strings.Join(v.lines, " "), true, nil
		} else {
			v.lines = append(v.lines, line)
			return line, false, nil
		}
	default:
		panic("invalid join state")
	}
}

type IgnoreLineVisitor struct {
	Directive Markers
	markers   []string
	state     int
}

func (v *IgnoreLineVisitor) VisitLine(line string) (string, bool, error) {
	const outside = 0
	const inside = 1

	if len(v.markers) == 0 {
		v.markers = []string{
			v.Directive.Wrap("ignore"),
			v.Directive.Wrap("endignore"),
		}
	}

	startMarker, endMarker := v.markers[0], v.markers[1]

	lineTrimmed := strings.TrimSpace(line)
	switch v.state {
	case outside:
		if lineTrimmed == startMarker {
			v.state = inside
			return line, false, nil
		} else {
			return line, true, nil
		}
	case inside:
		if lineTrimmed == endMarker {
			v.state = outside
		}
		return line, false, nil
	default:
		panic("invalid ignore state")
	}
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

func TemplateStringOption(content string, transformers ...Transformer) TemplateLoader {
	return func() ([]byte, error) {
		for _, transformer := range transformers {
			content = transformer(content)
		}
		return []byte(content), nil
	}
}

func TemplateBytesOption(source []byte) TemplateLoader {
	return func() ([]byte, error) {
		return source, nil
	}
}

func TrimNewlineSuffix(source string) (result string) {
	return strings.TrimSuffix(source, "\n")
}
