// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	AddTarget("generate-domain-system", "Generate system domain implementation", GenerateSystemDomain)
	AddTarget("generate-domain-tenant", "Generate tenant domain implementation", GenerateTenantDomain)
	AddTarget("generate-topic-publisher", "Generate publisher topic implementation", GenerateTopicPublisher)
	AddTarget("generate-topic-subscriber", "Generate subscriber topic implementation", GenerateTopicSubscriber)
	AddTarget("generate-timer", "Generate timer implementation", GenerateTimer)
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
	inflections := NewInflector(topicName)

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
			Name:       inflections[inflectionTitleSingular] + " Package",
			SourceFile: path.Join(topicPackageSource, "pkg.go"),
			DestFile:   path.Join(topicPackagePath, "pkg.go"),
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

	if err := initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", "stream", topicPackageName)); err != nil {
		return err
	}

	var deps = []string{"github.com/stretchr/testify@v1.7.0"}
	if err := AddDependencies(deps); err != nil {
		return err
	}

	topicPackageAbsPath := filepath.Join(skeletonConfig.TargetDirectory(), topicPackagePath)
	return GoGenerate(topicPackageAbsPath)
}

func GenerateTopicSubscriber(args []string) error {
	if len(args) == 0 {
		return errors.New("No Topic Name specified.  Please provide singular topic name.  Examples: 'employee' or 'device connection'")
	}
	topicName := strings.Join(args, " ")
	inflections := NewInflector(topicName)

	topicPackageName := inflections[inflectionLowerPlural]
	topicPackageSource := path.Join("code", "topic", "stream", "lowerplural")
	topicPackagePath := path.Join("internal", "stream", topicPackageName)
	apiPackagePath := path.Join("pkg", "api")
	apiPackageSource := path.Join("code", "topic", "api")
	apiPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, apiPackagePath)

	templates := TemplateSet{
		{
			Name:       inflections[inflectionTitleSingular] + " Subscriber",
			SourceFile: path.Join(topicPackageSource, "subscriber.go"),
			DestFile:   fmt.Sprintf(path.Join(topicPackagePath, "subscriber_%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Package",
			SourceFile: path.Join(topicPackageSource, "pkg.go"),
			DestFile:   path.Join(topicPackagePath, "pkg.go"),
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

	if err := initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", "stream", topicPackageName)); err != nil {
		return err
	}

	var deps []string
	deps = append(deps,
		"github.com/ThreeDotsLabs/watermill/message",
		"github.com/pkg/errors",
		"github.com/stretchr/testify@v1.7.0")

	if err := AddDependencies(deps); err != nil {
		return err
	}

	topicPackageAbsPath := filepath.Join(skeletonConfig.TargetDirectory(), topicPackagePath)
	return GoGenerate(topicPackageAbsPath)
}

func GenerateTimer(args []string) error {
	if len(args) == 0 {
		return errors.New("No Timer Name specified.  Please provide singular timer name.  Examples: 'employee' or 'device connection'")
	}
	timerName := strings.Join(args, " ")
	inflections := NewInflector(timerName)

	timerPackageName := inflections[inflectionLowerPlural]
	timerPackageSource := path.Join("code", "timer", "lowerplural")
	timerPackagePath := path.Join("internal", "timer", timerPackageName)

	templates := TemplateSet{
		{
			Name:       inflections[inflectionTitleSingular] + " Timer",
			SourceFile: path.Join(timerPackageSource, "timer.go"),
			DestFile:   fmt.Sprintf(path.Join(timerPackagePath, "timer_%s.go"), inflections[inflectionLowerSingular]),
		},
		{
			Name:       inflections[inflectionTitleSingular] + " Package",
			SourceFile: path.Join(timerPackageSource, "pkg.go"),
			DestFile:   path.Join(timerPackagePath, "pkg.go"),
		},
	}

	options := NewRenderOptions()
	options.AddStrings(inflections)

	// Add required configurations to bootstrap configuration file
	bootstrapFilePath := path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "bootstrap.yml")
	leaderElectionKey := "consul.leader.election"
	leaderElectionConfig := leaderElectionKey + ":\n  enabled: true\n"
	leaderElectionRegEx := regexp.MustCompile(`(?m)^` + leaderElectionKey + `:\n  enabled: (.*)\n$`)
	fixedIntervalKey := "scheduled.tasks." + inflections[inflectionLowerSingular] + ".fixed-interval"
	fixedIntervalConfig := fixedIntervalKey + ": 15m\n"
	fixedIntervalRegEx := regexp.MustCompile(`(?m)^` + fixedIntervalKey + `:(.*)\n$`)

	if err := addYamlConf(bootstrapFilePath, leaderElectionKey, leaderElectionConfig, leaderElectionRegEx); err != nil {
		return err
	}

	if err := addYamlConf(bootstrapFilePath, fixedIntervalKey, fixedIntervalConfig, fixedIntervalRegEx); err != nil {
		return err
	}

	if err := templates.Render(options); err != nil {
		return err
	}

	if err := initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", "timer", timerPackageName)); err != nil {
		return err
	}

	var deps = []string{"github.com/stretchr/testify@v1.7.0"}
	if err := AddDependencies(deps); err != nil {
		return err
	}

	timerPackageAbsPath := filepath.Join(skeletonConfig.TargetDirectory(), timerPackagePath)
	return GoGenerate(timerPackageAbsPath)
}

func generateDomain(name string, conditions map[string]bool) error {
	inflections := NewInflector(name)

	domainPackageName := inflections[inflectionLowerPlural]
	domainPackageSource := path.Join("code", "domain", inflectionLowerPlural)
	domainPackagePath := path.Join("internal", domainPackageName)
	apiPackagePath := path.Join("pkg", "api")
	apiPackageSource := path.Join("code", "domain", "api")
	apiPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, apiPackagePath)
	populatePackage := path.Join("internal", "populate", "usermanagement", "permission", "templates")
	migratePackageSource := path.Join("code", "domain", "migrate", "version")
	migratePackagePath := path.Join("internal", "migrate", "V"+strings.ReplaceAll(skeletonConfig.AppVersion, ".", "_"))
	migratePrefix, err := NextMigrationPrefix(migratePackagePath)
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
		{
			Name:       inflections[inflectionTitleSingular] + " Permissions",
			SourceFile: path.Join(populatePackage, "manifest.json"),
			DestFile:   path.Join(populatePackage, "manifest.json"),
			Format:     FileFormatJson,
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

func NextMigrationPrefix(folder string) (string, error) {
	prefix := "V" + strings.ReplaceAll(skeletonConfig.AppVersion, ".", "_")
	for i := 0; i < 100; i++ {
		glob := fmt.Sprintf("%s_%d__*.%s",
			prefix,
			i,
			skeletonConfig.RepositoryQueryFileExtension())
		matches, _ := filepath.Glob(filepath.Join(folder, glob))
		if len(matches) == 0 {
			return prefix + "_" + strconv.Itoa(i), nil
		}
	}

	return "", errors.Errorf("More than 100 migrations found for %q", prefix)
}
