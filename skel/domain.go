// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"fmt"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	cmd := AddTarget("generate-domain-system", "Generate system domain implementation", GenerateSystemDomain)
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Use = "generate-domain-system <domain name>"
	cmd.Aliases = []string{"sysdom"}

	cmd = AddTarget("generate-domain-tenant", "Generate tenant domain implementation", GenerateTenantDomain)
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Use = "generate-domain-tenant <domain name>"
	cmd.Aliases = []string{"tendom"}

	cmd = AddTarget("generate-topic-publisher", "Generate publisher topic implementation", GenerateTopicPublisher)
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Use = "generate-topic-publisher <topic name>"
	cmd.Aliases = []string{"toppub"}

	cmd = AddTarget("generate-topic-subscriber", "Generate subscriber topic implementation", GenerateTopicSubscriber)
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Use = "generate-topic-subscriber <topic name>"
	cmd.Aliases = []string{"topsub"}

	cmd = AddTarget("generate-timer", "Generate timer implementation", GenerateTimer)
	cmd.Aliases = []string{"timer"}

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
	inflections := text.NewInflector(topicName)

	topicPackageName := inflections[text.InflectionLowerPlural]
	topicPackageSource := path.Join("code", "topic", "stream", "lowerplural")
	topicPackagePath := path.Join("internal", "stream", topicPackageName)
	apiPackagePath := path.Join("pkg", "api")
	apiPackageSource := path.Join("code", "topic", "api")
	apiPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, apiPackagePath)

	templates := TemplateSet{
		{
			Name:       inflections[text.InflectionTitleSingular] + " Publisher",
			SourceFile: path.Join(topicPackageSource, "publisher.go"),
			DestFile:   fmt.Sprintf(path.Join(topicPackagePath, "publisher_%s.go"), inflections[text.InflectionLowerSingular]),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Package",
			SourceFile: path.Join(topicPackageSource, "pkg.go"),
			DestFile:   path.Join(topicPackagePath, "pkg.go"),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " DTOs",
			SourceFile: path.Join(apiPackageSource, text.InflectionLowerPlural+"_message.go"),
			DestFile:   fmt.Sprintf(path.Join(apiPackagePath, "%s_message.go"), inflections[text.InflectionLowerSingular]),
		},
	}

	options := NewRenderOptions()
	options.AddStrings(inflections)
	options.AddString("cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/topic/api", apiPackageUrl)

	if err := templates.Render(options); err != nil {
		return err
	}

	if err := InitializePackageFromFile(
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
	inflections := text.NewInflector(topicName)

	topicPackageName := inflections[text.InflectionLowerPlural]
	topicPackageSource := path.Join("code", "topic", "stream", "lowerplural")
	topicPackagePath := path.Join("internal", "stream", topicPackageName)
	apiPackagePath := path.Join("pkg", "api")
	apiPackageSource := path.Join("code", "topic", "api")
	apiPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, apiPackagePath)

	templates := TemplateSet{
		{
			Name:       inflections[text.InflectionTitleSingular] + " Subscriber",
			SourceFile: path.Join(topicPackageSource, "subscriber.go"),
			DestFile:   fmt.Sprintf(path.Join(topicPackagePath, "subscriber_%s.go"), inflections[text.InflectionLowerSingular]),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Package",
			SourceFile: path.Join(topicPackageSource, "pkg.go"),
			DestFile:   path.Join(topicPackagePath, "pkg.go"),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " DTOs",
			SourceFile: path.Join(apiPackageSource, text.InflectionLowerPlural+"_message.go"),
			DestFile:   fmt.Sprintf(path.Join(apiPackagePath, "%s_message.go"), inflections[text.InflectionLowerSingular]),
		},
	}

	options := NewRenderOptions()
	options.AddStrings(inflections)
	options.AddString("cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/topic/api", apiPackageUrl)

	if err := templates.Render(options); err != nil {
		return err
	}

	if err := InitializePackageFromFile(
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
	inflections := text.NewInflector(timerName)

	timerPackageName := inflections[text.InflectionLowerPlural]
	timerPackageSource := path.Join("code", "timer", "lowerplural")
	timerPackagePath := path.Join("internal", "timer", timerPackageName)

	templates := TemplateSet{
		{
			Name:       inflections[text.InflectionTitleSingular] + " Timer",
			SourceFile: path.Join(timerPackageSource, "timer.go"),
			DestFile:   fmt.Sprintf(path.Join(timerPackagePath, "timer_%s.go"), inflections[text.InflectionLowerSingular]),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Package",
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
	fixedIntervalKey := "scheduled.tasks." + inflections[text.InflectionLowerSingular] + ".fixed-interval"
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

	if err := InitializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", "timer", timerPackageName)); err != nil {
		return err
	}

	var deps = []string{"github.com/stretchr/testify@v1.8.1"}
	if err := AddDependencies(deps); err != nil {
		return err
	}

	timerPackageAbsPath := filepath.Join(skeletonConfig.TargetDirectory(), timerPackagePath)
	return GoGenerate(timerPackageAbsPath)
}

func generateDomain(name string, conditions map[string]bool) error {
	inflections := text.NewInflector(name)

	domainPackageName := inflections[text.InflectionLowerPlural]
	domainPackageSource := path.Join("code", "domain", text.InflectionLowerPlural)
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
			Name:       inflections[text.InflectionTitleSingular] + " Log",
			SourceFile: path.Join(domainPackageSource, "log.go"),
			DestFile:   path.Join(domainPackagePath, "log.go"),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Context",
			SourceFile: path.Join(domainPackageSource, "context.go"),
			DestFile:   path.Join(domainPackagePath, "context.go"),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Controller",
			SourceFile: path.Join(domainPackageSource, "controller.go"),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "controller_%s.go"), inflections[text.InflectionLowerSingular]),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Converter",
			SourceFile: path.Join(domainPackageSource, "converter.go"),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "converter_%s.go"), inflections[text.InflectionLowerSingular]),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Service",
			SourceFile: path.Join(domainPackageSource, "service.go"),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "service_%s.go"), inflections[text.InflectionLowerSingular]),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Model",
			SourceFile: path.Join(domainPackageSource, "model.go"),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "model_%s.go"), inflections[text.InflectionLowerSingular]),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Repository",
			SourceFile: fmt.Sprintf(path.Join(domainPackageSource, "repository_%s.go"), skeletonConfig.Repository),
			DestFile:   fmt.Sprintf(path.Join(domainPackagePath, "repository_%s.go"), inflections[text.InflectionLowerSingular]),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " DTOs",
			SourceFile: path.Join(apiPackageSource, text.InflectionLowerPlural+".go"),
			DestFile:   fmt.Sprintf(path.Join(apiPackagePath, "%s.go"), inflections[text.InflectionLowerSingular]),
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Migration",
			SourceFile: path.Join(migratePackageSource, "table.cql"),
			DestFile: fmt.Sprintf(
				path.Join(migratePackagePath, "%s__CREATE_TABLE_%s.%s"),
				migratePrefix,
				inflections[text.InflectionScreamingSnakeSingular],
				queryFileExtension),
			Format: text.FileFormatSql,
		},
		{
			Name:       inflections[text.InflectionTitleSingular] + " Permissions",
			SourceFile: path.Join(populatePackage, "manifest.json"),
			DestFile:   path.Join(populatePackage, "manifest.json"),
			Format:     text.FileFormatJson,
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

	return InitializePackageFromFile(
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
