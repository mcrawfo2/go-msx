// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
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

var (
	ErrFileExistsAlready = errors.New("file exists already, it should not")
	ErrFileDoesNotExist  = errors.New("file does not exist, it should")
)

var printOptionsDeltas = false // only need to print the NewRenderOptions() values once

type RenderOptions struct {
	text.TemplateOptions
	IncFiles    []string // glob patterns of files to consider (empty=all)
	ExcFiles    []string // glob patterns of files to exclude (empty=none)
	NoOverwrite bool     // RenderTo will not alter existing files, no matter what

	printDeltas  bool // print the changed items only
	DeltaOptions text.TemplateOptions
}

func (r RenderOptions) AddString(source, dest string) {
	r.TemplateOptions.AddString(source, dest)
	r.DeltaOptions.AddString(source, dest)
}

func (r RenderOptions) AddStrings(strings map[string]string) {
	r.TemplateOptions.AddStrings(strings)
	r.DeltaOptions.AddStrings(strings)
}

func (r RenderOptions) AddVariable(source, dest string) {
	r.TemplateOptions.AddVariable(source, dest)
	r.DeltaOptions.AddVariable(source, dest)
}

func (r RenderOptions) AddVariables(variables map[string]string) {
	r.TemplateOptions.AddVariables(variables)
	r.DeltaOptions.AddVariables(variables)
}

func (r RenderOptions) AddCondition(condition string, value bool) {
	r.TemplateOptions.AddCondition(condition, value)
	r.DeltaOptions.AddCondition(condition, value)
}

func (r RenderOptions) AddConditions(conditions map[string]bool) {
	r.TemplateOptions.AddConditions(conditions)
	r.DeltaOptions.AddConditions(conditions)
}

func NewEmptyRenderOptions() RenderOptions {
	return RenderOptions{
		TemplateOptions: text.NewTemplateOptions(),
		IncFiles:        incFiles,
		ExcFiles:        excFiles,
		DeltaOptions:    text.NewTemplateOptions(),
	}
}

func NewTemplateOptions() text.TemplateOptions {
	return text.TemplateOptions{
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
			"repository.cockroach.enabled": strconv.FormatBool(skeletonConfig.Repository == "cockroach"),
			"jenkins.publish.trunk":        strconv.FormatBool(skeletonConfig.KubernetesGroup != "platformms"),
			"generator":                    skeletonConfig.Archetype,
			"beat.protocol":                skeletonConfig.BeatProtocol,
			"service.type":                 skeletonConfig.ServiceType,
			"slack.channel":                skeletonConfig.SlackChannel,
			"trunk":                        skeletonConfig.Trunk,
			"gomsx.version":                GoMsxVersion,
		},
		Conditions: map[string]bool{
			"REPOSITORY_COCKROACH":   skeletonConfig.Repository == "cockroach",
			"GENERATOR_APP":          skeletonConfig.Archetype == archetypeKeyApp,
			"NOT_GENERATOR_BEAT":     skeletonConfig.Archetype != archetypeKeyBeat,
			"GENERATOR_BEAT":         skeletonConfig.Archetype == archetypeKeyBeat,
			"GENERATOR_SP":           skeletonConfig.Archetype == archetypeKeyServicePack,
			"GENERATOR_SPUI":         skeletonConfig.Archetype == archetypeKeySPUI,
			"UI":                     hasUI(),
			"K8S_GROUP_DATAPLATFORM": skeletonConfig.KubernetesGroup == "dataplatform",
			"EXTERNAL":               IsExternal,
		},
		Strings: map[string]string{},
	}
}

func NewRenderOptions() RenderOptions {
	return RenderOptions{
		TemplateOptions: NewTemplateOptions(),
		IncFiles:        incFiles,
		ExcFiles:        excFiles,
		DeltaOptions:    text.NewTemplateOptions(),
	}
}

func (r RenderOptions) String() string {
	if printOptionsDeltas {
		if len(r.DeltaOptions.Strings) > 0 || len(r.DeltaOptions.Variables) > 0 || len(r.DeltaOptions.Conditions) > 0 {
			return fmt.Sprintf("Changed\n%s%s%s\n",
				mapStr("Variables", r.DeltaOptions.Variables),
				mapStr("Conditions", r.DeltaOptions.Conditions),
				mapStr("Strings", r.DeltaOptions.Strings))
		}
		return ""
	}
	printOptionsDeltas = true
	return fmt.Sprintf("%s%s%s%s%s\n",
		mapStr("Variables", r.Variables),
		mapStr("Conditions", r.Conditions),
		mapStr("Strings", r.Strings),
		sliceStr("Include", r.IncFiles),
		sliceStr("Exclude: ", r.ExcFiles))
}

// renders a map as a string with a title but returns "" if the map is empty
func mapStr[t interface{ string | bool }](name string, m map[string]t) string {
	if len(m) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString(name + ":\n")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("  %s = %v\n", k, v))
	}

	return buf.String() + "\n"
}

// renders a slice of strings as a string with a title but returns "" if the slice is empty
func sliceStr(name string, s []string) string {

	if len(s) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString(name + ":")
	for _, v := range s {
		buf.WriteString(fmt.Sprintf("%s  ", v))
	}

	return buf.String() + "\n"
}

type Template struct {
	Name       string
	DestFile   string
	SourceFile string
	SourceData []byte
	Format     text.FileFormat
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

	data, err := ReadStaticFile(sourceFile)
	if err != nil {
		return "", errors.Wrap(err, "Template file not found")
	}

	return string(data), nil
}

func SubstituteVariables(content string, variables map[string]string) string {
	t := text.NewTemplate("vars", text.FileFormatOther, text.TemplateLanguageSkel,
		text.TemplateStringOption(content))

	o := text.NewTemplateOptions()
	o.Variables = variables

	result, _ := t.Render(o)
	return result
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

	includeIt, err := shouldInclude(destFile, options)
	if err != nil {
		return err
	}

	if !includeIt {
		logger.Infof("o (skip on inc/exc) %s (%s)", t.Name, destFile)
		return nil
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

	if options.NoOverwrite && targetExists {
		logger.Infof("  🔒️ (skip) %s (%s) in no-overwrite mode", t.Name, targetFileName)
		return nil
	}

	switch t.Operation {

	case OpGone: // remove if it exists
		if targetExists {
			err = os.Remove(targetFileName)
			if err != nil {
				return errors.Wrapf(err, "removing old target %q failed", targetFileName)
			}
			logger.Infof("  ♻️️ (removed) %s (%s)", t.Name, targetFileName)
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
		logger.Infof("  🗑️ (removed) %s (%s)", t.Name, targetFileName)
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
			logger.Infof("  🔒️ (skip) %s (%s)", t.Name, targetFileName)
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
	if t.Format == text.FileFormatBash { // +x for bash scripts
		perm = os.FileMode(0755)
	}
	err = os.WriteFile(targetFileName, []byte(newcontents), perm)
	if err != nil {
		return errors.Wrapf(err, "writing %q failed", targetFileName)
	}

	if t.Format == text.FileFormatGo {
		logger.Infof("  🎨 Reformatting %q", path.Base(destFile))
		err = exec.ExecutePipes(
			exec.ExecQuiet("go", []string{"fmt", targetFileName}))
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

	tt := text.NewTemplate(t.Name, t.Format, text.TemplateLanguageSkel, text.TemplateStringOption(contents))
	out, err := tt.Render(options.TemplateOptions)
	if err != nil {
		return "", err
	}

	return out, nil
}

type TemplateSet []Template

func (t TemplateSet) Render(options RenderOptions) error {
	return t.RenderTo(skeletonConfig.TargetDirectory(), options)
}

// shouldInclude returns true if the destination file should be considered for output
func shouldInclude(destination string, options RenderOptions) (includeIt bool, err error) {

	includeIt = true
	excludeIt := false

	if len(options.IncFiles) > 0 || len(options.ExcFiles) > 0 {
		logger.Tracef("Considering file %s against (inc/exc) %s %s", destination, options.IncFiles, options.ExcFiles)
	}

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
	optionsStr := options.String()
	if len(optionsStr) > 0 {
		logger.Debugf("Template rendering options: \n%s\n", optionsStr)
	}
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

func InitializePackageFromFile(fileName, packageUrl string) error {
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
