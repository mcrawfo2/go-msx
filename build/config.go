package build

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/pflagprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"time"
)

var logger = log.NewLogger("build")

const (
	// build.yml
	configRootMsx        = "msx"
	configRootExecutable = "executable"
	configRootBuild      = "build"
	configRootDocker     = "docker"
	configRootKubernetes = "kubernetes"
	configRootManifest   = "manifest"
	configRootGo         = "go"

	// bootstrap.yml
	configRootAppInfo = "info.app"
	configRootServer  = "server"

	// Output directories
	configOutputPath = "dist"

	configOutputRootPath   = configOutputPath + "/root"
	configOutputConfigPath = configOutputRootPath + "/etc"
	configOutputBinaryPath = configOutputRootPath + "/usr/bin"
	configOutputStaticPath = configOutputRootPath + "/var/lib"

	configTestPath = "test"
)

var (
	defaultConfigs = map[string]string{
		"msx.platform.includegroups":   "com.cisco.**",
		"msx.platform.swaggerartifact": "com.cisco.nfv:nfv-swagger",
		"msx.platform.swaggerwebjar":   "org.webjars:swagger-ui:3.22.2",
		"build.number":                 "SNAPSHOT",
		"build.group":                  "com.cisco.msx",
		"manifest.folder":              "Build-Stable",
		"kubernetes.group":             "platformms",
		"docker.baseimage":             "vms-base-stretch:3.8.0-370",
		"docker.repository":            "dockerhub.cisco.com/vms-platform-dev-docker",
		"docker.username":              "",
		"docker.password":              "",
		"go.env.all.GOPRIVATE":         "cto-github.cisco.com/NFV-BU",
		"go.env.all.GOPROXY":           "https://proxy.golang.org,direct",
		"go.env.linux.GOFLAGS":         `-buildmode=pie -i -ldflags="-extldflags=-Wl,-z,now,-z,relro"`,
		"go.env.darwin.GOFLAGS":        `-i`,
	}
)

type AppInfo struct {
	Name       string
	Attributes struct {
		DisplayName string
	}
}

func (p AppInfo) OutputConfigPath() string {
	return path.Join(configOutputConfigPath, p.Name)
}

func (p AppInfo) OutputStaticPath() string {
	return path.Join(configOutputStaticPath, p.Name)
}

type Server struct {
	Port        int
	ContextPath string
}

func (p Server) PortString() string {
	return strconv.Itoa(p.Port)
}

type Executable struct {
	Cmd         string // refers to `cmd/<name>/main.go`
	ConfigFiles []string
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
	BaseImage  string
	Repository string
	Username   string
	Password   string
}

type Kubernetes struct {
	Group string
}

type Config struct {
	Timestamp  time.Time
	Msx        MsxParams
	Go         Go
	Executable Executable
	Build      Build
	App        AppInfo
	Server     Server
	Docker     Docker
	Kubernetes Kubernetes
	Manifest   Manifest
	Cfg        *config.Config
}

func (p Config) FullBuildNumber() string {
	return fmt.Sprintf("%s-%s", p.Msx.Release, p.Build.Number)
}

func (p Config) OutputBinaryPath() string {
	return configOutputBinaryPath
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

var BuildConfig = new(Config)

func LoadBuildConfig(ctx context.Context, configFiles []string) (err error) {
	var providers = []config.Provider{
		config.NewStatic("defaults", defaultConfigs),
	}

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

	if err = cfg.Populate(&BuildConfig.Msx, configRootMsx); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Executable, configRootExecutable); err != nil {
		return
	}

	if err = cfg.Populate(&BuildConfig.Go, configRootGo); err != nil {
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

	BuildConfig.Cfg = cfg

	return nil
}
