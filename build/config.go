package build

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/pflagprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"path"
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

	// bootstrap.yml
	configRootAppInfo = "info.app"
	configRootServer  = "server"

	// Output directories
	configOutputPath       = "dist"
	configOutputRootPath   = configOutputPath + "/root"
	configOutputConfigPath = configOutputRootPath + "/etc"
	configOutputBinaryPath = configOutputRootPath + "/usr/bin"
)

type AppInfo struct {
	Name       string `config:"default="`
	Attributes struct {
		DisplayName string `config:"default="`
	}
}

func (p AppInfo) OutputConfigPath() string {
	return path.Join(configOutputConfigPath, p.Name)
}

type Server struct {
	Port int `config:"default=0"`
}

func (p Server) PortString() string {
	return strconv.Itoa(p.Port)
}

type Executable struct {
	Cmd         string // refers to `cmd/<name>/main.go`
	ConfigFiles []string
}

type MsxParams struct {
	Release  string
	Platform struct {
		ParentArtifacts []string
		Version         string
		IncludeGroups   string `config:"default=com.cisco.**"`
	}
}

type Build struct {
	Number string `config:"default=SNAPSHOT"`
	Group  string `config:"default=com.cisco.msx"`
}

type Config struct {
	Timestamp  time.Time
	Msx        MsxParams
	Executable Executable
	Build      Build
	App        AppInfo
	Server     Server
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

func (p Config) InputCommandRoot() string {
	return path.Join("cmd", p.Executable.Cmd)
}

func (p Config) Port() string {
	return strconv.Itoa(p.Server.Port)
}

var BuildConfig = new(Config)

func LoadBuildConfig(ctx context.Context, configFiles []string) (err error) {
	var providers []config.Provider
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

	BuildConfig.Cfg = cfg

	return nil
}
