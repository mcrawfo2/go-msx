// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"path"

	"golang.org/x/text/cases"
)

func init() {
	AddTarget("generate-service-pack", "Generate service pack implementation", GenerateServicePack)
}

func GenerateServicePack(args []string) error {
	caser := cases.Title(TitlingLanguage)
	inflections := map[string]string{
		inflectionAppTitle: caser.String(skeletonConfig.AppName),
	}

	apiPackageSource := path.Join("code", "sp", "api")
	apiPackagePath := path.Join("pkg", "api")
	apiPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, apiPackagePath)

	subscriptionPackageSource := path.Join("code", "sp", "subscription")
	subscriptionPackagePath := path.Join("internal", "subscription")

	slmPackageSource := path.Join("code", "sp", "platform-common", "servicelifecycle")
	slmPackagePath := path.Join("platform-common", "servicelifecycle")

	templates := TemplateSet{
		{
			Name:       inflections[inflectionAppTitle] + " Subscription DTO",
			SourceFile: path.Join(apiPackageSource, "subscription.go"),
			DestFile:   path.Join(apiPackagePath, "subscription.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Subscription Context",
			SourceFile: path.Join(subscriptionPackageSource, "context.go.tpl"),
			DestFile:   path.Join(subscriptionPackagePath, "context.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Subscription Controller",
			SourceFile: path.Join(subscriptionPackageSource, "controller.go.tpl"),
			DestFile:   path.Join(subscriptionPackagePath, "controller.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Subscription Controller",
			SourceFile: path.Join(subscriptionPackageSource, "converter.go.tpl"),
			DestFile:   path.Join(subscriptionPackagePath, "converter.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Subscription Model",
			SourceFile: path.Join(subscriptionPackageSource, "model.go.tpl"),
			DestFile:   path.Join(subscriptionPackagePath, "model.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Subscription Service",
			SourceFile: path.Join(subscriptionPackageSource, "service.go.tpl"),
			DestFile:   path.Join(subscriptionPackagePath, "service.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Service Lifecycle Manifest",
			SourceFile: path.Join(slmPackageSource, "manifest.json"),
			DestFile:   path.Join(slmPackagePath, "manifest.json"),
			Format:     FileFormatJson,
		},
		{
			Name:       inflections[inflectionAppTitle] + " Service Lifecycle Deployment Manifest",
			SourceFile: path.Join(slmPackageSource, "manifest.yml"),
			DestFile:   path.Join(slmPackagePath, "manifest.yml"),
			Format:     FileFormatYaml,
		},
	}

	options := NewRenderOptions()
	options.AddStrings(inflections)
	options.AddString("cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/sp/api", apiPackageUrl)

	if err := templates.Render(options); err != nil {
		return err
	}

	return initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", "subscription"))
}
