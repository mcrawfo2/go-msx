// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"context"
	"github.com/pkg/errors"
)

type configContextKey int

const (
	contextKeyConfig configContextKey = iota
	contextKeyLatestVersion
)

func ContextWithConfig(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, contextKeyConfig, cfg)
}

func FromContext(ctx context.Context) *Config {
	configInterface := ctx.Value(contextKeyConfig)
	if configInterface != nil {
		if cfg, ok := configInterface.(*Config); ok {
			return cfg
		} else {
			logger.WithContext(ctx).Warnf("Context config value wrong type: %v", cfg)
			return nil
		}
	}

	logger.WithContext(ctx).Warn("Context config not found")
	return nil
}

func MustFromContext(ctx context.Context) *Config {
	cfg := FromContext(ctx)
	if cfg == nil {
		panic(errors.New("Config missing from context"))
	}
	return cfg
}

// Inject during start of operation for consistency
func ContextWithLatestValues(ctx context.Context) context.Context {
	cfg := FromContext(ctx)
	if cfg == nil {
		return ctx
	}

	return context.WithValue(ctx, contextKeyLatestVersion, cfg.values.latest)
}

// Retrieve during lifetime of operation
func LatestValuesFromContext(ctx context.Context) SnapshotValues {
	if raw := ctx.Value(contextKeyLatestVersion); raw != nil {
		if result, ok := raw.(SnapshotValues); ok {
			return result
		}

		logger.Warnf("Context config snapshot wrong type: %v", raw)
	}

	cfg := FromContext(ctx)
	if cfg == nil {
		return emptySnapshotValues
	}

	return cfg.LatestValues()
}
