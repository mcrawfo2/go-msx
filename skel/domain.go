package skel

import (
	"bufio"
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"fmt"
	"github.com/gedex/inflector"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	inflectionTitleSingular          = "Title Singular"
	inflectionTitlePlural            = "Title Plural"
	inflectionUpperCamelSingular     = "UpperCamelSingular"
	inflectionUpperCamelPlural       = "UpperCamelPlural"
	inflectionLowerCamelSingular     = "lowerCamelSingular"
	inflectionLowerCamelPlural       = "lowerCamelPlural"
	inflectionLowerSnakeSingular     = "lower_snake_singular"
	inflectionScreamingSnakePlural   = "SCREAMING_SNAKE_PLURAL"
	inflectionScreamingSnakeSingular = "SCREAMING_SNAKE_SINGULAR"
	inflectionLowerSingular          = "lowersingular"
	inflectionLowerPlural            = "lowerplural"
)

func init() {
	AddTarget("generate-domain-system", "Generate system domain implementation", GenerateSystemDomain)
	AddTarget("generate-domain-tenant", "Generate tenant domain implementation", GenerateTenantDomain)
}

func GenerateSystemDomain(args []string) error {
	if len(args) == 0 {
		return errors.New("No Domain Name specified.  Please provide singular domain name.  Examples: 'employee' or 'device connection'")
	}

	domainName := strings.Join(args, " ")
	conditionals := map[string]bool{
		"TENANT_DOMAIN":        false,
		"REPOSITORY_COCKROACH": skeletonConfig.Repository == "cockroach",
		"REPOSITORY_CASSANDRA": skeletonConfig.Repository == "cassandra",
	}

	return generateDomain(domainName, conditionals)
}

func GenerateTenantDomain(args []string) error {
	if len(args) == 0 {
		return errors.New("No Domain Name specified.  Please provide singular domain name.  Examples: 'employee' or 'device connection'")
	}

	domainName := strings.Join(args, " ")
	conditionals := map[string]bool{
		"TENANT_DOMAIN":        true,
		"REPOSITORY_COCKROACH": skeletonConfig.Repository == "cockroach",
		"REPOSITORY_CASSANDRA": skeletonConfig.Repository == "cassandra",
	}

	return generateDomain(domainName, conditionals)
}

func inflect(title string) map[string]string {
	titleSingular := strings.Title(inflector.Singularize(title))
	titlePlural := strings.Title(inflector.Pluralize(titleSingular))
	upperCamelSingular := strcase.ToCamel(titleSingular)
	upperCamelPlural := strcase.ToCamel(titlePlural)
	lowerCamelSingular := strcase.ToLowerCamel(titleSingular)
	lowerCamelPlural := strcase.ToLowerCamel(titlePlural)
	lowerSingular := strings.ToLower(lowerCamelSingular)
	lowerPlural := strings.ToLower(lowerCamelPlural)
	lowerSnakeSingular := strcase.ToSnake(titleSingular)
	screamingSnakeSingular := strcase.ToScreamingSnake(titleSingular)
	screamingSnakePlural := strcase.ToScreamingSnake(titlePlural)

	return map[string]string{
		inflectionTitleSingular:          titleSingular,
		inflectionTitlePlural:            titlePlural,
		inflectionUpperCamelSingular:     upperCamelSingular,
		inflectionUpperCamelPlural:       upperCamelPlural,
		inflectionLowerCamelSingular:     lowerCamelSingular,
		inflectionLowerCamelPlural:       lowerCamelPlural,
		inflectionLowerSingular:          lowerSingular,
		inflectionLowerPlural:            lowerPlural,
		inflectionLowerSnakeSingular:     lowerSnakeSingular,
		inflectionScreamingSnakeSingular: screamingSnakeSingular,
		inflectionScreamingSnakePlural:   screamingSnakePlural,
	}
}

type domainDefinitionFile struct {
	Name     string
	Template Template
}

func generateDomain(name string, conditions map[string]bool) error {
	inflections := inflect(name)

	domainPackageName := inflections[inflectionLowerPlural]
	domainPackageSource := path.Join("code", "domain", inflectionLowerPlural)
	domainPackagePath := path.Join("internal", domainPackageName)
	apiPackagePath := path.Join("pkg", "api")
	apiPackageSource := path.Join("code", "domain", "api")
	apiPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, apiPackagePath)
	migratePackageSource := path.Join("code", "domain", "migrate", "version")
	migratePackagePath := path.Join("internal", "migrate", "V"+strings.ReplaceAll(skeletonConfig.AppVersion, ".", "_"))
	migratePrefix, err := nextMigrationPrefix(migratePackagePath)
	if err != nil {
		return err
	}

	queryFileExtension := skeletonConfig.RepositoryQueryFileExtension()

	files := []domainDefinitionFile{
		{
			Name: inflections[inflectionTitleSingular] + " Log",
			Template: Template{
				SourceFile: path.Join(domainPackageSource, "log.go"),
				DestFile:   path.Join(domainPackagePath, "log.go"),
			},
		},
		{
			Name: inflections[inflectionTitleSingular] + " Context",
			Template: Template{
				SourceFile: path.Join(domainPackageSource, "context.go"),
				DestFile:   path.Join(domainPackagePath, "context.go"),
			},
		},
		{
			Name: inflections[inflectionTitleSingular] + " Controller",
			Template: Template{
				SourceFile: path.Join(domainPackageSource, "controller.go"),
				DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "controller_%s.go"), inflections[inflectionLowerSingular]),
			},
		},
		{
			Name: inflections[inflectionTitleSingular] + " Converter",
			Template: Template{
				SourceFile: path.Join(domainPackageSource, "converter.go"),
				DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "converter_%s.go"), inflections[inflectionLowerSingular]),
			},
		},
		{
			Name: inflections[inflectionTitleSingular] + " Service",
			Template: Template{
				SourceFile: path.Join(domainPackageSource, "service.go"),
				DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "service_%s.go"), inflections[inflectionLowerSingular]),
			},
		},
		{
			Name: inflections[inflectionTitleSingular] + " Model",
			Template: Template{
				SourceFile: path.Join(domainPackageSource, "model.go"),
				DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "model_%s.go"), inflections[inflectionLowerSingular]),
			},
		},
		{
			Name: inflections[inflectionTitleSingular] + " Repository",
			Template: Template{
				SourceFile: fmt.Sprintf(path.Join(domainPackageSource, "repository_%s.go"), skeletonConfig.Repository),
				DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "repository_%s.go"), inflections[inflectionLowerSingular]),
			},
		},
		{
			Name: inflections[inflectionTitleSingular] + " DTOs",
			Template: Template{
				SourceFile: path.Join(apiPackageSource, inflectionLowerPlural+".go"),
				DestFile:   fmt.Sprintf(path.Join(apiPackagePath, "%s.go"), inflections[inflectionLowerSingular]),
			},
		},
		{
			Name: inflections[inflectionTitleSingular] + " Migration",
			Template: Template{
				SourceFile: path.Join(migratePackageSource, "table.cql"),
				DestFile: fmt.Sprintf(
					path.Join(migratePackagePath, "%s__CREATE_TABLE_%s.%s"),
					migratePrefix,
					inflections[inflectionScreamingSnakeSingular],
					queryFileExtension),
			},
		},
	}

	packagePaths := map[string]string{
		"cto-github.cisco.com/NFV-BU/go-msx/skel/templates/code/domain/api": apiPackageUrl,
	}

	err = renderDomain(files, inflections, conditions, packagePaths)
	if err != nil {
		return err
	}

	return initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", domainPackageName))
}

func renderDomain(files []domainDefinitionFile, inflections map[string]string, conditions map[string]bool, packagePaths map[string]string) error {
	variableValues := variables()

	for _, file := range files {
		// Load the target
		sourceFile := substituteVariables(file.Template.SourceFile, variableValues)
		bytes, err := readTemplate(sourceFile)
		if err != nil {
			return err
		}
		fileData := string(bytes)

		// Substitute inflections
		for k, v := range inflections {
			fileData = strings.ReplaceAll(fileData, k, v)
		}

		// Substitute API package path
		for sourcePath, destPath := range packagePaths {
			fileData = strings.ReplaceAll(fileData, sourcePath, destPath)
		}

		// Substitute generator variables
		fileData = substituteVariables(fileData, variableValues)

		// Process conditional blocks
		for condition, output := range conditions {
			fileData, err = processConditionalBlocks(fileData, condition, output)
			if err != nil {
				return err
			}
		}

		// Write static file
		err = writeStaticFiles(map[string]StaticFile{
			file.Name: {
				Data:     []byte(fileData),
				DestFile: file.Template.DestFile,
			},
		})
		if err != nil {
			return err
		}

		destFile := path.Join(skeletonConfig.TargetDirectory(), file.Template.DestFile)
		if path.Ext(destFile) == ".go" {
			err = exec.ExecutePipes(
				exec.Info("Reformatting %q", path.Base(destFile)),
				exec.ExecSimple("go", "fmt", destFile))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func nextMigrationPrefix(folder string) (string, error) {
	prefix := "V" + strings.ReplaceAll(skeletonConfig.AppVersion, ".", "_")
	for i := 0; i < 128; i++ {
		glob := fmt.Sprintf("%s/%s_%d__*.%s",
			folder,
			prefix,
			i,
			skeletonConfig.RepositoryQueryFileExtension())
		matches, _ := filepath.Glob(glob)
		if len(matches) == 0 {
			return prefix + "_" + strconv.Itoa(i), nil
		}
	}

	return "", errors.Errorf("More than 128 migrations found for %q", prefix)
}

func processConditionalBlocks(data, condition string, output bool) (result string, err error) {
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
	startMarker := "//#if " + condition
	middleMarker := "//#else " + condition
	endMarker := "//#endif " + condition

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

		switch d.(type) {
		case *ast.GenDecl:
			dd := d.(*ast.GenDecl)

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

	return ioutil.WriteFile(fileName, buffer.Bytes(), 0644)
}
