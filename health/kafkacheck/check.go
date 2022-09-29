// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package kafkacheck

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/kafka"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
)

func Check(ctx context.Context) health.CheckResult {
	kafkaPool := kafka.PoolFromContext(ctx)
	if kafkaPool == nil {
		return errorCheckResult(nil, "Kafka pool not found in context")
	}
	connectionConfig := kafkaPool.ConnectionConfig()

	var healthResult health.CheckResult
	err := trace.Operation(ctx, "kafka.healthCheck", func(ctx context.Context) error {
		return kafkaPool.WithConnection(ctx, func(conn *kafka.Connection) error {
			broker, err := conn.Client().Controller()
			if err != nil {
				healthResult = errorCheckResult(connectionConfig, err.Error())
				return nil
			}

			connectedResult, err := broker.Connected()
			if err != nil {
				healthResult = errorCheckResult(connectionConfig, err.Error())
			} else if !connectedResult {
				healthResult = errorCheckResult(connectionConfig, "Kafka not connected")
			} else {
				healthResult = health.CheckResult{
					Status:  health.StatusUp,
					Details: nil,
				}
			}

			return nil
		})
	})

	if err != nil {
		healthResult = health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}
	}

	return healthResult
}

func errorCheckResult(connectionConfig *kafka.ConnectionConfig, s string) health.CheckResult {
	details := map[string]interface{}{
		"error": s,
	}

	if connectionConfig != nil {
		details = map[string]interface{}{
			"error":   s,
			"brokers": connectionConfig.Brokers,
			"zk":      connectionConfig.ZkNodes,
			"tls": map[string]interface{}{
				"enabled":     connectionConfig.Tls.Enabled,
				"certificate": connectionConfig.Tls.CertificateSource,
			},
		}
	}

	return health.CheckResult{
		Status:  health.StatusDown,
		Details: details,
	}
}
