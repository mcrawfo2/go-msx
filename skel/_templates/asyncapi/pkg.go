package ${async.channel.package}

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
)

type contextKeyNamed string

var channel *streamops.Channel

func init() {
	app.OnEvent(app.EventConfigure, app.PhaseAfter, func(ctx context.Context) (err error) {
		channel, err = streamops.NewChannel(ctx, "${async.channel}")
		if err != nil {
			return
		}

		channel.WithDocumentor(new(asyncapi.ChannelDocumentor).
			WithChannelItem(new(asyncapi.ChannelItem).
				WithDescription("Description of the ${async.channel} channel.")))

		return nil
	})
}
