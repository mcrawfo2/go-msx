package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
)

const configRootFileSystem = "fs"
const configModeDetect = "detect"
const configModeRelease = "release"

var logger = log.NewLogger("msx.fs")

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
