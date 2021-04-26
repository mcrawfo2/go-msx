package mocks

import (
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"
	"testing"
)

func TestImplementations(t *testing.T) {
	var _ sqldb.CrudRepositoryApi = &CrudRepositoryApi{}
	var _ sqldb.CrudRepositoryFactoryApi = &CrudRepositoryFactory{}
}
