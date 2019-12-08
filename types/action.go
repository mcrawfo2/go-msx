package types

import "context"

type ActionFunc func(ctx context.Context) error
