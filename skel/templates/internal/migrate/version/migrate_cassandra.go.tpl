package ${app.migrateVersion}

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/migrate"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
)

func init() {
	app.OnEvent(app.EventCommand, app.CommandMigrate, func(ctx context.Context) error {
		app.OnEvent(app.EventStart, app.PhaseDuring, func(ctx context.Context) error {
			return migrate.
				ManifestFromContext(ctx).
				AddCqlResourceMigrations(
					resource.References("*.cql")...
				)
		})
		return nil
	})
}
