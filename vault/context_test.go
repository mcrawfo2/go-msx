package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContextWithPool(t *testing.T) {
	ctx := context.Background()
	o := PoolFromContext(ctx)
	assert.Nil(t, o)
	pool = new(ConnectionPool)
	ctx = ContextWithPool(ctx)
	p := PoolFromContext(ctx)
	assert.Equal(t, pool, p)
}

func TestContextWithConnection(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	c := ConnectionFromContext(ctx)
	assert.Nil(t, c)
	conn := new(Connection)
	ctx = ContextWithConnection(ctx, conn)
	conn2 := ConnectionFromContext(ctx)
	assert.Equal(t, conn, conn2)
	ctx = configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.cloud.vault.enabled": "false",
	})
	c = ConnectionFromContext(ctx)
	assert.Nil(t, c)
}
