package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/migrate"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	repomigrate "cto-github.cisco.com/NFV-BU/go-msx/repository/migrate"
)

func init() {
	OnEvent(EventCommand, CommandMigrate, func(ctx context.Context) error {
		// Only during migrate command
		OnEvent(EventConfigure, PhaseAfter, cassandra.CreateKeyspaceForPool)

		OnEvent(EventStart, PhaseBefore, func(ctx context.Context) error {
			manifest, err := migrate.NewManifest(config.FromContext(ctx))
			if err != nil {
				return err
			}

			RegisterContextInjector(func(ctx context.Context) context.Context {
				return migrate.ContextWithManifest(ctx, manifest)
			})
			return nil
		})

		return nil
	})

	repomigrate.RegisterMigrator(migrate.Migrate)
}
