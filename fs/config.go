package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
)

const configRootFileSystem = "fs"

var logger = log.NewLogger("msx.fs")

type FileSystemConfig struct {
	Root      string
	Resources string
	Configs   string
}

func NewFileSystemConfig(cfg *config.Config) (*FileSystemConfig, error) {
	var fsConfig FileSystemConfig
	if err := cfg.Populate(&fsConfig, configRootFileSystem); err != nil {
		return nil, err
	}
	return &fsConfig, nil
}
