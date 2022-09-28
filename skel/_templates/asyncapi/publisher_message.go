package ${async.channel.package}

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	${imports}
)

// Dependencies

//go:generate mockery --inpackage --name=${async.upmsgtype}Publisher --structname=Mock${async.upmsgtype}Publisher --filename mock_${async.upmsgtype}Publisher.go

${dependencies}

const contextKey${async.upmsgtype}Publisher = contextKeyNamed("${async.upmsgtype}Publisher")

func Context${async.upmsgtype}Publisher() types.ContextKeyAccessor[${async.upmsgtype}Publisher] {
	return types.NewContextKeyAccessor[${async.upmsgtype}Publisher](contextKey${async.upmsgtype}Publisher)
}

// Implementation

type ${async.msgtype}Publisher struct {
	messagePublisher *streamops.MessagePublisher
}

${implementation}

func New${async.upmsgtype}Publisher(ctx context.Context) (${async.upmsgtype}Publisher, error) {
	svc := Context${async.upmsgtype}Publisher().Get(ctx)
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

		mpb, err := streamops.NewMessagePublisherBuilder(ctx, cp, "${async.msgtype}", ${async.msgtype}Output{})

		mp, err := mpb.
			WithDocumentor(doc).
			Build()
		if err != nil {
			return nil, err
		}

		svc = &${async.msgtype}Publisher{
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
			_, err = New${async.upmsgtype}Publisher(ctx)
			return
		})
}
