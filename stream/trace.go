package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
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
		return nil
	}

	operationName := "kafka receive " + a.cfg.Destination
	ctx, span := trace.NewSpan(msg.Context(), operationName, ext.RPCServerOption(incomingContext))
	defer span.Finish()
	msg.SetContext(ctx)

	return a.action(msg)
}

func TraceActionDecorator(cfg *BindingConfiguration, action ListenerAction) ListenerAction {
	traceAction := &TraceSubscriberAction{
		action: action,
		cfg:    cfg,
	}
	return traceAction.Call
}
