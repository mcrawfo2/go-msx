// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package lru

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupMocks() (ContextCache, context.Context, string, string) {
	mockClock := types.NewMockClock()
	cconfig := CacheConfig{}
	heap := NewCache2(cconfig.Ttl, cconfig.ExpireLimit, cconfig.ExpireFrequency, cconfig.DeAgeOnAccess, mockClock, false, "")
	heap.expired = make(chan struct{})

	var cwc ContextCache
	cwc = ContextCacheAdapter{Lru: heap}

	ctx := context.Background()

	mockKeyId := "somekeyid"
	mockData := "somevalue"
	return cwc, ctx, mockKeyId, mockData
}

func TestContextCache_SetGet(t *testing.T) {
	cwc, ctx, mockKeyId, mockData := setupMocks()

	err := cwc.Set(ctx, mockKeyId, mockData)
	assert.NoError(t, err)

	val, exists, err := cwc.Get(ctx, mockKeyId)
	assert.NoError(t, err)
	assert.Equal(t, exists, true)
	assert.Equal(t, mockData, val)
}

func TestContextCache_Clear(t *testing.T) {
	cwc, ctx, mockKeyId, mockData := setupMocks()

	err := cwc.Set(ctx, mockKeyId, mockData)
	assert.NoError(t, err)

	err = cwc.Clear(ctx)
	assert.NoError(t, err)

	val, exists, err := cwc.Get(ctx, mockKeyId)
	assert.NoError(t, err)
	assert.Equal(t, exists, false)
	assert.Nil(t, val)
}
