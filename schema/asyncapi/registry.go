package asyncapi

func RegisterChannel(topic string, channelItem ChannelItem) {
	Reflector.SpecEns().WithChannelsItem(topic, channelItem)
}

func RegisterMessage(name string, message Message) Reference {
	Reflector.SpecEns().ComponentsEns().WithMessagesItem(name, MessageChoices{
		Message: &message,
	})

	return Reference{
		Ref: "#/components/messages/" + name,
	}
}

func RegisterChannelSubscribers(topic string, operationId string, messageRefs ...Reference) {
	channelItem := Reflector.SpecEns().Channels[topic]
	operation := channelItem.PublishEns() // Opposite action in spec
	operation.WithID(operationId)
	for _, messageRef := range messageRefs {
		registerOperation(topic, channelItem, operation, messageRef)
	}
}

func RegisterChannelPublishers(topic string, operationId string, messageRefs ...Reference) {
	channelItem := Reflector.SpecEns().Channels[topic]
	operation := channelItem.SubscribeEns() // Opposite action in spec
	operation.WithID(operationId)
	for _, messageRef := range messageRefs {
		registerOperation(topic, channelItem, operation, messageRef)
	}
}

func registerOperation(topic string, channelItem ChannelItem, operation *Operation, messageRef Reference) {
	operation.WithBindings(BindingsObject{
		AdditionalProperties: map[string]interface{}{
			"$ref": "#/components/operationBindings/cpx",
		},
	})
	messageChoices := operation.MessageEns()

	// Add this message type
	options := messageChoices.MessageOptionsEns().OneOf
	// TODO: ensure no repetition of message choices
	options = append(options, MessageChoices{
		Reference: &messageRef,
	})
	messageChoices.MessageOptions.WithOneOf(options...)

	RegisterChannel(topic, channelItem)
}
