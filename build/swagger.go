// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import "cto-github.cisco.com/NFV-BU/go-msx/build/npm"

func init() {
	AddTarget("install-swagger-ui", "Installs Swagger-UI package", InstallSwaggerUi)
}

func InstallSwaggerUi(args []string) error {
	return npm.InstallNodePackageContents(
		BuildConfig.Msx.Platform.Swagger.Artifact,
		BuildConfig.Msx.Platform.Swagger.Version,
		"package/dist",
		BuildConfig.OutputStaticPath())
}
