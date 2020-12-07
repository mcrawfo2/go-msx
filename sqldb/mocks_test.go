//go:generate mockery --name CrudRepositoryApi

package sqldb

import (
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb/mocks"
	"testing"
)

func TestImplementations(t *testing.T) {
	var _ CrudRepositoryApi = &mocks.CrudRepositoryApi{}
}
