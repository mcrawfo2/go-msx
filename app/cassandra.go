package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/migrate"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	repomigrate "cto-github.cisco.com/NFV-BU/go-msx/repository/migrate"
)

func init() {
	repomigrate.RegisterMigrator(migrate.Migrate)

	OnEvent(EventStart, PhaseBefore, func(ctx context.Context) error {
		manifest, err := migrate.NewManifest(config.FromContext(ctx))
		if err != nil {
			return err
		}

		contextInjectors.Register(func(ctx context.Context) context.Context {
			return migrate.ContextWithManifest(ctx, manifest)
		})
		return nil
	})
}
