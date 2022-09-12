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
