// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/gochannel"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/kafka"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/redis"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/sql"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, registerKafkaStreamProvider)
	OnEvent(EventConfigure, PhaseAfter, registerGoChannelStreamProvider)
	OnEvent(EventConfigure, PhaseAfter, registerRedisStreamProvider)
	OnEvent(EventConfigure, PhaseAfter, registerSqlStreamProvider)
	OnEvent(EventStart, PhaseDuring, asyncapi.DocumentStreams)
	OnEvent(EventStart, PhaseAfter, stream.StartRouter)
	OnEvent(EventStop, PhaseBefore, stream.StopRouter)
}

func registerKafkaStreamProvider(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	logger.WithContext(ctx).Info("Registering kafka stream provider")
	if err := kafka.RegisterProvider(cfg); err != nil && err != kafka.ErrDisabled {
		return err
	} else if err == kafka.ErrDisabled {
		logger.WithContext(ctx).WithError(err).Warn("Kafka disabled.  Stream provider not registered.")
	}
	return nil
}

func registerGoChannelStreamProvider(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	logger.WithContext(ctx).Info("Registering gochannel stream provider")
	if err := gochannel.RegisterProvider(cfg); err != nil {
		return err
	}
	return nil
}

func registerRedisStreamProvider(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	logger.WithContext(ctx).Info("Registering redis stream provider")
	if err := redis.RegisterProvider(cfg); err != nil && err != redis.ErrDisabled {
		return err
	} else if err == redis.ErrDisabled {
		logger.WithContext(ctx).WithError(err).Warn("Redis disabled.  Stream provider not registered.")
	}
	return nil
}

func registerSqlStreamProvider(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	logger.WithContext(ctx).Info("Registering sql stream provider")
	if err := sql.RegisterProvider(cfg); err != nil && err != sql.ErrDisabled {
		return err
	} else if err == sql.ErrDisabled {
		logger.WithContext(ctx).WithError(err).Warn("SQL disabled.  Stream provider not registered.")
	}
	return nil
}
