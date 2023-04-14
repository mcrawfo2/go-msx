//#id channelpackage ${async.channel.package}
//#id AsyncMessageSubscriber ${async.upmsgtype}Subscriber
//#id asyncMessageSubscriber ${async.msgtype}Subscriber
//#id newAsyncMessageSubscriber new${async.upmsgtype}Subscriber
//#id contextKeyAsyncMessageSubscriber contextKey${async.upmsgtype}Subscriber
//#id ContextAsyncMessageSubscriber Context${async.upmsgtype}Subscriber
//#id AsyncMessageHandler ${async.upmsgtype}Handler
//#id dropAsyncMessageHandler drop${async.upmsgtype}Handler
//#id asyncMessageInput ${async.msgtype}Input
package channelpackage

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

//go:generate mockery --name=AsyncMessageHandler --testonly --case=snake --inpackage --with-expecter

//#var dependencies
//#ignore
type AsyncMessageHandler interface {
	OnAsyncMessage(ctx context.Context, payload api.AsyncMessage) error
}

type dropAsyncMessageHandler struct{}

func (_ dropAsyncMessageHandler) OnAsyncMessage(ctx context.Context, payload api.AsyncMessage) error {
	logger.Error("No handler assigned to AsyncMessage message subscription.  Dropping message.")
	return nil
}
//#endignore

// Context

const contextKeyAsyncMessageSubscriber = contextKeyNamed("AsyncMessageSubscriber")

func ContextAsyncMessageSubscriber() types.ContextKeyAccessor[*streamops.MessageSubscriber] {
	return types.NewContextKeyAccessor[*streamops.MessageSubscriber](contextKeyAsyncMessageSubscriber)
}

// Constructor

//#var implementation
//#ignore

type asyncMessageInput struct {
	Payload   api.AsyncMessage `in:"body"`
}

//#endignore

func newAsyncMessageSubscriber(ctx context.Context) (*streamops.MessageSubscriber, error) {
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
	var handler AsyncMessageHandler = dropAsyncMessageHandler{}

	sb, err := streamops.NewMessageSubscriberBuilder(ctx, cs, "${async.upmsgtype}")
	if err != nil {
		return nil, err
	}

	svc, err := sb.
		WithInputs(asyncMessageInput{}).
		WithDecorator(service.DefaultServiceAccount).
		WithHandler(
			//#var handler suffix=,
			//#ignore
			func(ctx context.Context, in *asyncMessageInput) error {
				return delegate.OnAsyncMessage(ctx, in.Payload)
			},
			//#endignore
		).
		WithDocumentor(doc).
		Build()
	if err != nil {
		return nil, err
	}

	return svc, nil
}

// Singleton

var asyncMessageSubscriber = types.NewSingleton(
	newAsyncMessageSubscriber,
	ContextAsyncMessageSubscriber)

// Instantiate

func init() {
	app.OnCommandsEvent(
		[]string{app.CommandRoot, app.CommandAsyncApi},
		app.EventStart,
		app.PhaseBefore,
		func(ctx context.Context) (err error) {
			_, err = asyncMessageSubscriber.Factory(ctx)
			return
		})
}
