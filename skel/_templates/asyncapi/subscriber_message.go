package ${async.channel.package}

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"${app.packageurl}/internal/stream/${async.channel.package}/api"
)

// Dependencies

//go:generate mockery --inpackage --name=${async.upmsgtype}Handler --structname=Mock${async.upmsgtype}Handler --filename mock_${async.upmsgtype}Handler.go

${dependencies}

// Context

const contextKey${async.upmsgtype}Subscriber = contextKeyNamed("${async.upmsgtype}Subscriber")

func Context${async.upmsgtype}Subscriber() types.ContextKeyAccessor[*streamops.MessageSubscriber] {
	return types.NewContextKeyAccessor[*streamops.MessageSubscriber](contextKey${async.upmsgtype}Subscriber)
}

// Constructor

${implementation}

func new${async.upmsgtype}Subscriber(ctx context.Context) (*streamops.MessageSubscriber, error) {
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

	// TODO: specify your message handler here
	var handler ${async.upmsgtype}Handler = drop${async.upmsgtype}Handler{}

	sb, err := streamops.NewMessageSubscriberBuilder(ctx, cs, "${async.upmsgtype}")
	if err != nil {
		return nil, err
	}

	svc, err := sb.
		WithInputs(${async.msgtype}Input{}).
		WithDecorator(service.DefaultServiceAccount).
		WithHandler(${handler}).
		WithDocumentor(doc).
		Build()
	if err != nil {
		return nil, err
	}

	return svc, nil
}

// Singleton

var ${async.msgtype}Subscriber = types.NewSingleton(
	new${async.upmsgtype}Subscriber,
	Context${async.upmsgtype}Subscriber)

// Instantiate

func init() {
	app.OnCommandsEvent(
		[]string{app.CommandRoot, app.CommandAsyncApi},
		app.EventStart,
		app.PhaseBefore,
		func(ctx context.Context) (err error) {
			_, err = ${async.msgtype}Subscriber.Factory(ctx)
			return
		})
}
