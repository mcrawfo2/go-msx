// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapiprovider

import "cto-github.cisco.com/NFV-BU/go-msx/resource"

type StaticFileSpecProvider struct {
	cfg DocumentationResourcesConfig
}

func (p StaticFileSpecProvider) Spec() ([]byte, error) {
	return resource.
		Reference(p.cfg.YamlSpecFile).
		ReadAll()
}
