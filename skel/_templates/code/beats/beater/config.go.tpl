package beater

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"time"
)

const (
	configRoot = "${app.name}"
)

type BeatConfig struct {
	Period  time.Duration `config:"default=60s"`
	Timeout time.Duration `config:"default=5s"`
	Batch   BatchConfig
}

type BatchConfig struct {
	Size  int           `config:"default=1"`
	Delay time.Duration `config:"default=6ms"`
}

func newConfig(cfg *config.Config) (*BeatConfig, error) {
	var beatConfig BeatConfig
	if err := cfg.Populate(&beatConfig, configRoot); err != nil {
		return nil, err
	}
	return &beatConfig, nil
}
