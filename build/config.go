// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/pflagprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var logger = log.NewLogger("build")

const (
	// build.yml
	configRootMsx         = "msx"
	configRootLibrary     = "library"
	configRootTool        = "tool"
	configRootExecutable  = "executable"
	configRootBuild       = "build"
	configRootDocker      = "docker"
	configRootKubernetes  = "kubernetes"
	configRootManifest    = "manifest"
	configRootGo          = "go"
	configRootGenerate    = "generate"
	configRootResources   = "resources"
	configRootArtifactory = "artifactory"
	configRootLicense     = "license"
	configRootSonatype    = "ossi"

	// bootstrap.yml
	configRootAppInfo = "info.app"
	configRootServer  = "server"

	// Output directories
	configOutputRootPath = "dist/root"
	configOutputToolPath = "dist/tools"

	configDistPath = "dist"
	configTestPath = "test"
)

var (
	defaultConfigs = map[string]string{
		"spring.application.name":        "build",
		"msx.platform.includegroups":     "com.cisco.**",
		"msx.platform.swagger.artifact":  "@msx/nfv-swagger-ui",
		"msx.platform.swagger.version":   "1.1.0",
		"msx.platform.asyncapi.artifact": "@msx/asyncapi-studio",
		"msx.platform.asyncapi.version":  "0.12.11",
		"msx.release":                    "5.0.0",
		"msx.platform.version":           "5.0.0-STABLE",
		"build.number":                   "SNAPSHOT",
		"build.group":                    "com.cisco.msx",
		"manifest.folder":                "Build-Stable",
		"kubernetes.group":               "platformms",
		"docker.dockerfile":              "build/package/Dockerfile",
		"docker.baseimage":               "msx-base-buster:3.9.0-70",
		"docker.repository":              "dockerhub.cisco.com/vms-platform-dev-docker",
		"docker.username":                "",
		"docker.password":                "",
		"go.env.all.GOPRIVATE":           "cto-github.cisco.com/NFV-BU",
		"go.env.all.GOPROXY":             "https://proxy.golang.org,direct",
		"go.env.linux.GOFLAGS":           `-buildmode=pie -ldflags=-extldflags=-Wl,-z,now,-z,relro -ldflags=-s -ldflags=-w`,
		"go.env.darwin.GOFLAGS":          ``,
		"library.name":                   "",
		"tool.cmd":                       "",
		"tool.name":                      "",
		"artifactory.repository":         "",
		"artifactory.username":           "",
		"artifactory.password":           "",
		"ossi.username":                  "",
		"ossi.token":                     "",
	}
	defaultConfigProvider = config.NewInMemoryProvider("defaults", defaultConfigs)
	defaultConfigCache    = config.NewCacheProvider(defaultConfigProvider)
)

type AppInfo struct {
	Name       string
	Attributes struct {
		DisplayName string
	}
}

type Server struct {
	Port        int
	ContextPath string
	StaticPath  string
}

func (p Server) PortString() string {
	return strconv.Itoa(p.Port)
}

type Executable struct {
	Cmd         string // refers to `cmd/<name>/main.go`
	ConfigFiles []string
}

type Library struct {
	Name string
}

type Tool struct {
	Cmd       string
	Name      string
	Resources Resources
}

func (t Tool) PublishUrl(goos string) string {
	return strings.Join([]string{
		BuildConfig.Binaries.Repository,
		t.Name,
		BuildConfig.Build.Number,
		t.PublishArtifactName(goos)},
		"/")
}

func (t Tool) PublishLatestUrl(goos string) string {
	return strings.Join([]string{
		BuildConfig.Binaries.Repository,
		t.Name,
		"latest",
		t.PublishArtifactLatestName(goos)},
		"/")
}

func (t Tool) PackageFolder(goos string) string {
	return filepath.Join(BuildConfig.OutputToolPath(), t.Name, goos)
}

func (t Tool) PublishArtifactName(goos string) string {
	return fmt.Sprintf("%s-%s-%s.tar.gz",
		t.Name,
		goos,
		BuildConfig.Build.Number)
}

func (t Tool) PublishArtifactLatestName(goos string) string {
	return fmt.Sprintf("%s-%s-%s.tar.gz",
		t.Name,
		goos,
		"latest")
}

func (t Tool) PublishArtifactPath(goos string) string {
	return filepath.Join(BuildConfig.OutputToolPath(), t.Name, t.PublishArtifactName(goos))
}

type Go struct {
	Env struct {
		All    map[string]string
		Linux  map[string]string
		Darwin map[string]string
	}
	Vet struct {
		Options []string
	}
}

func (g Go) Environment() map[string]string {
	result := make(map[string]string)
	copyMap := func(source map[string]string) {
		for k, v := range source {
			result[k] = v
		}
	}
	copyMap(g.Env.All)
	switch runtime.GOOS {
	case "linux":
		copyMap(g.Env.Linux)
	case "darwin":
		copyMap(g.Env.Darwin)
	}
	return result
}

type MsxParams struct {
	Release  string
	Platform struct {
		ParentArtifacts []string
		Version         string
		IncludeGroups   string
		Swagger         struct {
			Artifact string
			Version  string
		}
		AsyncApi struct {
			Artifact string
			Version  string
		}
	}
}

type Build struct {
	Number string
	Group  string
}

type Manifest struct {
	Folder string
}

type Docker struct {
	Dockerfile string
	BaseImage  string
	Repository string
	Username   string
	Password   string
}

type Kubernetes struct {
	Group string
}

type Generate struct {
	Path    string
	Command string   `config:"default="`
	BuiltIn []string `config:"default="`
	VfsGen  *GenerateVfs
}

type GenerateVfs struct {
	Root         string `config:"default="`
	Filename     string `config:"default=assets.go"`
	VariableName string `config:"default=assets"`
	Includes     []string
	Excludes     []string `config:"default="`
}

type Resources struct {
	Includes []string
	Excludes []string
	Mappings []PathMapping
}

type PathMapping struct {
	From string
	To   string
}

type Binaries struct {
	Repository string // Root URL of artifactory binaries repository
	Username   string // Injected from Jenkins credentials store via ARTIFACTORY_USERNAME
	Password   string // Injected from Jenkins credentials store via ARTIFACTORY_PASSWORD
}

func (b Binaries) Authorization() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(b.Username+":"+b.Password))
}

type Module struct {
	GoModPath    string
	ModulePath   string
	MinGoVersion string
}

type License struct {
	Excludes []string
}

type Sonatype struct {
	Username string
	Token    string
}

type Config struct {
	Timestamp  time.Time
	Library    Library
	Tool       Tool
	Msx        MsxParams
	Go         Go
	Executable Executable
	Build      Build
	App        AppInfo
	Server     Server
	Docker     Docker
	Kubernetes Kubernetes
	Manifest   Manifest
	Generate   []Generate
	Resources  Resources
	Binaries   Binaries
	Module     Module
	License    License
	Ossi       Sonatype
	Fs         *fs.FileSystemConfig
	Cfg        *config.Config
}

func (p Config) FullBuildNumber() string {
	return fmt.Sprintf("%s-%s", p.Msx.Release, p.Build.Number)
}

func (p Config) OutputRoot() string {
	return configOutputRootPath
}

func (p Config) TestPath() string {
	return configTestPath
}

func (p Config) DistPath() string {
	return configDistPath
}

func (p Config) InputCommandRoot() string {
	return path.Join("cmd", p.Executable.Cmd)
}

func (p Config) Port() string {
	return strconv.Itoa(p.Server.Port)
}

func (p Config) OutputConfigPath() string {
	return path.Join(configOutputRootPath, p.Fs.Root, p.Fs.Configs)
}

func (p Config) OutputResourcesPath() string {
	return path.Join(configOutputRootPath, p.Fs.Root, p.Fs.Resources)
}

func (p Config) OutputBinaryPath() string {
	return path.Join(configOutputRootPath, p.Fs.Root, p.Fs.Binaries)
}

func (p Config) OutputStaticPath() string {
	return path.Join(p.OutputResourcesPath(), "www")
}

func (p Config) OutputToolPath() string {
	return configOutputToolPath
}

var BuildConfig = new(Config)

func LoadGoModule(m *Module) error {
	sourceDir, err := types.FindSourceDirFromStack()
	if err != nil {
		return err
	}

	goModPath := path.Join(sourceDir, "go.mod")
	goModBytes, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return err
	}

	goModFile, err := modfile.Parse(goModPath, goModBytes, nil)
	if err != nil {
		return err
	}

	m.GoModPath = goModPath
	m.MinGoVersion = goModFile.Go.Version
	m.ModulePath = goModFile.Module.Mod.Path
	return nil
}

func LoadAppBuildConfig(ctx context.Context, cfg *config.Config, providers []config.Provider) (finalConfig *config.Config, err error) {
	if err = cfg.Populate(&BuildConfig.Msx, configRootMsx); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Executable, configRootExecutable); err != nil {
		return
	}

	for _, v := range BuildConfig.Executable.ConfigFiles {
		filePath := path.Join(BuildConfig.InputCommandRoot(), v)
		fileProvider := config.NewFileProvider(v, filePath)
		providers = append(providers, fileProvider)
	}

	cfg = config.NewConfig(providers...)
	if err = cfg.Load(ctx); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.App, configRootAppInfo); err != nil {
		return
	}

	// Set the spring app name if it is not set
	springAppName, _ := cfg.StringOr("spring.application.name", "build")
	if springAppName == "build" {
		defaultConfigs["spring.application.name"] = BuildConfig.App.Name
		defaultConfigCache.Invalidate()
		_ = cfg.Load(ctx)
	}

	if err = cfg.Populate(&BuildConfig.Server, configRootServer); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Docker, configRootDocker); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Kubernetes, configRootKubernetes); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Manifest, configRootManifest); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Resources, configRootResources); err != nil {
		return
	}

	return cfg, nil
}

func LoadBuildConfig(ctx context.Context, cmd *cobra.Command, configFiles []string) (err error) {
	var providers = []config.Provider{
		defaultConfigCache,
	}

	defaultFilesProviders := config.EmbeddedDefaultsProviders
	providers = append(providers, defaultFilesProviders...)

	for _, configFile := range configFiles {
		fileProvider := config.NewFileProvider("Build", configFile)
		providers = append(providers, fileProvider)
	}

	envProvider := config.NewEnvironmentProvider("Environment")
	providers = append(providers, envProvider)

	cliProvider := pflagprovider.NewProvider("CommandLine", cmd.Flags(), "cli.flag")
	providers = append(providers, cliProvider)

	cfg := config.NewConfig(providers...)
	if err = cfg.Load(ctx); err != nil {
		return
	}

	BuildConfig.Timestamp = time.Now().UTC()

	if err = cfg.Populate(&BuildConfig.Library, configRootLibrary); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Tool, configRootTool); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Go, configRootGo); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Generate, configRootGenerate); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Binaries, configRootArtifactory); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Build, configRootBuild); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.License, configRootLicense); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Ossi, configRootSonatype); err != nil {
		return
	}

	if BuildConfig.Library.Name == "" && BuildConfig.Tool.Name == "" {
		if newCfg, err := LoadAppBuildConfig(ctx, cfg, providers); err != nil {
			return err
		} else {
			cfg = newCfg
		}
	}

	if BuildConfig.Fs, err = fs.NewFileSystemConfig(cfg); err != nil {
		return err
	}

	if err = LoadGoModule(&BuildConfig.Module); err != nil {
		return err
	}

	BuildConfig.Cfg = cfg

	return nil
}
