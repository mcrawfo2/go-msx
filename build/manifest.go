// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import "cto-github.cisco.com/NFV-BU/go-msx/exec"

func init() {
	AddTarget("publish-installer-manifest", "Deploy the installer manifests", PublishInstallerManifest)
}

func PublishInstallerManifest(args []string) error {
	return exec.ExecutePipes(
		exec.WithDir("manifest",
			exec.ExecSimple("mvn",
				"-B", "clean", "deploy",
				"-Dversion="+BuildConfig.Msx.Release,
				"-Dbuild_number="+BuildConfig.Build.Number,
				"-Dfolder="+BuildConfig.Manifest.Folder)))
}
