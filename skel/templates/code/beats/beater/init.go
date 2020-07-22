package beater

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx-beats/meta"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
)

func init() {
	var beater *Beater

	meta.SetFieldsResource("/internal/_meta/fields.yml")

	app.OnEvent(app.EventCommand, app.CommandRoot, func(ctx context.Context) error {
		app.OnEvent(app.EventStart, app.PhaseBefore, func(ctx context.Context) (err error) {
			beater, err = newBeater(ctx)
			return
		})

		app.OnEvent(app.EventReady, app.PhaseDuring, func(ctx context.Context) error {
			go beater.Run(ctx)
			return nil
		})

		return nil
	})
}
