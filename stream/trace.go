package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type TraceSubscriberAction struct {
	action ListenerAction
	cfg    *BindingConfiguration
}

func (a *TraceSubscriberAction) Call(msg *message.Message) (err error) {
	textMap := opentracing.TextMapCarrier(msg.Metadata)
	incomingContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, textMap)
	if err != nil {
		logger.WithError(err).Error("Invalid trace context.")
	}

	operationName := fmt.Sprintf("%s.receive.%s", a.cfg.Binder, a.cfg.Destination)

	ctx, span := trace.NewSpan(msg.Context(), operationName,
		opentracing.FollowsFrom(incomingContext),
		ext.SpanKindConsumer)
	defer span.Finish()
	msg.SetContext(ctx)

	span.SetTag(trace.FieldDirection, "receive")
	span.SetTag(trace.FieldTopic, a.cfg.Destination)

	err = a.action(msg)
	if err != nil {
		span.LogFields(trace.Error(err))
	}

	return err
}

func TraceActionInterceptor(cfg *BindingConfiguration, action ListenerAction) ListenerAction {
	traceAction := &TraceSubscriberAction{
		action: action,
		cfg:    cfg,
	}
	return traceAction.Call
}

type TracePublisher struct {
	publisher Publisher
	cfg       *BindingConfiguration
}

func (t *TracePublisher) Publish(msg *message.Message) error {
	if msg == nil {
		return nil
	}

	operationName := fmt.Sprintf("%s.send.%s", t.cfg.Binder, t.cfg.Destination)
	ctx, span := trace.NewSpan(msg.Context(), operationName, ext.SpanKindProducer)
	span.SetTag(trace.FieldDirection, "send")
	span.SetTag(trace.FieldTopic, t.cfg.Destination)
	span.SetTag(trace.FieldTransport, t.cfg.Binder)

	defer span.Finish()
	msg.SetContext(ctx)

	// Decorate all of the messages with the trace metadata
	if msg.Metadata == nil {
		msg.Metadata = make(message.Metadata)
	}

	err := opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(msg.Metadata))
	if err != nil {
		logger.WithError(err).WithContext(ctx).Warn("Failed to apply trace context to outgoing message")
	}

	return t.publisher.Publish(msg)
}

func (t *TracePublisher) Close() error {
	return t.publisher.Close()
}

func NewTracePublisher(publisher Publisher, cfg *BindingConfiguration) Publisher {
	return &TracePublisher{
		publisher: publisher,
		cfg:       cfg,
	}
}
