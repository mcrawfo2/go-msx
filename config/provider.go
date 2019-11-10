package config

import (
	"context"
)

type Provider interface {
	Description() string
	Load(ctx context.Context) (map[string]string, error)
}
