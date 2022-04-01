// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContextWithConfig(t *testing.T) {
	var cfg = new(Config)
	var ctx = ContextWithConfig(context.Background(), cfg)
	assert.NotNil(t, ctx)
	assert.Equal(t, cfg, ctx.Value(contextKeyConfig))
}

func TestContextWithLatestValues(t *testing.T) {
	var cfg = new(Config)
	cfg.values.latest = emptySnapshotValues

	var ctx = ContextWithLatestValues(context.Background())
	assert.Equal(t, context.Background(), ctx)

	ctx = ContextWithConfig(ctx, cfg)

	ctx = ContextWithLatestValues(ctx)
	assert.NotNil(t, ctx)
	assert.Equal(t, emptySnapshotValues, ctx.Value(contextKeyLatestVersion))
}

func TestFromContext(t *testing.T) {
	var cfg = new(Config)
	var ctx = ContextWithConfig(context.Background(), cfg)
	assert.Equal(t, cfg, FromContext(ctx))
}

func TestLatestValuesFromContext(t *testing.T) {
	var cfg = new(Config)
	cfg.values.latest = emptySnapshotValues

	var ctx = ContextWithConfig(context.Background(), cfg)
	assert.NotNil(t, ctx)

	// Test auto-calculated latest values
	assert.Equal(t, emptySnapshotValues, LatestValuesFromContext(ctx))

	// Test injected latest values
	cfg.values.latest = SnapshotValues{
		index:   map[string]int{},
		entries: nil,
	}

	ctx = ContextWithLatestValues(ctx)
	assert.NotNil(t, ctx)
	assert.Equal(t, cfg.values.latest, LatestValuesFromContext(ctx))
}

func TestMustFromContext(t *testing.T) {
	defer func() {
		v := recover()
		assert.Error(t, v.(error))
	}()

	MustFromContext(context.Background())
	t.Failed()
}
