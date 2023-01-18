// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package lru

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/idempotency/cache"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func setupMocks() (ContextCache, context.Context, string, cache.CachedWebData) {
	mockClock := types.NewMockClock()
	cconfig := CacheConfig{}
	heap := NewCache2(cconfig.Ttl, cconfig.ExpireLimit, cconfig.ExpireFrequency, cconfig.DeAgeOnAccess, mockClock, false, "")
	heap.expired = make(chan struct{})

	var cwc ContextCache
	cwc = ContextCacheAdapter{Lru: heap}

	ctx := context.Background()

	mockKeyId := "somekeyid"
	mockCreq := cache.CachedRequest{
		Method:     "somemethod",
		RequestURI: "somerequesturi",
	}
	mockCresp := cache.CachedResponse{
		StatusCode: 999,
		Data:       []byte("some contents"),
		Header:     http.Header{},
	}
	mockCData := cache.CachedWebData{
		Req:  mockCreq,
		Resp: mockCresp,
	}
	return cwc, ctx, mockKeyId, mockCData
}

func TestContextCache_SetGet(t *testing.T) {
	cwc, ctx, mockKeyId, mockCData := setupMocks()

	err := cwc.Set(ctx, mockKeyId, mockCData)
	assert.NoError(t, err)

	val, exists, err := cwc.Get(ctx, mockKeyId)
	assert.NoError(t, err)
	assert.Equal(t, exists, true)
	assert.Equal(t, mockCData, val)
}

func TestContextCache_Clear(t *testing.T) {
	cwc, ctx, mockKeyId, mockCData := setupMocks()

	err := cwc.Set(ctx, mockKeyId, mockCData)
	assert.NoError(t, err)

	err = cwc.Clear(ctx)
	assert.NoError(t, err)

	val, exists, err := cwc.Get(ctx, mockKeyId)
	assert.NoError(t, err)
	assert.Equal(t, exists, false)
	assert.Nil(t, val)
}
