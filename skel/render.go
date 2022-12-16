// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type FileFormat int

const (
	FileFormatGo FileFormat = iota
	FileFormatMakefile
	FileFormatJson
	FileFormatSql
	FileFormatYaml
	FileFormatXml
	FileFormatGroovy
	FileFormatProperties
	FileFormatMarkdown
	FileFormatGoMod
	FileFormatDocker
	FileFormatBash
	FileFormatJavaScript
	FileFormatTypeScript
	FileFormatJenkins
	FileFormatOther
)

var (
	ErrFileExistsAlready = errors.New("file exists already, it should not")
	ErrFileDoesNotExist  = errors.New("file does not exist, it should")
)

type RenderOptions struct {
	Variables  map[string]string
	Conditions map[string]bool
	Strings    map[string]string
	IncFiles   []string // glob patterns of files to consider (empty=all)
	ExcFiles   []string // glob patterns of files to exclude (empty=none)
}

func (r RenderOptions) AddString(source, dest string) {
	r.Strings[source] = dest
}

func (r RenderOptions) AddStrings(strings map[string]string) {
	for k, v := range strings {
		r.Strings[k] = v
	}
}

func (r RenderOptions) AddVariable(source, dest string) {
	r.Variables[source] = dest
}

func (r RenderOptions) AddVariables(variables map[string]string) {
	for k, v := range variables {
		r.Variables[k] = v
	}
}

func (r RenderOptions) AddCondition(condition string, value bool) {
	r.Conditions[condition] = value
}

func (r RenderOptions) AddConditions(conditions map[string]bool) {
	for k, v := range conditions {
		r.Conditions[k] = v
	}
}

func NewRenderOptions() RenderOptions {
	return RenderOptions{
		Variables: map[string]string{
			"app.name":                     skeletonConfig.AppName,
			"app.shortname":                strings.TrimSuffix(skeletonConfig.AppName, "service"),
			"app.uuid":                     skeletonConfig.AppUUID,
			"app.description":              skeletonConfig.AppDescription,
			"app.displayname":              skeletonConfig.AppDisplayName,
			"app.version":                  skeletonConfig.AppVersion,
			"app.migrateversion":           skeletonConfig.AppMigrateVersion(),
			"app.packageurl":               skeletonConfig.AppPackageUrl(),
			"deployment.group":             skeletonConfig.DeploymentGroup,
			"server.port":                  strconv.Itoa(skeletonConfig.ServerPort),
			"debug.port":                   strconv.Itoa(skeletonConfig.DebugPort),
			"server.contextpath":           path.Clean("/" + skeletonConfig.ServerContextPath),
			"server.contextpath.noroot":    strings.TrimPrefix(path.Clean(skeletonConfig.ServerContextPath), "/"),
			"kubernetes.group":             skeletonConfig.KubernetesGroup,
			"target.dir":                   skeletonConfig.TargetDirectory(),
			"repository.cassandra.enabled": strconv.FormatBool(skeletonConfig.Repository == "cassandra"),
			"repository.cockroach.enabled": strconv.FormatBool(skeletonConfig.Repository == "cockroach"),
			"jenkins.publish.trunk":        strconv.FormatBool(skeletonConfig.KubernetesGroup != "platformms"),
			"generator":                    skeletonConfig.Archetype,
			"beat.protocol":                skeletonConfig.BeatProtocol,
			"service.type":                 skeletonConfig.ServiceType,
			"slack.channel":                skeletonConfig.SlackChannel,
			"trunk":                        skeletonConfig.Trunk,
		},
		Conditions: map[string]bool{
			"REPOSITORY_COCKROACH":   skeletonConfig.Repository == "cockroach",
			"REPOSITORY_CASSANDRA":   skeletonConfig.Repository == "cassandra",
			"GENERATOR_APP":          skeletonConfig.Archetype == archetypeKeyApp,
			"NOT_GENERATOR_BEAT":     skeletonConfig.Archetype != archetypeKeyBeat,
			"GENERATOR_BEAT":         skeletonConfig.Archetype == archetypeKeyBeat,
			"GENERATOR_SP":           skeletonConfig.Archetype == archetypeKeyServicePack,
			"GENERATOR_SPUI":         skeletonConfig.Archetype == archetypeKeySPUI,
			"UI":                     hasUI(),
			"K8S_GROUP_DATAPLATFORM": skeletonConfig.KubernetesGroup == "dataplatform",
		},
		Strings:  make(map[string]string),
		IncFiles: incFiles,
		ExcFiles: excFiles,
	}
}

func (r RenderOptions) String() string {
	return fmt.Sprintf("Variables: \n%s\nConditions: \n%s\nStrings: %s\nInclude: %s\nExclude: %s\n",
		mapStr(r.Variables), mapStr(r.Conditions), mapStr(r.Strings),
		r.IncFiles, r.ExcFiles)
}

func mapStr[t interface{ string | bool }](m map[string]t) string {
	var buf bytes.Buffer
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("%s = %v\n", k, v))
	}
	return buf.String()
}

type TemplateOp int

type Template struct {
	Name       string
	DestFile   string
	SourceFile string
	SourceData []byte
	Format     FileFormat
	Operation  TemplateOperation
}

type TemplateOperation int

const (
	OpAdd            TemplateOperation = iota // add or replace
	OpNew                                     // must not exist before being added
	OpAddNoOverwrite                          // add if not exists
	OpReplace                                 // must already exist
	OpDelete                                  // must exist before removal
	OpGone                                    // might not exist before removal
)

func (t Template) source(options RenderOptions) (string, error) {
	if t.SourceData != nil {
		return string(t.SourceData), nil
	}

	sourceFile := SubstituteVariables(t.SourceFile, options.Variables)

	f, ok := staticFiles[sourceFile]
	if !ok {
		return "", errors.Errorf("Template file not found: %s", sourceFile)
	}

	var reader io.Reader
	if f.size != 0 {
		var err error
		reader, err = gzip.NewReader(strings.NewReader(f.data))
		if err != nil {
			return "", err
		}
	} else {
		reader = strings.NewReader(f.data)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (t Template) destinationFile(options RenderOptions) (string, error) {
	if len(t.DestFile) == 0 {
		if len(t.SourceFile) == 0 {
			return "", errors.New("Missing destination filename")
		}
		return SubstituteVariables(t.SourceFile, options.Variables), nil
	}
	return SubstituteVariables(t.DestFile, options.Variables), nil
}

// Render does the original Render to skeletonConfig.TargetDirectory()
func (t Template) Render(options RenderOptions) error {
	return t.RenderTo(skeletonConfig.TargetDirectory(), options)
}

// RenderTo allows rendering to a directory root other than skeletonConfig.TargetDirectory()
func (t Template) RenderTo(directory string, options RenderOptions) error {

	// Find the destination
	destFile, err := t.destinationFile(options)
	if err != nil {
		return err
	}

	targetFileName := path.Join(directory, destFile)
	targetDirectory := path.Dir(targetFileName)

	logger.Tracef("Considering file %s against (inc/exc) %s %s", destFile, options.IncFiles, options.ExcFiles)

	includeIt, err := shouldInclude(destFile, options)
	if err != nil {
		return err
	}

	if !includeIt {
		logger.Infof("o (skip on inc/exc) %s (%s)", t.Name, destFile)
		return nil
	}

	if t.Operation == OpDelete || t.Operation == OpGone {
		if t.Operation == OpDelete { // we *must* delete
			_, err := os.Stat(targetFileName)
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("- %s - erroneously gone - (%s)", t.Name, targetFileName)
			}
		}
		err := os.Remove(targetFileName)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) { // ok if not there already
				logger.Infof("- %s - already gone - (%s)", t.Name, targetFileName)
				return nil
			} else {
				return fmt.Errorf("unable to remove %s: %w", targetFileName, err)
			}
		}
		logger.Infof("- %s (%s)", t.Name, targetFileName)
		return nil
	}

	if t.Operation == OpReplace { // we insist it exists before we replace it
		_, err := os.Stat(targetFileName)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("target to be replaced does not exist %s: %w", targetFileName, err)
		} else if err != nil {
			return fmt.Errorf("error checking target file %s: %w", targetFileName, err)
		}
	}

	targetExists := true
	_, err = os.Stat(targetFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			targetExists = false
		} else {
			return errors.Wrapf(err, "checking for target %q failed", targetFileName)
		}
	}

	switch t.Operation {

	case OpGone: // remove if it exists
		if targetExists {
			err = os.Remove(targetFileName)
			if err != nil {
				return errors.Wrapf(err, "removing old target %q failed", targetFileName)
			}
			logger.Infof("X %s (%s)", t.Name, targetFileName)
		}
		return nil

	case OpDelete: // it must exist and be removed
		if !targetExists {
			return errors.Errorf("file %q required but missing", targetFileName)
		}
		err = os.Remove(targetFileName)
		if err != nil {
			return errors.Wrapf(err, "unable to remove %q", targetFileName)
		}
		logger.Infof("- %s (%s)", t.Name, targetFileName)
		return nil

	case OpAdd: // exists? we care not

	case OpReplace: // we insist it exists before we replace it
		if !targetExists {
			return errors.Wrapf(ErrFileDoesNotExist, "failed replacing %q", targetFileName)
		}

	case OpNew: // we insist it does not exist before we add it
		if targetExists {
			return errors.Wrapf(ErrFileExistsAlready, "failed making new %q", targetFileName)
		}

	case OpAddNoOverwrite: // it may exist or not, but we don't overwrite it
		if targetExists {
			logger.Infof("= (skip) %s (%s)", t.Name, targetFileName)
			return nil
		}
	}

	newcontents, err := t.RenderContents(options)
	if err != nil {
		return err
	}

	// Ensure the target parent directory exists
	err = os.MkdirAll(targetDirectory, 0755)
	if err != nil {
		return err
	}

	// Write the rendered contents to the destination file
	perm := os.FileMode(0644)
	if t.Format == FileFormatBash { // +x for bash scripts
		perm = os.FileMode(0755)
	}
	err = os.WriteFile(targetFileName, []byte(newcontents), perm)
	if err != nil {
		return errors.Wrapf(err, "writing %q failed", targetFileName)
	}

	if t.Format == FileFormatGo {
		err = exec.ExecutePipesStderr(
			exec.Info("  - Reformatting %q", path.Base(destFile)),
			exec.ExecSimple("go", "fmt", targetFileName))
		if err != nil {
			return errors.Wrapf(err, "reformatting %q failed", targetFileName)
		}
	}

	logger.Infof("+ %s (%s)", t.Name, targetFileName)

	return nil
}

func (t Template) RenderContents(options RenderOptions) (result string, err error) {
	// Load the source
	contents, err := t.source(options)
	if err != nil {
		return
	}

	if len(contents) < 3 {
		logger.Warnf("improbably short template `%s` (%s, %d bytes)", t.Name, t.SourceFile, len(contents))
	}

	// Replace strings
	for sourceString, destString := range options.Strings {
		contents = strings.ReplaceAll(contents, sourceString, destString)
	}

	// Substitute variables
	contents = SubstituteVariables(contents, options.Variables)

	// Execute conditions
	for condition, value := range options.Conditions {
		contents, err = processConditionalBlocks(contents, t.Format, condition, value)
	}
	if err != nil {
		return
	}

	result = contents
	return
}

type TemplateSet []Template

func (t TemplateSet) Render(options RenderOptions) error {
	return t.RenderTo(skeletonConfig.TargetDirectory(), options)
}

// shouldInclude returns true if the destination file should be considered for output
func shouldInclude(destination string, options RenderOptions) (includeIt bool, err error) {

	includeIt = true
	excludeIt := false

	if len(options.IncFiles) > 0 {
		includeIt = false
		for _, m := range options.IncFiles {
			j, err := doublestar.Match(m, destination)
			if err != nil {
				return false, err
			}
			includeIt = includeIt || j
		}
		if includeIt {
			logger.Tracef("Including %s", destination)
		} else {
			logger.Tracef("Not including %s", destination)
		}
	}

	if len(options.ExcFiles) > 0 {
		for _, m := range options.ExcFiles {
			j, err := doublestar.Match(m, destination)
			if err != nil {
				return false, err
			}
			excludeIt = excludeIt || j
		}
		if excludeIt {
			logger.Tracef("Excluding %s", destination)
			includeIt = false
		}
	}

	return includeIt && !excludeIt, nil
}

func (t TemplateSet) RenderTo(directory string, options RenderOptions) (err error) {
	logger.Debugf("Template rendering options: \n%s\n", options)
	for _, template := range t {

		if err := template.RenderTo(directory, options); err != nil {
			return err
		}
	}
	return nil
}

func (t TemplateSet) Dirs(opts RenderOptions) (results []string) {
	dirs := types.NewStringSet()
	for _, template := range t {
		destFile := SubstituteVariables(template.DestFile, opts.Variables)
		destDir := path.Dir(destFile)
		dirs.Add(destDir)
	}
	return dirs.Values()
}

func SubstituteVariables(source string, variableValues map[string]string) string {
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

func conditionalMarkers(format FileFormat) (string, string) {
	switch format {
	case FileFormatMakefile, FileFormatYaml, FileFormatProperties, FileFormatDocker, FileFormatBash:
		return "#", ""
	case FileFormatSql:
		return "--#", ""
	case FileFormatXml, FileFormatMarkdown:
		return "<--#", "-->"
	default:
		return "//#", ""
	}
}

func processConditionalBlocks(data string, format FileFormat, condition string, output bool) (result string, err error) {
	type parserState int
	const outside parserState = 0
	const insideIf parserState = 1
	const insideElse parserState = 2

	sb := strings.Builder{}
	write := func(out bool, line string) {
		if !out {
			return
		}
		sb.WriteString(line)
		sb.WriteRune('\n')
	}
	insideCondition := outside
	prefix, suffix := conditionalMarkers(format)
	startMarker := prefix + "if " + condition + suffix
	middleMarker := prefix + "else " + condition + suffix
	endMarker := prefix + "endif " + condition + suffix

	scanner := bufio.NewScanner(strings.NewReader(data))
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

	if err := scanner.Err(); err != nil {
		return "", errors.Wrap(err, "Failed to process conditional blocks")
	}

	return sb.String(), nil
}

func initializePackageFromFile(fileName, packageUrl string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return err
	}

	// Add the imports
	for i := 0; i < len(f.Decls); i++ {
		d := f.Decls[i]

		switch dd := d.(type) {
		case *ast.GenDecl:
			//dd := d.(*ast.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				// Add the new import
				iSpec := &ast.ImportSpec{
					Path: &ast.BasicLit{Value: strconv.Quote(packageUrl)},
					Name: ast.NewIdent("_"),
				}

				dd.Specs = append(dd.Specs, iSpec)
			}
		}
	}

	// Sort the imports
	ast.SortImports(fset, f)

	var output []byte
	buffer := bytes.NewBuffer(output)
	if err := printer.Fprint(buffer, fset, f); err != nil {
		return err
	}

	return os.WriteFile(fileName, buffer.Bytes(), 0644)
}

func hasUI() bool {
	uiPath := filepath.Join(skeletonConfig.TargetDirectory(), "ui", "package.json")
	if st, err := os.Stat(uiPath); err != nil {
		return false
	} else {
		return !st.IsDir() && st.Size() > 0
	}
}

func iff(cond bool, truth, falsehood string) string {
	if cond {
		return truth
	}
	return falsehood
}

// addYamlConf attempts to add conf to the file at filePath.
// If confKey exists in the file, the existing configs that match regEx are replaced
// and the result is written back to the file.
// If the confKey does not exist in the file, we append conf to the end of the file.
func addYamlConf(filePath, confKey, conf string, regEx *regexp.Regexp) error {
	logger.Infof("Adding configuration for %s to %s", confKey, filePath)
	config, err := getYamlConf(filePath, confKey)
	if err != nil {
		return err
	}

	if config == nil {
		if err = appendYaml(filePath, []byte("\n"+conf)); err != nil {
			return err
		}
	} else if err = replaceYaml(filePath, []byte(conf), regEx); err != nil {
		return err
	}

	return nil
}

// getYamlConf retrieves an interface mapped to a given conf
// within the file at the given filePath.
func getYamlConf(filePath, conf string) (interface{}, error) {
	sourceData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err = yaml.Unmarshal(sourceData, &result); err != nil {
		return "", err
	}

	return result[conf], nil
}

// appendYaml appends the yaml to the end of the file at
// the given filePath. The resulting data is written back to the file.
func appendYaml(filePath string, yaml []byte) error {
	sourceData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	sourceData = append(sourceData, yaml...)
	if err = os.WriteFile(filePath, sourceData, 0644); err != nil {
		return err
	}

	return nil
}

// replaceYaml replaces all strings that match regEx in
// the file at filePath with yaml, and logs a warning
// if there are no matches. The resulting data is written back to the file.
func replaceYaml(filePath string, yaml []byte, regEx *regexp.Regexp) error {
	sourceData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if ok := regEx.Match(sourceData); ok {
		sourceData = regEx.ReplaceAll(sourceData, yaml)
	} else {
		logger.Warnf("Failed to add the following configuation to %s:\n%sAs it already exists with a different configuration in %s", filePath, yaml, filePath)
	}

	if err = os.WriteFile(filePath, sourceData, 0644); err != nil {
		return err
	}

	return nil
}
