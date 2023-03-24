//#id channelpackage ${async.channel.package}
//#id asyncMessagePublisher ${async.msgtype}Publisher
//#id AsyncMessagePublisher ${async.upmsgtype}Publisher
//#id NewAsyncMessagePublisher New${async.upmsgtype}Publisher
//#id contextKeyAsyncMessagePublisher  contextKey${async.upmsgtype}Publisher
//#id ContextAsyncMessagePublisher  Context${async.upmsgtype}Publisher
//#id asyncMessageOutput  ${async.msgtype}Output
package channelpackage

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	//#var imports
)

// Dependencies

//go:generate mockery --name=AsyncMessagePublisher --testonly --case=snake --inpackage --with-expecter

//#var dependencies
//#ignore

type AsyncMessagePublisher interface {
	PublishAsyncMessage(ctx context.Context, payload api.AsyncMessage) error
}

//#endignore

// Context

const contextKeyAsyncMessagePublisher = contextKeyNamed("AsyncMessagePublisher")

func ContextAsyncMessagePublisher() types.ContextKeyAccessor[AsyncMessagePublisher] {
	return types.NewContextKeyAccessor[AsyncMessagePublisher](contextKeyAsyncMessagePublisher)
}

// Implementation

type asyncMessagePublisher struct {
	messagePublisher *streamops.MessagePublisher
}

//#var implementation
//#ignore

type asyncMessagePublisher struct {
	messagePublisher *streamops.MessagePublisher
}

type asyncMessageOutput struct {
	EventType string           `out:"header=eventType" const:"statusChange"`
	Payload   api.AsyncMessage `out:"body"`
}

func (p asyncMessagePublisher) PublishStatusChangeResponse(ctx context.Context, payload api.AsyncMessage) error {
	return p.messagePublisher.Publish(ctx, asyncMessageOutput{
		Payload: payload,
	})
}

//#endignore

// Constructor

func NewAsyncMessagePublisher(ctx context.Context) (AsyncMessagePublisher, error) {
	svc := ContextAsyncMessagePublisher().Get(ctx)
	if svc == nil {
		doc := new(asyncapi.MessagePublisherDocumentor).
			WithMessage(new(asyncapi.Message).
				WithTitle("${async.msgtype.human}").
				WithSummary("Notifies subscribers of ${async.msgtype.human}.").
				WithTags(
					*asyncapi.NewTag("${async.msgtype}"),
				))

		cp, err := newChannelPublisher(ctx)
		if err != nil {
			return nil, err
		}

		mpb, err := streamops.NewMessagePublisherBuilder(ctx, cp, "${async.upmsgtype}", asyncMessageOutput{})

		mp, err := mpb.
			WithDocumentor(doc).
			Build()
		if err != nil {
			return nil, err
		}

		svc = &asyncMessagePublisher{
			messagePublisher: mp,
		}
	}

	return svc, nil
}

// Instantiation

func init() {
	app.OnCommandsEvent(
		[]string{app.CommandRoot, app.CommandAsyncApi},
		app.EventStart,
		app.PhaseBefore,
		func(ctx context.Context) (err error) {
			_, err = NewAsyncMessagePublisher(ctx)
			return
		})
}
