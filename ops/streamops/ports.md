# Ports

_Ports_ describe the various components of an incoming or outgoing message.  
An early incarnation of Ports is currently used to describe REST Controller 
Endpoint inputs, so hopefully they seem familiar.

## Input Ports

Input Ports specify fields to be extracted from an incoming stream message.
For example:

```go
type driftCheckResponseInput struct {
    EventType string                 `in:"header" const:"DriftCheck"`
    Payload   api.DriftCheckResponse `in:"body"`
}
```

In Java these would be implemented using a hypothetical function signature:
```java
void OnReceiveDriftCheckResponse(
  @Header(name="eventType") @Schema(constant="DriftCheck") String eventType, 
  @Body DriftCheckResponse body
)
```

### Struct Tags

Each field with an `in` struct tag will be automatically populated before being
passed to your Message Subscriber.  The full syntax of the `in` struct tag is as
follows:

    in:"<fieldGroup>[=<peerName>]"

The `in` struct tag contains the following subcomponents:

`<fieldGroup>`
: (Required) The name of the message part from which the value will be extracted.
Valid field groups for streaming operations are:

* `header` - Message metadata/headers (string-keyed map of strings).
* `body` - Message payload (JSON request body).  Max one per port struct.
* `messageId` - Unique id of the message (typically a random uuid)

`[=<peerName>]`
: (Optional)
    A _peer_ is the index within the field group of the data for each port field in the original message.
    Currently, only `header` fieldGroup has indexed content (individual header values).
    When not specified, the default peer in the metadata is the lowerCamelCase inflection of the field name: 
    eg the `EventType` struct field points to the header `eventType`. 

## Output Ports

Output ports specify parts of the published message to be populated from the
port struct.  For example:

```go
type driftCheckRequestOutput struct {
    EventType string                `out:"header" const:"DriftCheck"`
    Payload   api.DriftCheckRequest `out:"body"`
}
```

In Java this would be implemented using a hypothetical function signature:
```java
void PublishDriftCheckRequest(
  @Header(name="eventType") @Schema(constant="DriftCheck") String eventType, 
  @Body DriftCheckRequest body
)
```

### Struct Tags

Each field with an `out` struct tag will be automatically applied to the message
before the message is published.  The full syntax of the `out` struct tag is as
follows:

    out:"<fieldGroup>[=peerName]"

The subcomponents of the `out` struct tag are the same as `in` struct tag, above.

## Data Transfer Objects (DTOs)

Fields in a port specifying the `body` component will typically have a DTO struct
as their underlying type (eg. `api.DriftCheckRequest` above).  

By default, these are serialized using
the Content-Type of the stream (currently defaults to `application/json`).

