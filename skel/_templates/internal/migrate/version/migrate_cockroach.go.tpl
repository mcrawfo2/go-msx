package ${app.migrateVersion}

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb/migrate"
)

func init() {
	app.OnEvent(app.EventCommand, app.CommandMigrate, func(ctx context.Context) error {
		app.OnEvent(app.EventStart, app.PhaseDuring, func(ctx context.Context) error {
			return migrate.
				ManifestFromContext(ctx).
				AddSqlResourceMigrations(
					resource.References("*.sql")...
				)
		})
		return nil
	})
}
