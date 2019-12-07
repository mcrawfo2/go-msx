package main

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/migrate"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const (
	appName = "${app.name}"
)

func init() {
	app.OnEvent(app.EventCommand, app.CommandMigrate, func(ctx context.Context) error {
		app.OnEvent(app.EventStart, app.PhaseDuring, addMigrations)
		return nil
	})
}

func addMigrations(ctx context.Context) error {
	manifest := migrate.ManifestFromContext(ctx)

	return types.ErrorList{
		manifest.AddCqlStringMigration("3.8.0.1", "Create first table", "CREATE TABLE first (value text PRIMARY KEY)"),
	}.Filter()
}

func main() {
	app.Run(appName)
}
