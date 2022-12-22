// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package schema

import "cto-github.cisco.com/NFV-BU/go-msx/config"

const (
	configKeyAppName        = "info.app.name"
	configKeyAppDescription = "info.app.description"
	configKeyBuildVersion   = "info.build.version"
	configKeyDisplayName    = "info.app.attributes.display-name"
)

type AppInfo struct {
	Name        string
	DisplayName string
	Description string
	Version     string
}

func AppInfoFromConfig(cfg *config.Config) (*AppInfo, error) {
	var appInfo AppInfo
	var err error
	if appInfo.Name, err = cfg.String(configKeyAppName); err != nil {
		return nil, err
	}
	if appInfo.Description, err = cfg.String(configKeyAppDescription); err != nil {
		return nil, err
	}
	if appInfo.Version, err = cfg.String(configKeyBuildVersion); err != nil {
		return nil, err
	}
	if appInfo.DisplayName, err = cfg.String(configKeyDisplayName); err != nil {
		return nil, err
	}
	return &appInfo, nil
}
