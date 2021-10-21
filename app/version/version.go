package version

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
)

var logger = log.NewLogger("msx.app.version")

type info struct {
	App struct {
		Name string
	}
	Build struct {
		Version       string
		BuildNumber   string
		BuildDateTime string
		CommitHash    string `config:"default="`
		DiffHash      string `config:"default="`
	}
}

func Version(ctx context.Context, args []string) (err error) {
	var i info
	if err = config.FromContext(ctx).Populate(&i, "info"); err != nil {
		return err
	}

	logger.WithContext(ctx).Infof("App Name: %s", i.App.Name)
	logger.WithContext(ctx).Infof("Build Version: %s", i.Build.Version)
	logger.WithContext(ctx).Infof("Build Timestamp: %s", i.Build.BuildDateTime)
	logger.WithContext(ctx).Infof("Commit Hash: %s", i.Build.CommitHash)
	logger.WithContext(ctx).Infof("Diff Hash: %s", i.Build.DiffHash)
	return nil
}
