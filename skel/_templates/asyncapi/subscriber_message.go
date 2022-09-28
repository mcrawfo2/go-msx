package ${async.channel.package}

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"${app.packageurl}/internal/stream/${async.channel.package}/api"
	// "${app.packageurl}/internal/${async.domain.package}"
)

// Dependencies

//go:generate mockery --inpackage --name=${async.upmsgtype}Handler --structname=Mock${async.upmsgtype}Handler --filename mock_${async.upmsgtype}Handler.go

${dependencies}

const contextKey${async.upmsgtype}Subscriber = contextKeyNamed("${async.upmsgtype}Subscriber")

func Context${async.upmsgtype}Subscriber() types.ContextKeyAccessor[*streamops.MessageSubscriber] {
	return types.NewContextKeyAccessor[*streamops.MessageSubscriber](contextKey${async.upmsgtype}Subscriber)
}

// Constructor

${implementation}

func new${async.upmsgtype}Subscriber(ctx context.Context) (*streamops.MessageSubscriber, error) {
	svc := Context${async.upmsgtype}Subscriber().Get(ctx)
	if svc == nil {

		doc := new(asyncapi.MessageSubscriberDocumentor).
			WithMessage(new(asyncapi.Message).
				WithTitle("${async.msgtype.human}").
				WithSummary("${async.msgtype.human}").
				WithTags(
					*asyncapi.NewTag("${async.upmsgtype}"),
				))

		cs, err := newChannelSubscriber(ctx)
		if err != nil {
			return nil, err
		}

		var handler ${async.upmsgtype}Handler
		// handler, err = ${async.domain.package}.New${async.domain.uppercamelsingular}Service(ctx)
		// if err != nil {
		//	   return nil, err
		// }

		sb, err := streamops.NewMessageSubscriberBuilder(ctx, cs, "${async.upmsgtype}")
		if err != nil {
			return nil, err
		}

		svc, err = sb.
			WithInputs(${async.msgtype}Input{}).
			WithDecorator(service.DefaultServiceAccount).
			WithHandler(${handler}).
			WithDocumentor(doc).
			Build()
		if err != nil {
			return nil, err
		}
	}

	return svc, nil
}

// Instantiate

func init() {
	app.OnCommandsEvent(
		[]string{app.CommandRoot, app.CommandAsyncApi},
		app.EventStart,
		app.PhaseBefore,
		func(ctx context.Context) error {
			_, err := new${async.upmsgtype}Subscriber(ctx)
			return err
		})
}
