package skel

import (
	"path"
	"strings"
)

func init() {
	AddTarget("generate-service-pack", "Generate service pack implementation", GenerateServicePack)
}

func GenerateServicePack(args []string) error {
	inflections := map[string]string{
		inflectionAppTitle: strings.Title(skeletonConfig.AppName),
	}

	apiPackageSource := path.Join("code", "sp", "api")
	apiPackagePath := path.Join("pkg", "api")
	apiPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, apiPackagePath)

	subscriptionPackageSource := path.Join("code", "sp", "subscription")
	subscriptionPackagePath := path.Join("internal", "subscription")

	slmPackageSource := path.Join("code", "sp", "platform-common", "servicelifecycle")
	slmPackagePath := path.Join("platform-common", "servicelifecycle")

	files := []domainDefinitionFile{
		{
			Name: inflections[inflectionAppTitle] + " Subscription DTO",
			Template: Template{
				SourceFile: path.Join(apiPackageSource, "subscription.go"),
				DestFile:   path.Join(apiPackagePath, "subscription.go"),
			},
		},
		{
			Name: inflections[inflectionAppTitle] + " Subscription Context",
			Template: Template{
				SourceFile: path.Join(subscriptionPackageSource, "context.go.tpl"),
				DestFile:   path.Join(subscriptionPackagePath, "context.go"),
			},
		},
		{
			Name: inflections[inflectionAppTitle] + " Subscription Controller",
			Template: Template{
				SourceFile: path.Join(subscriptionPackageSource, "controller.go.tpl"),
				DestFile:   path.Join(subscriptionPackagePath, "controller.go"),
			},
		},
		{
			Name: inflections[inflectionAppTitle] + " Subscription Controller",
			Template: Template{
				SourceFile: path.Join(subscriptionPackageSource, "converter.go.tpl"),
				DestFile:   path.Join(subscriptionPackagePath, "converter.go"),
			},
		},
		{
			Name: inflections[inflectionAppTitle] + " Subscription Model",
			Template: Template{
				SourceFile: path.Join(subscriptionPackageSource, "model.go.tpl"),
				DestFile:   path.Join(subscriptionPackagePath, "model.go"),
			},
		},
		{
			Name: inflections[inflectionAppTitle] + " Subscription Service",
			Template: Template{
				SourceFile: path.Join(subscriptionPackageSource, "service.go.tpl"),
				DestFile:   path.Join(subscriptionPackagePath, "service.go"),
			},
		},
		{
			Name: inflections[inflectionAppTitle] + " Service Lifecycle Manifest",
			Template: Template{
				SourceFile: path.Join(slmPackageSource, "manifest.json"),
				DestFile:   path.Join(slmPackagePath, "manifest.json"),
			},
		},
		{
			Name: inflections[inflectionAppTitle] + " Service Lifecycle Deployment Manifest",
			Template: Template{
				SourceFile: path.Join(slmPackageSource, "manifest.yml"),
				DestFile:   path.Join(slmPackagePath, "manifest.yml"),
			},
		},
	}

	packagePaths := map[string]string{
		"cto-github.cisco.com/NFV-BU/go-msx/skel/templates/code/sp/api": apiPackageUrl,
	}

	err := renderDomain(files, inflections, nil, packagePaths)
	if err != nil {
		return err
	}

	return initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", "subscription"))
}
