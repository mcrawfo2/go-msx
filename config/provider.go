package config

import (
	"context"
)

type Provider interface {
	Load(ctx context.Context) (map[string]string, error)
}
