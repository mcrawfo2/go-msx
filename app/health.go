package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

const (
	rootKeyHealth = "health"
)

var (
	healthLogger *HealthLogger
)

type HealthLoggerConfig struct {
	Enabled   bool          `config:"default=true"`
	Frequency time.Duration `config:"default=15s"`
}

type HealthLogger struct {
	ctx  context.Context
	cfg  *HealthLoggerConfig
	done chan struct{}
}

func (l *HealthLogger) LogHealth() {
	healthReport := health.GenerateReport(l.ctx)

	if bytes, err := json.Marshal(&healthReport); err != nil {
		logger.Error(err)
	} else {
		logger.Info("Health report: ", string(bytes))
	}
}

func (l *HealthLogger) Run() {
	ticker := time.NewTicker(l.cfg.Frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if l.ctx.Err() != nil {
				l.Stop()
			} else {
				l.LogHealth()
			}

		case <-l.done:
			break
		}
	}
}

func (l *HealthLogger) Stop() {
	close(l.done)
}

func NewHealthLogger(ctx context.Context, cfg *HealthLoggerConfig) *HealthLogger {
	return &HealthLogger{
		ctx:  ctx,
		cfg:  cfg,
		done: make(chan struct{}),
	}
}

func init() {
	OnEvent(EventStart, PhaseAfter, createHealthLogger)
	OnEvent(EventStop, PhaseBefore, closeHealthLogger)
}

func createHealthLogger(ctx context.Context) error {
	logger.Info("Starting health logger")

	cfg := config.FromContext(ctx)
	if cfg == nil {
		return errors.New("Config not found in context")
	}

	healthLoggerConfig := &HealthLoggerConfig{}
	if err := cfg.Populate(healthLoggerConfig, rootKeyHealth); err != nil {
		return err
	}

	if !healthLoggerConfig.Enabled {
		return nil
	}

	healthLogger = NewHealthLogger(ctx, healthLoggerConfig)
	go healthLogger.Run()
	return nil
}

func closeHealthLogger(ctx context.Context) error {
	logger.Info("Stopping health logger")

	if healthLogger != nil {
		healthLogger.Stop()
	}
	return nil
}
