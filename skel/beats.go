package skel

import (
	"github.com/iancoleman/strcase"
	"path"
	"strings"
)

const (
	inflectionAppTitle           = "App Title"
	inflectionProtocolUpperCamel = "ProtocolUpperCamel"
	inflectionProtocolLowerCamel = "protocolLowerCamel"
)

func init() {
	AddTarget("generate-domain-beats", "Generate beats domain implementation", GenerateBeatsDomain)
}

func GenerateBeatsDomain(args []string) error {
	inflections := map[string]string{
		inflectionAppTitle:           strings.Title(skeletonConfig.AppName),
		inflectionProtocolUpperCamel: strcase.ToCamel(skeletonConfig.BeatProtocol),
		inflectionProtocolLowerCamel: strcase.ToLowerCamel(skeletonConfig.BeatProtocol),
	}

	apiPackageSource := path.Join("code", "beats", "api")
	apiPackagePath := path.Join("pkg", "api")
	apiPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, apiPackagePath)

	beaterPackageSource := path.Join("code", "beats", "beater")
	beaterPackagePath := path.Join("internal", "beater")
	//beaterPackageUrl := path.Join("cto-github.cisco.com/NFV-BU", skeletonConfig.AppName, beaterPackagePath)

	metaPackageSource := path.Join("code", "beats", "_meta")
	metaPackagePath := path.Join("internal", "_meta")

	templates := TemplateSet{
		{
			Name:       inflections[inflectionAppTitle] + " Field Descriptors",
			SourceFile: path.Join(metaPackageSource, "fields.yml"),
			DestFile:   path.Join(metaPackagePath, "fields.yml"),
			Format:     FileFormatYaml,
		},
		{
			Name:       inflections[inflectionAppTitle] + " DTO",
			SourceFile: path.Join(apiPackageSource, "device.go"),
			DestFile:   path.Join(apiPackagePath, "device.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Beater Init",
			SourceFile: path.Join(beaterPackageSource, "init.go.tpl"),
			DestFile:   path.Join(beaterPackagePath, "init.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Beater Config",
			SourceFile: path.Join(beaterPackageSource, "config.go.tpl"),
			DestFile:   path.Join(beaterPackagePath, "config.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Beater State",
			SourceFile: path.Join(beaterPackageSource, "state.go.tpl"),
			DestFile:   path.Join(beaterPackagePath, "state.go"),
		},
		{
			Name:       inflections[inflectionAppTitle] + " Beater Implementation",
			SourceFile: path.Join(beaterPackageSource, "beater.go.tpl"),
			DestFile:   path.Join(beaterPackagePath, "beater.go"),
		},
	}

	options := NewRenderOptions()
	options.AddStrings(inflections)
	options.AddString("cto-github.cisco.com/NFV-BU/go-msx/skel/templates/code/beats/api", apiPackageUrl)

	if err := templates.Render(options); err != nil {
		return err
	}

	return initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd", "app", "main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal", "beater"))
}
