# Subscribers

Stream Operations Subscribers are used to publish messages on streams.
They consist of a number of components:

| Component                         | AsyncApi Documentor                  | Documentation Model |
|-----------------------------------|--------------------------------------|---------------------|
| Your Message Subscriber (service) | -                                    | -                   |
| Your Input Port (struct)          | asyncapi.MessageSubscriberDocumentor | jsonschema.Schema   |
| Your Payload DTO (struct)         | asyncapi.MessageSubscriberDocumentor | jsonschema.Schema   |
| streamops.MessageSubscriber       | asyncapi.MessageSubscriberDocumentor | asyncapi.Message    |
| streamops.ChannelSubscriber       | asyncapi.ChannelSubscriberDocumentor | asyncapi.Operation  |
| streamops.Channel                 | asyncapi.ChannelDocumentor           | asyncapi.Channel    |

## Components

**Channel**

The channel component represents the stream itself (SQS or Kafka topic, Redis stream, Go channel, SQLDB table, etc).
It is implemented as a singleton that should be created after configuration but before start-up.

**Channel Subscriber**

The channel subscriber component represents the set of subscribable messages for a given stream.
It is implemented as a service, and should have one of your application services as a dependency.

**Message Subscriber**

The message subscriber component represents one of the publishable messages for a given stream.
It is implemented as a service created after configuration but before start-up.  
Notice that it has a defined API interface for mocking, and should be mocked by dependent services
during testing.

**Input Port**

The message port contains a mapping of fields to be set from the incoming message.
Each field will be mapped from a header or body field based on the struct tags.

**Payload DTO**

The payload DTO will contain the parsed message that is subscribed.
Before dispatch to your subscriber, the message will be validated using
the JSON-schema annotations and any `Validatable` interface implementation
on your DTO.

## Generation

It is strongly advised to auto-generate these components and customize them afterwards.
See [Channels](../../skel/asyncapi/channels.md) and [AsyncApi](../../skel/asyncapi/spec.md)
for details about generation.
