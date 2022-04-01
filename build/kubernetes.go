// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

func init() {
	AddTarget("install-kubernetes-manifests", "Install the distribution kubernetes manifests", InstallKubernetesManifests)
}

func InstallKubernetesManifests(args []string) error {
	return nil
}
