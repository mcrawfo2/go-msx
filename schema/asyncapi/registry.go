// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
)

func DocumentStreams(ctx context.Context) error {
	for _, channel := range streamops.RegisteredChannels() {
		err := ops.DocumentorWithType[streamops.Channel](channel, DocType).
			OrElse(ChannelDocumentor{}).
			Document(channel)
		if err != nil {
			return errors.Wrapf(err, "Failed to document channel %q", channel.Name())
		}
	}
	return nil
}

func channelBindingsEns(channel *ChannelItem) {
	// TODO: Bindings
}

func operationKafkaBindings(operation *Operation) {
	operation.WithBindings(*new(BindingsObject).WithKafka(types.Pojo{
		"groupId":  "{APP_NAME}-{TOPIC_NAME}-GP",
		"clientId": "{APP_NAME}-{TOPIC_NAME}-{APP_INSTANCE_ID}",
	}))
}

func operationRedisBindings(operation *Operation) {
	operation.WithBindings(*new(BindingsObject))
}

func operationBindingsEns(operation *Operation, binder string) {
	if operation.Bindings == nil {
		switch binder {
		case "kafka":
			operationKafkaBindings(operation)
		case "redis":
			operationRedisBindings(operation)
		}
	}
}

func addPublisherOperation(topic string, operation Operation) {
	channelItem := Reflector.SpecEns().Channels[topic]
	channelItem.WithSubscribe(operation) // Opposite action in spec
	Reflector.SpecEns().WithChannelsItem(topic, channelItem)
}

func addPublisherOperationMessageChoice(topic string, messageRef Reference) {
	channelItem := Reflector.SpecEns().Channels[topic]
	operation := channelItem.SubscribeEns() // Opposite action in spec
	addOperationMessage(operation, messageRef)
}

func addSubscriberOperation(topic string, operation Operation) {
	channelItem := Reflector.SpecEns().Channels[topic]
	channelItem.WithPublish(operation) // Opposite action in spec
	Reflector.Spec.Channels[topic] = channelItem
}

func addSubscriberOperationMessageChoice(topic string, messageRef Reference) {
	channelItem := Reflector.SpecEns().Channels[topic]
	operation := channelItem.PublishEns() // Opposite action in spec
	addOperationMessage(operation, messageRef)
}

func addOperationMessage(operation *Operation, messageRef Reference) {
	messageChoices := operation.MessageEns()

	// Add this message type
	options := messageChoices.MessageOptionsEns().OneOf
	options = append(options, MessageChoices{
		Reference: &messageRef,
	})
	messageChoices.MessageOptions.WithOneOf(options...)
}

func createMessageReference(name string, message Message) Reference {
	Reflector.SpecEns().ComponentsEns().WithMessagesItem(name, MessageChoices{
		Message: &message,
	})

	return Reference{
		Ref: "#/components/messages/" + name,
	}
}
