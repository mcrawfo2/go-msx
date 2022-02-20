package skel

import (
	"fmt"
	"github.com/gedex/inflector"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
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
	AddTarget("generate-topic-publisher", "Generate publisher topic implementation", GenerateTopicPublisher)
}

func GenerateSystemDomain(args []string) error {
	if len(args) == 0 {
		return errors.New("No Domain Name specified.  Please provide singular domain name.  Examples: 'employee' or 'device connection'")
	}

	domainName := strings.Join(args, " ")
	conditionals := map[string]bool{
		"TENANT_DOMAIN": false,
	}

	return generateDomain(domainName, conditionals)
}

func GenerateTenantDomain(args []string) error {
	if len(args) == 0 {
		return errors.New("No Domain Name specified.  Please provide singular domain name.  Examples: 'employee' or 'device connection'")
	}

	domainName := strings.Join(args, " ")
	conditionals := map[string]bool{
		"TENANT_DOMAIN": true,
	}

	return generateDomain(domainName, conditionals)
}

func GenerateTopicPublisher(args []string) error {
	if len(args) == 0 {
		return errors.New("No Topic Name specified.  Please provide singular topic name.  Examples: 'employee' or 'device connection'")
	}
	topicName := strings.Join(args, " ")
	inflections := inflect(topicName)

	topicPackageName := inflections[inflectionLowerPlural]
	topicPackageSource := path.Join("code", "topic", "stream", "lowerplural")
	topicPackagePath := path.Join("internal", "stream", topicPackageName)
	apiPackagePath := path.Join("pkg", "api")
	apiPackageSource := path.Join("code", "topic", "api")
	apiPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, apiPackagePath)

	templates := TemplateSet{
		{
			Name:       inflections[inflectionTitleSingular] + " Publisher",
			SourceFile: path.Join(topicPackageSource, "publisher.go"),
			DestFile:   fmt.Sprintf(path.Join(topicPackagePath, "publisher_%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Producer",
			SourceFile: path.Join(topicPackageSource, "producer.go"),
			DestFile:   fmt.Sprintf(path.Join(topicPackagePath, "producer_%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Context",
			SourceFile: path.Join(topicPackageSource, "context.go"),
			DestFile:   path.Join(topicPackagePath, "context.go"),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " DTOs",
			SourceFile: path.Join(apiPackageSource, inflectionLowerPlural+"_message.go"),
			DestFile:   fmt.Sprintf(path.Join(apiPackagePath, "%s_message.go"), inflections[inflectionLowerSingular]),
		},
	}

	options := NewRenderOptions()
	options.AddStrings(inflections)
	options.AddString("cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/topic/api", apiPackageUrl)

	if err := templates.Render(options); err != nil {
		return err
	}

	return initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", "stream", topicPackageName))
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

	templates := TemplateSet{
		{
			Name:       inflections[inflectionTitleSingular] + " Log",
			SourceFile: path.Join(domainPackageSource, "log.go"),
			DestFile:   path.Join(domainPackagePath, "log.go"),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Context",
			SourceFile: path.Join(domainPackageSource, "context.go"),
			DestFile:   path.Join(domainPackagePath, "context.go"),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Controller",
			SourceFile: path.Join(domainPackageSource, "controller.go"),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "controller_%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Converter",
			SourceFile: path.Join(domainPackageSource, "converter.go"),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "converter_%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Service",
			SourceFile: path.Join(domainPackageSource, "service.go"),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "service_%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Model",
			SourceFile: path.Join(domainPackageSource, "model.go"),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "model_%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Repository",
			SourceFile: fmt.Sprintf(path.Join(domainPackageSource, "repository_%s.go"), skeletonConfig.Repository),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "repository_%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " DTOs",
			SourceFile: path.Join(apiPackageSource, inflectionLowerPlural+".go"),
			DestFile:   fmt.Sprintf(path.Join(apiPackagePath, "%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Migration",
			SourceFile: path.Join(migratePackageSource, "table.cql"),
			DestFile: fmt.Sprintf(
				path.Join(migratePackagePath, "%s__CREATE_TABLE_%s.%s"),
				migratePrefix,
				inflections[inflectionScreamingSnakeSingular],
				queryFileExtension),
			Format: FileFormatSql,
		},
	}

	options := NewRenderOptions()
	options.AddStrings(inflections)
	options.AddString("cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/domain/api", apiPackageUrl)
	options.AddConditions(conditions)

	err = templates.Render(options)
	if err != nil {
		return err
	}

	return initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", domainPackageName))
}

func nextMigrationPrefix(folder string) (string, error) {
	prefix := "V" + strings.ReplaceAll(skeletonConfig.AppVersion, ".", "_")
	for i := 0; i < 128; i++ {
		glob := fmt.Sprintf("%s_%d__*.%s",
			prefix,
			i,
			skeletonConfig.RepositoryQueryFileExtension())
		matches, _ := filepath.Glob(filepath.Join(folder, glob))
		if len(matches) == 0 {
			return prefix + "_" + strconv.Itoa(i), nil
		}
	}

	return "", errors.Errorf("More than 128 migrations found for %q", prefix)
}
