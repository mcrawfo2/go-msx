// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/swaggest/refl"
	"reflect"
	"sort"
)

type MessageSubscriberBuilder struct {
	Name                 string
	ChannelSubscriber    *ChannelSubscriber
	Inputs               interface{}
	Handler              interface{}
	Filters              types.ActionFilters
	Documentors          ops.Documentors[MessageSubscriber]
	MetadataFilterValues map[string][]string
}

func (o *MessageSubscriberBuilder) WithInputs(portStruct interface{}) *MessageSubscriberBuilder {
	o.Inputs = portStruct
	return o
}

func (o *MessageSubscriberBuilder) WithHandler(fn interface{}) *MessageSubscriberBuilder {
	o.Handler = fn
	return o
}

func (o *MessageSubscriberBuilder) WithDecorator(deco types.ActionFuncDecorator) *MessageSubscriberBuilder {
	return o.WithFilter(types.NewOrderedDecorator(o.Filters.NextCustomOrder(), deco))
}

func (o *MessageSubscriberBuilder) WithFilter(filter types.ActionFilter) *MessageSubscriberBuilder {
	filters := append(types.ActionFilters{}, o.Filters...)
	o.Filters = append(filters, filter)
	sort.Sort(o.Filters)
	return o
}

func (o *MessageSubscriberBuilder) WithDocumentor(doc ops.Documentor[MessageSubscriber]) *MessageSubscriberBuilder {
	o.Documentors = o.Documentors.WithDocumentor(doc)
	return o
}

func (o *MessageSubscriberBuilder) WithMetadataFilterValues(headerName string, values ...string) *MessageSubscriberBuilder {
	if o.MetadataFilterValues == nil {
		o.MetadataFilterValues = make(map[string][]string)
	}
	o.MetadataFilterValues[headerName] = values
	return o
}

var ErrMessageSubscriberBuildFailure = errors.New("Missing value for subscriber field")

func (o *MessageSubscriberBuilder) Build() (ms *MessageSubscriber, err error) {
	if o.Handler == nil {
		return nil, errors.Wrap(ErrMessageSubscriberBuildFailure, "func")
	}

	argsTypeSet := types.NewTypeSet(
		messageMessageType,
		channelType)

	var inputPort *ops.Port
	if o.Inputs != nil {
		portStructType := refl.DeepIndirect(reflect.TypeOf(o.Inputs))
		inputPort, err = PortReflector{}.ReflectInputPort(portStructType)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to reflect inputs from port struct for operation %q", o.Name)
		}

		argsTypeSet.WithType(portStructType)
	}

	handler, err := types.NewHandler(o.Handler,
		types.NewHandlerValueTypeReflector(
			argsTypeSet,
			types.DefaultHandlerArgumentValueTypeSet,
		),
		nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create handler")
	}

	call := types.
		NewOperation(handler.Call).
		WithDecorator(types.RecoverErrorDecorator).
		Run

	result := &MessageSubscriber{
		name:                 o.Name,
		channelSubscriber:    o.ChannelSubscriber,
		inputPort:            inputPort,
		handler:              call,
		filters:              o.Filters,
		documentors:          o.Documentors,
		metadataFilterValues: o.MetadataFilterValues,
	}

	if err = result.channelSubscriber.AddMessageConsumer(result); err != nil {
		return nil, err
	}

	RegisterMessageSubscriber(result)

	return result, nil
}

func NewMessageSubscriberBuilder(_ context.Context, channelSubscriber *ChannelSubscriber, name string) (*MessageSubscriberBuilder, error) {
	if name == "" {
		return nil, errors.Wrap(ErrMessageSubscriberBuildFailure, "name")
	}
	if channelSubscriber == nil {
		return nil, errors.Wrapf(ErrMessageSubscriberBuildFailure, "channelSubscriber")
	}

	return &MessageSubscriberBuilder{
		Name:              name,
		ChannelSubscriber: channelSubscriber,
	}, nil
}

// MessageConsumer provides a mockable interface to consuming messages
type MessageConsumer interface {
	Name() string
	MetadataFilterValues(headerName string) ([]string, error)
	stream.MessageListener
}

// MessageSubscriber maps to asyncapi.Operation
type MessageSubscriber struct {
	name                 string
	channelSubscriber    *ChannelSubscriber
	inputPort            *ops.Port
	handler              types.ActionFunc
	filters              types.ActionFilters
	documentors          ops.Documentors[MessageSubscriber]
	metadataFilterValues map[string][]string
}

func (o MessageSubscriber) Name() string {
	return o.name
}

func (o MessageSubscriber) Channel() *Channel {
	return o.channelSubscriber.Channel()
}

func (o MessageSubscriber) InputPort() *ops.Port {
	return o.inputPort
}

func (o MessageSubscriber) ContentType() string {
	return o.channelSubscriber.Channel().binding.ContentType
}

func (o MessageSubscriber) Documentor(pred ops.DocumentorPredicate[MessageSubscriber]) ops.Documentor[MessageSubscriber] {
	return o.documentors.Documentor(pred)
}

// MetadataFilterValues returns a list of values matching this subscriber for the specified header name
func (o *MessageSubscriber) MetadataFilterValues(headerName string) ([]string, error) {
	if o.metadataFilterValues != nil {
		filterValues, ok := o.metadataFilterValues[headerName]
		if ok {
			return filterValues, nil
		}
	}

	if o.inputPort == nil {
		return nil, errors.Errorf(
			"Metadata filter field %q values not defined for message subscriber %q",
			headerName,
			o.name)
	}

	header := o.inputPort.Fields.First(
		ops.PortFieldHasGroup(FieldGroupStreamHeader),
		ops.PortFieldHasPeer(headerName),
	)

	if header == nil {
		return nil, errors.Errorf(
			"Metadata filter field %q not defined in input ports for message subscriber %q",
			headerName,
			o.name)
	}

	e := header.Enum()
	if e != nil {
		var results []string
		for _, v := range e {
			results = append(results, cast.ToString(v))
		}
		return results, nil
	}

	c := header.Const()
	if c != nil {
		return []string{
			cast.ToString(c),
		}, nil
	}

	return nil, errors.Errorf(
		"Metadata filter field %q values not defined in input port %q for message subscriber %q",
		headerName,
		header.Name,
		o.name)
}

func (o MessageSubscriber) inputs(msg *message.Message) (result interface{}, err error) {
	if o.inputPort == nil {
		return
	}

	source := NewMessageDataSource(
		o.Channel().Name(),
		msg)

	decoder := NewMessageDecoder(source,
		o.Channel().DefaultContentType(),
		o.Channel().DefaultContentEncoding())

	populator := ops.NewInputsPopulator(o.inputPort, decoder)

	result, err = populator.PopulateInputs()
	if err != nil {
		return
	}

	messageValidator := NewMessageValidator(o.inputPort, decoder)
	err = messageValidator.ValidateMessage()
	if err != nil {
		return
	}

	return
}

func (o MessageSubscriber) requestActionDecorator(msg *message.Message) types.ActionFuncDecorator {
	return func(action types.ActionFunc) types.ActionFunc {
		return func(ctx context.Context) error {
			inp, err := o.inputs(msg)
			if err != nil {
				return err
			}

			var inputType reflect.Type
			if o.inputPort != nil {
				inputType = o.inputPort.StructType
			}

			handlerContext := &MessageSubscriberHandlerContext{
				msg:       msg,
				channel:   o.Channel(),
				inputType: inputType,
				inputs:    inp,
			}

			ctx = types.ContextWithHandlerContext(ctx, handlerContext)
			ctx = contextMessageSubscriberHandlerContext.Set(ctx, handlerContext)

			return action(ctx)
		}
	}
}

func (o MessageSubscriber) OnMessage(msg *message.Message) (err error) {
	if o.handler == nil {
		return
	}

	ctx := msg.Context()

	handler := types.NewOperation(o.handler).
		WithFilters(o.filters).
		Run

	return trace.NewOperation("On"+o.name, handler).
		WithDecorator(o.requestActionDecorator(msg)).
		Run(ctx)
}

var contextContextInstance context.Context
var contextContextType = reflect.TypeOf(&contextContextInstance).Elem()
var messageMessageInstance *message.Message
var messageMessageType = reflect.TypeOf(&messageMessageInstance).Elem()
var channelInstance *Channel
var channelType = reflect.TypeOf(&channelInstance).Elem()
var errorInstance error
var errorType = reflect.TypeOf(&errorInstance).Elem()

// MessageSubscriberHandlerContext holds the data to be injected into the subscriber's handler function
type MessageSubscriberHandlerContext struct {
	msg       *message.Message
	channel   *Channel
	inputType reflect.Type
	inputs    interface{}
}

func (m *MessageSubscriberHandlerContext) Message() *message.Message {
	return m.msg
}

func (m *MessageSubscriberHandlerContext) Channel() *Channel {
	return m.channel
}

func (m *MessageSubscriberHandlerContext) Inputs() interface{} {
	return m.inputs
}

func (m *MessageSubscriberHandlerContext) GenerateArgument(ctx context.Context, t types.HandlerValueType) (result reflect.Value, err error) {
	switch t.ValueType {
	case contextContextType:
		result = reflect.ValueOf(ctx)
	case messageMessageType:
		result = reflect.ValueOf(m.msg)
	case channelType:
		result = reflect.ValueOf(m.channel)
	case m.inputType:
		result = reflect.ValueOf(m.inputs).Elem()
	default:
		err = errors.Wrapf(types.ErrUnknownValueType, "%v", t)
	}
	return

}

func (m *MessageSubscriberHandlerContext) HandleResult(t types.HandlerValueType, v reflect.Value) (err error) {
	switch t.ValueType {
	case errorType:
		erri := v.Interface()
		if erri != nil {
			err = erri.(error)
		}
	default:
		err = errors.Wrapf(types.ErrUnknownValueType, "%v", t)
	}
	return
}

// For handlers that want to access the handler context contents directly
const contextKeyMessageSubscriberHandlerContext = contextKey("MessageSubscriberHandlerContext")

var contextMessageSubscriberHandlerContext = types.NewContextKeySetter[*MessageSubscriberHandlerContext](contextKeyMessageSubscriberHandlerContext)

func ContextMessageSubscriberHandlerContext() types.ContextKeyGetter[*MessageSubscriberHandlerContext] {
	return types.NewContextKeyGetter[*MessageSubscriberHandlerContext](contextKeyMessageSubscriberHandlerContext)
}

type messageConsumerDispatcher struct {
	messageConsumer MessageConsumer
}

func (m messageConsumerDispatcher) Dispatch(msg *message.Message) error {
	return m.messageConsumer.OnMessage(msg)
}

type messageSubscribersList []*MessageSubscriber

func (m messageSubscribersList) Lookup(channel string, message string) *MessageSubscriber {
	for _, mp := range m {
		if mp.Name() == message && mp.Channel().Name() == channel {
			return mp
		}
	}
	return nil
}

func (m messageSubscribersList) AllByChannel(channel string) []*MessageSubscriber {
	var results []*MessageSubscriber
	for _, mp := range m {
		if channel == mp.Channel().Name() {
			results = append(results, mp)
		}
	}
	return results
}

var registeredMessageSubscribers = messageSubscribersList{}

func RegisterMessageSubscriber(p *MessageSubscriber) {
	// Do not add registration twice
	if nil != registeredMessageSubscribers.Lookup(p.Channel().Name(), p.Name()) {
		return
	}

	registeredMessageSubscribers = append(registeredMessageSubscribers, p)
}

func RegisteredMessageSubscribers(channel string) []*MessageSubscriber {
	return registeredMessageSubscribers.AllByChannel(channel)
}
