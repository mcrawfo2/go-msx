# MSX Stream Operations

A Streaming Operation library compatible with [AsyncApi](https://asyncapi.com) 2.x documentation.

## Terminology

Message
: A discrete unit of communication between a publisher and
a set of subscribers. Must include data (payload) and may
include metadata (headers).

Channel
: A source or destination for message delivery between
publishers and subscribers.

Subscriber
: A receiver of a sequence of messages from a channel.

Publisher
: A sender of a sequence of messages to a channel.

AsyncApi
: Documentation standard for describing event-based and streaming
  message transports such as Kafka, Redis Streams, Amazon SQS.
  Describes messages, channels, publishers, subscribers, servers,
  security, and other related concerns.  Comparable to OpenApi, which
  describes REST message transports.

DTO
: Data Transfer Object.  Used for serialization and deserialization of
  externally sourced or directed structured values.

Port
: Description of the interface between the Stream Operations subsystem and
  your message publisher or subscriber. Can include headers, filters, and
  must include a payload DTO.

## Components

The following table compares the pattern of components across AsyncApi,
Stream Operations, and HTTP components.

| AsyncApi         | Stream Publisher | Stream Subscriber | HTTP              | Purpose                   |
|------------------|------------------|-------------------|-------------------|---------------------------|
| Channel          | Channel          | Channel           | Controller        | Domain ingress and egress |
| Operation        | ChannelPublisher | ChannelSubscriber | Router            | Dispatch to endpoints     |
| Message          | MessagePublisher | MessageSubscriber | Endpoint          | Event processing          |
| Header, Payload  | Output Port      | Input Port        | Request/Response  | Exchanged data            | 
