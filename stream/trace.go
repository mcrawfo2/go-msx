// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
)

type TraceSubscriberAction struct {
	action ListenerAction
	cfg    *BindingConfiguration
}

func (a *TraceSubscriberAction) Call(msg *message.Message) (err error) {
	ctx := msg.Context()

	incomingSpanContext, err := trace.TextMapCarrier(msg.Metadata).Extract()
	if err != nil {
		logger.WithError(err).Error("Missing or invalid trace context.")
	} else {
		ctx = trace.ContextWithParentContext(ctx, incomingSpanContext)
	}

	operationName := fmt.Sprintf("%s.receive.%s", a.cfg.Binder, a.cfg.Destination)

	options := []trace.StartSpanOption{
		trace.StartWithTag(trace.FieldSpanKind, trace.SpanKindConsumer),
		trace.StartWithTag(trace.FieldDirection, "receive"),
		trace.StartWithTag(trace.FieldTopic, a.cfg.Destination),
		trace.StartWithTag(trace.FieldSpanType, "stream"),
	}

	if incomingSpanContext != nil {
		options = append(options,
			trace.StartWithRelated(trace.RefFollowsFrom, incomingSpanContext))
	}

	ctx, span := trace.NewSpan(ctx, operationName, options...)
	defer span.Finish()
	msg.SetContext(ctx)

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
	ctx, span := trace.NewSpan(msg.Context(), operationName,
		trace.StartWithTag(trace.FieldSpanKind, trace.SpanKindProducer),
		trace.StartWithTag(trace.FieldDirection, "send"),
		trace.StartWithTag(trace.FieldTopic, t.cfg.Destination),
		trace.StartWithTag(trace.FieldTransport, t.cfg.Binder),
		trace.StartWithTag(trace.FieldSpanType, "stream"),
	)

	defer span.Finish()
	msg.SetContext(ctx)

	// Decorate all the messages with the trace metadata
	if msg.Metadata == nil {
		msg.Metadata = make(message.Metadata)
	}

	err := trace.TextMapCarrier(msg.Metadata).Inject(span.Context())
	if err != nil {
		return err
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
