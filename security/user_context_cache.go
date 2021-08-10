package security

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

var logger = log.NewLogger("msx.security")

type UserContextCacheApi interface {
	Latest(ctx context.Context) (result *UserContext, err error)
	Expire()
}

type UserContextFactory func(ctx context.Context) (*UserContext, error)

type UserContextCache struct {
	factory UserContextFactory
	latest  *UserContext
	worker  *types.Worker
}

func (c *UserContextCache) Latest(ctx context.Context) (result *UserContext, err error) {
	err = c.worker.Run(
		func(ctx context.Context) error {
			if c.latest != nil {
				result = c.latest
				return nil
			}

			result, err = c.factory(ctx)
			if err != nil {
				result = nil
			}
			return err
		},
		types.JobContext(ctx),
		types.JobDecorator(log.RecoverLogDecorator(logger)),
	)
	return
}

func (c *UserContextCache) Expire() {
	_ = c.worker.Run(func(ctx context.Context) error {
		c.latest = nil
		return nil
	})
}

func NewUserContextCache(ctx context.Context, factory UserContextFactory) *UserContextCache {
	return &UserContextCache{
		factory: factory,
		worker:  types.NewWorker(ctx),
	}
}
