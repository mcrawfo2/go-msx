// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/swaggest/refl"
	"reflect"
	"sort"
)

type MessagePublisherBuilder struct {
	Name             string
	ChannelPublisher *ChannelPublisher
	Outputs          interface{}
	Filters          types.ActionFilters
	Documentors      ops.Documentors[MessagePublisher]
}

func (b *MessagePublisherBuilder) WithDocumentor(documentor ops.Documentor[MessagePublisher]) *MessagePublisherBuilder {
	b.Documentors = b.Documentors.WithDocumentor(documentor)
	return b
}

func (b *MessagePublisherBuilder) WithDecorator(deco types.ActionFuncDecorator) *MessagePublisherBuilder {
	return b.WithFilter(types.NewOrderedDecorator(b.Filters.NextCustomOrder(), deco))
}

func (b *MessagePublisherBuilder) WithFilter(filter types.ActionFilter) *MessagePublisherBuilder {
	filters := append(types.ActionFilters{}, b.Filters...)
	b.Filters = append(filters, filter)
	sort.Sort(b.Filters)
	return b
}

var ErrMessagePublisherBuildFailure = errors.New("Missing value for publisher field")

func (b *MessagePublisherBuilder) Build() (mp *MessagePublisher, err error) {
	portStructType := refl.DeepIndirect(reflect.TypeOf(b.Outputs))
	outputPort, err := PortReflector{}.ReflectOutputPort(portStructType)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to reflect outputs from port struct for operation %q", b.Name)
	}

	result := &MessagePublisher{
		name:             b.Name,
		channelPublisher: b.ChannelPublisher,
		filters:          b.Filters,
		outputPort:       outputPort,
		documentors:      b.Documentors,
	}

	RegisterMessagePublisher(result)

	return result, nil
}

func NewMessagePublisherBuilder(_ context.Context, channelPublisher *ChannelPublisher, name string, outputs interface{}) (*MessagePublisherBuilder, error) {
	if nil == channelPublisher {
		return nil, errors.Wrap(ErrMessagePublisherBuildFailure, "channelPublisher")
	} else if nil == outputs {
		return nil, errors.Wrap(ErrMessagePublisherBuildFailure, "outputs")
	}

	result := &MessagePublisherBuilder{
		Name:             name,
		ChannelPublisher: channelPublisher,
		Outputs:          outputs,
	}

	return result, nil
}

// MessagePublisher maps to asyncapi.Message
type MessagePublisher struct {
	name             string
	channelPublisher *ChannelPublisher
	filters          types.ActionFilters
	outputPort       *ops.Port
	documentors      ops.Documentors[MessagePublisher]
}

func (o MessagePublisher) Name() string {
	return o.name
}

func (o MessagePublisher) ContentType() string {
	return o.channelPublisher.Channel().binding.ContentType
}

func (o MessagePublisher) Publish(ctx context.Context, outputs interface{}) error {
	// Configure the response populator
	sink := new(MessageDataSink)
	encoder := WatermillMessageEncoder{Sink: sink}
	populator := &OutputsPopulator{
		Outputs:         &outputs,
		OutputPort:      o.outputPort,
		Channel:         types.OptionalOf(o.Channel().name),
		ContentType:     o.Channel().binding.ContentType,
		ContentEncoding: o.Channel().binding.ContentEncoding,
		Encoder:         encoder,
	}

	// Execute custom outputs struct validation if defined
	if err := o.ValidateOutputs(outputs); err != nil {
		return err
	}

	// Populate message from outputs
	if err := populator.PopulateOutputs(); err != nil {
		return errors.Wrap(err, "Failed to populate stream message")
	}

	// Validate output message matches schema
	if err := o.ValidateMessage(sink.Message(ctx)); err != nil {
		return err
	}

	// Publish the message
	return trace.NewOperation(
		o.name,
		func(ctx context.Context) error {
			return o.channelPublisher.Publish(ctx,
				sink.Payload,
				sink.Metadata)
		}).
		WithFilters(o.filters).
		Run(ctx)
}

func (o MessagePublisher) ValidateMessage(channelName types.Optional[string], msg *message.Message) error {
	messageDataSource := NewMessageDataSource(
		channelName.OrElse(o.Channel().Name()),
		msg)

	messageDecoder := NewMessageDecoder(messageDataSource,
		o.Channel().binding.ContentType,
		o.Channel().binding.ContentEncoding)

	messageValidator := NewMessageValidator(o.outputPort, messageDecoder)

	return messageValidator.ValidateMessage()
}

func (o MessagePublisher) ValidateOutputs(outputs any) error {
	// Auto-validation for validatable Port Struct
	portStructValue := reflect.ValueOf(outputs)
	if err := validate.ValidateValue(portStructValue); err != nil {
		errs := &ops.ValidationFailure{
			Path:     "message",
			Children: make(map[string]*ops.ValidationFailure),
		}

		return errs.Apply(err)
	}

	return nil
}

func (o MessagePublisher) Documentor(pred ops.DocumentorPredicate[MessagePublisher]) ops.Documentor[MessagePublisher] {
	return o.documentors.Documentor(pred)
}

func (o MessagePublisher) Channel() *Channel {
	return o.channelPublisher.Channel()
}

func (o MessagePublisher) OutputPort() *ops.Port {
	return o.outputPort
}

type messagePublishersList []*MessagePublisher

func (m messagePublishersList) Lookup(channel string, message string) *MessagePublisher {
	for _, mp := range m {
		if mp.Name() == message && mp.Channel().Name() == channel {
			return mp
		}
	}
	return nil
}

func (m messagePublishersList) AllByChannel(channel string) []*MessagePublisher {
	var results []*MessagePublisher
	for _, mp := range m {
		if channel == mp.Channel().Name() {
			results = append(results, mp)
		}
	}
	return results
}

var registeredMessagePublishers = messagePublishersList{}

func RegisterMessagePublisher(p *MessagePublisher) {
	// Do not add registration twice
	if nil != registeredMessagePublishers.Lookup(p.Channel().Name(), p.Name()) {
		return
	}

	registeredMessagePublishers = append(registeredMessagePublishers, p)
}

func RegisteredMessagePublishers(channel string) []*MessagePublisher {
	return registeredMessagePublishers.AllByChannel(channel)
}
