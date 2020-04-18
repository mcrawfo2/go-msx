package build

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app/appconfig"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/pflagprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"fmt"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"time"
)

var logger = log.NewLogger("build")

const (
	// build.yml
	configRootMsx        = "msx"
	configRootLibrary    = "library"
	configRootExecutable = "executable"
	configRootBuild      = "build"
	configRootDocker     = "docker"
	configRootKubernetes = "kubernetes"
	configRootManifest   = "manifest"
	configRootGo         = "go"
	configRootGenerate   = "generate"
	configRootResources  = "resources"

	// bootstrap.yml
	configRootAppInfo = "info.app"
	configRootServer  = "server"

	// Output directories
	configOutputRootPath = "dist/root"

	configTestPath = "test"
)

var (
	defaultConfigs = map[string]string{
		"spring.application.name":      "build",
		"msx.platform.includegroups":   "com.cisco.**",
		"msx.platform.swaggerartifact": "com.cisco.nfv:nfv-swagger",
		"msx.platform.swaggerwebjar":   "org.webjars:swagger-ui:3.23.11",
		"build.number":                 "SNAPSHOT",
		"build.group":                  "com.cisco.msx",
		"manifest.folder":              "Build-Stable",
		"kubernetes.group":             "platformms",
		"docker.dockerfile":            "docker/Dockerfile", // TODO: v1.0.0: switch to default 'build/package/Dockerfile'
		"docker.baseimage":             "vms-base-stretch:3.8.0-370",
		"docker.repository":            "dockerhub.cisco.com/vms-platform-dev-docker",
		"docker.username":              "",
		"docker.password":              "",
		"go.env.all.GOPRIVATE":         "cto-github.cisco.com/NFV-BU",
		"go.env.all.GOPROXY":           "https://proxy.golang.org,direct",
		"go.env.linux.GOFLAGS":         `-buildmode=pie -i -ldflags="-extldflags=-Wl,-z,now,-z,relro" -ldflags=-s -ldflags=-w`,
		"go.env.darwin.GOFLAGS":        `-i`,
		"library.name":                 "",
	}
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

type Go struct {
	Env struct {
		All    map[string]string
		Linux  map[string]string
		Darwin map[string]string
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
		SwaggerArtifact string
		SwaggerWebJar   string
		Version         string
		IncludeGroups   string
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
	Command string `config:"default="`
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
}

type Config struct {
	Timestamp  time.Time
	Library    Library
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

var BuildConfig = new(Config)

func newHttpFileProviders(name string, fs http.FileSystem, files []string) []config.Provider {
	var providers []config.Provider
	for _, file := range files {
		providers = append(providers, config.NewHttpFileProvider(name, fs, file))
	}
	return providers
}

func LoadAppBuildConfig(ctx context.Context, cfg *config.Config, providers []config.Provider) (finalConfig *config.Config, err error) {
	if err = cfg.Populate(&BuildConfig.Msx, configRootMsx); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Executable, configRootExecutable); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Build, configRootBuild); err != nil {
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

func LoadBuildConfig(ctx context.Context, configFiles []string) (err error) {
	var providers = []config.Provider{
		config.NewStatic("defaults", defaultConfigs),
	}

	defaultFiles := appconfig.FindConfigHttpFilesGlob(resource.Defaults, "**/defaults-*")
	defaultFilesProviders := newHttpFileProviders("Defaults", resource.Defaults, defaultFiles)
	providers = append(providers, defaultFilesProviders...)

	for _, configFile := range configFiles {
		fileProvider := config.NewFileProvider("Build", configFile)
		providers = append(providers, fileProvider)
	}

	envProvider := config.NewEnvironment("Environment")
	providers = append(providers, envProvider)

	cliProvider := pflagprovider.NewPflagSource("CommandLine", cli.RootCmd().Flags(), "cli.flag.")
	providers = append(providers, cliProvider)

	cfg := config.NewConfig(providers...)
	if err = cfg.Load(ctx); err != nil {
		return
	}

	BuildConfig.Timestamp = time.Now().UTC()

	if err = cfg.Populate(&BuildConfig.Library, configRootLibrary); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Go, configRootGo); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Generate, configRootGenerate); err != nil {
		return
	}

	if BuildConfig.Library.Name == "" {
		if newCfg, err := LoadAppBuildConfig(ctx, cfg, providers); err != nil {
			return err
		} else {
			cfg = newCfg
		}
	}

	if BuildConfig.Fs, err = fs.NewFileSystemConfig(cfg); err != nil {
		return err
	}

	BuildConfig.Cfg = cfg

	return nil
}
