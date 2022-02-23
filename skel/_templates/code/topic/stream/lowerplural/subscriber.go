//go:generate mockery --inpackage --name=UpperCamelSingularSubscriberApi --structname=MockUpperCamelSingularSubscriber --filename mock_subscriber_lowersingular.go
package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/topic/api"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

const (
	contextKeyUpperCamelSingularSubscriber = contextKey("UpperCamelSingularSubscriber")
)

type UpperCamelSingularSubscriberApi interface {
	OnUpperCamelSingularMessage(ctx context.Context, message api.UpperCamelSingularMessage) error
}

type lowerCamelSingularSubscriber struct {
}

func (s *lowerCamelSingularSubscriber) OnUpperCamelSingularMessage(ctx context.Context, message api.UpperCamelSingularMessage) error {
	logger.WithContext(ctx).Debugf("Handling message for lowersingular %q", message.Id.String())
	return nil
}

func newUpperCamelSingularSubscriber(ctx context.Context) UpperCamelSingularSubscriberApi {
	service := UpperCamelSingularSubscriberFromContext(ctx)
	if service == nil {
		service = &lowerCamelSingularSubscriber{}
	}
	return service
}

func UpperCamelSingularSubscriberFromContext(ctx context.Context) UpperCamelSingularSubscriberApi {
	value, _ := ctx.Value(contextKeyUpperCamelSingularSubscriber).(UpperCamelSingularSubscriberApi)
	return value
}

func ContextWithUpperCamelSingularSubscriber(ctx context.Context, subscriber UpperCamelSingularSubscriberApi) context.Context {
	return context.WithValue(ctx, contextKeyUpperCamelSingularSubscriber, subscriber)
}

func init() {
	app.OnRootEvent(app.EventStart, app.PhaseDuring, func(ctx context.Context) error {
		subscriber := newUpperCamelSingularSubscriber(ctx)
		return stream.AddListener(topicUpperCamelSingular, func(msg *message.Message) error {
			ctx := msg.Context()
			var lowerCamelSingularMessage api.UpperCamelSingularMessage
			if err := json.Unmarshal(msg.Payload, &lowerCamelSingularMessage); err != nil {
				err = errors.Wrap(err, "Failed to decode lowersingular message")
				logger.WithError(err).WithContext(ctx).Errorf("Invalid message received on topic %s", topicUpperCamelSingular)
				return err
			}
			return subscriber.OnUpperCamelSingularMessage(ctx, lowerCamelSingularMessage)
		})
	})
}
