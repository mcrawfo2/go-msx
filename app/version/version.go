// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package version

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"runtime/debug"
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

	bi, ok := debug.ReadBuildInfo()
	if ok {
		logger.WithContext(ctx).Infof("Go Version: %s", bi.GoVersion)
		logger.WithContext(ctx).Info("Build Settings:")
		for _, setting := range bi.Settings {
			logger.WithContext(ctx).Infof("    %s: %s", setting.Key, setting.Value)
		}
		logger.WithContext(ctx).Info("Modules:")
		for _, dep := range bi.Deps {
			if dep.Replace == nil {
				logger.WithContext(ctx).Infof("    %s: %s", dep.Path, dep.Version)
			} else {
				logger.WithContext(ctx).Infof("    %s replaced with %s: %s", dep.Path, dep.Replace.Path, dep.Replace.Version)
			}
		}
	} else {
		logger.WithContext(ctx).Warn("No embedded go BuildInfo found.")
	}
	return nil
}
