package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const configRootFileSystem = "fs"
const configModeDetect = "detect"
const configModeRelease = "release"

var logger = log.NewLogger("msx.fs")
var fsConfig *FileSystemConfig

type FileSystemConfig struct {
	Root      string `config:"default=/"`
	Resources string `config:"default=/var/lib/${spring.application.name}"`
	Configs   string `config:"default=/etc/${spring.application.name}"`
	Binaries  string `config:"default=/usr/bin"`
	Sources   string `config:"default="`
	Mode      string `config:"default=detect"`
}

func NewFileSystemConfig(cfg *config.Config) (*FileSystemConfig, error) {
	var fsConfig FileSystemConfig
	if err := cfg.Populate(&fsConfig, configRootFileSystem); err != nil {
		return nil, err
	}
	return &fsConfig, nil
}

func ConfigureFileSystem(cfg *config.Config) (err error) {
	fsConfig, err = NewFileSystemConfig(cfg)
	if err != nil {
		return err
	}

	if fsConfig.Mode == configModeDetect {
		if fsConfig.Sources == "" {
			fsConfig.Sources, err = types.FindSourceDirFromStack()
			if err == types.ErrSourceDirUnavailable {
				logger.WithError(err).Warningf("Did not detect source directory.")
			} else if err != nil {
				return err
			}
		}
	}

	return nil
}

func Config() *FileSystemConfig {
	if fsConfig == nil {
		panic("FileSystemConfig not created")
	}
	return fsConfig
}

func Sources() string {
	if fsConfig == nil {
		panic("FileSystemConfig not created")
	}
	return fsConfig.Sources
}

func Resources() string {
	if fsConfig == nil {
		panic("FileSystemConfig not created")
	}
	return fsConfig.Resources
}

func Binaries() string {
	if fsConfig == nil {
		panic("FileSystemConfig not created")
	}
	return fsConfig.Binaries
}

func Root() string {
	if fsConfig == nil {
		panic("FileSystemConfig not created")
	}
	return fsConfig.Root
}

func Mode() string {
	if fsConfig == nil {
		panic("FileSystemConfig not created")
	}
	return fsConfig.Mode
}
