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
		return health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": "Kafka pool not found in context",
			},
		}
	}

	var healthResult health.CheckResult
	err := trace.Operation(ctx, "kafka.healthCheck", func(ctx context.Context) error {
		return kafkaPool.WithConnection(ctx, func(conn *kafka.Connection) error {
			broker, err := conn.Client().Controller()
			if err != nil {
				healthResult = health.CheckResult{
					Status: health.StatusDown,
					Details: map[string]interface{}{
						"error": err.Error(),
					},
				}
				return nil
			}

			connectedResult, err := broker.Connected()
			if err != nil {
				healthResult = health.CheckResult{
					Status: health.StatusDown,
					Details: map[string]interface{}{
						"error": err.Error(),
					},
				}
			} else if !connectedResult {
				healthResult = health.CheckResult{
					Status: health.StatusDown,
					Details: map[string]interface{}{
						"error": "Kafka not connected",
					},
				}
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
