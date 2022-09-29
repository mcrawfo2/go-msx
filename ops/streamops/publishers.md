# Publishers

Stream Operations Publishers are used to publish messages on streams. 
They consist of a number of components:

| Component                        | AsyncApi Documentor                 | Documentation Model |
|----------------------------------|-------------------------------------|---------------------|
| Your Message Publisher (service) | -                                   | -                   |
| Your Output Port (struct)        | asyncapi.MessagePublisherDocumentor | jsonschema.Schema   |
| Your Payload DTO (struct)        | asyncapi.MessagePublisherDocumentor | jsonschema.Schema   |
| streamops.MessagePublisher       | asyncapi.MessagePublisherDocumentor | asyncapi.Message    |
| streamops.ChannelPublisher       | asyncapi.ChannelPublisherDocumentor | asyncapi.Operation  |
| streamops.Channel                | asyncapi.ChannelDocumentor          | asyncapi.Channel    |

## Components

**Channel**

The channel component represents the stream itself (SQS or Kafka topic, Redis stream, Go channel, SQLDB table, etc).
It is implemented as a singleton that should be created after configuration but before start-up.

**Channel Publisher**

The channel publisher component represents the set of publishable messages for a given stream.
It is implemented as a service that should be created as a dependency of your message publisher.

**Message Publisher**

The message publisher component represents one of the publishable messages for a given stream.
It is implemented as a service created after configuration but before start-up.  
Notice that it has a defined API interface for mocking, and should be mocked by dependent services
during testing.

**Output Port**

The message port contains a mapping of fields to be set on the outgoing message.
Each field will be mapped to a header or body field based on the struct tags.

**Message Payload DTO**

The payload DTO will contain the body of message that is to be published.
Before dispatch to the underlying stream, the message will be validated using
the JSON-schema annotations on your DTO.

## Generation

It is strongly advised to auto-generate these components and customize them afterwards.
See [Channels](../../skel/asyncapi/channels.md) and [AsyncApi](../../skel/asyncapi/spec.md)
for details about generation.
