package build

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

func init() {
	AddTarget("generate-build-info", "Create a build metadata file", GenerateBuildInfo)
}

const (
	configDefaultBuildInfoFile = "buildinfo.yml"
)

type MsxBuildInfo struct {
	Version       string `yaml:"version"`
	BuildNumber   string `yaml:"buildNumber"`
	BuildDateTime string `yaml:"buildDateTime"`
	Artifact      string `yaml:"artifact"`
	Name          string `yaml:"name"`
	Group         string `yaml:"group"`
}

func (m MsxBuildInfo) ToMap() map[string]string {
	return map[string]string{
		"version":       m.Version,
		"buildNumber":   m.BuildNumber,
		"buildDateTime": m.BuildDateTime,
		"artifact":      m.Artifact,
		"name":          m.Name,
		"group":         m.Group,
	}
}

func wrap(name string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		name: value,
	}
}

func (m MsxBuildInfo) MarshalYAML() (interface{}, error) {
	return wrap("info.build", m.ToMap()), nil
}

func GenerateBuildInfo(args []string) error {
	buildInfo := new(MsxBuildInfo)
	buildInfo.Version = BuildConfig.FullBuildNumber()
	buildInfo.BuildNumber = BuildConfig.Build.Number
	buildInfo.BuildDateTime = BuildConfig.Timestamp.Format("2006-01-02T15:04:05.999999Z")
	buildInfo.Artifact = BuildConfig.App.Name
	buildInfo.Name = BuildConfig.App.Attributes.DisplayName
	buildInfo.Group = BuildConfig.Build.Group

	buildInfoBytes, err := yaml.Marshal(buildInfo)
	if err != nil {
		return err
	}

	// Tack on source root for relative lookups
	sourceDir, err := os.Getwd()
	if err != nil {
		return err
	}
	fsSources := "\n" + `fs.sources: ` + sourceDir + "\n"
	buildInfoBytes = append(buildInfoBytes, []byte(fsSources)...)

	buildInfoFile := path.Join(BuildConfig.OutputConfigPath(), configDefaultBuildInfoFile)

	err = os.MkdirAll(BuildConfig.OutputConfigPath(), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(buildInfoFile, buildInfoBytes, 0644)
}
