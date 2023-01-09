# Input Ports

An Input Port is a go structure used to describe a source payload to be parsed,
such as an HTTP Request or an incoming Stream Message.  

Upon receipt of an incoming payload, go-msx will populate your data structure, execute 
validation, and if it is valid, pass it to your handler, such as an HTTP Controller Endpoint.

An example stream message input port:
```go
type driftCheckResponseInputs struct {
    EventType string                 `in:"header" const:"DriftCheck"`
    Payload   api.DriftCheckResponse `in:"body"`
}
```

An example HTTP request input port:
```go
type createEntityInputs struct {
    ControlPlaneId types.UUID              `req:"path"`
    Payload        api.CreateEntityRequest `req:"body"`	
}
```

## Struct Tags

Each field with an input struct tag will be automatically populated before being
passed to your handler.

Note that the struct tag prefix depends on the protocol being described:
- For HTTP Requests, the input struct tag must be `req` (for backwards compatibility)
- For Stream Messages, the input struct tag must be `in`.

The full syntax of the input struct tag is one of the following, appropriate for
the handling protocol:

    in:"<fieldGroup>[=<peerName>]"
    req:"<fieldGroup>[=<peerName>]"

The input struct tag contains the following subcomponents:

- `<fieldGroup>` (_Required_):
The name of the message/request part from which the value will be extracted.

- `[=<peerName>]` (_Optional_):
A _peer_ is a field or property in the source message 
For example an HTTP request may have a header with the name date, which can be requested
using the following input struct tag: `req:"header=Date"` 

See [HTTP Request Ports]() and [Stream Ports](../streamops/ports.md) for available field groups and peer
name conventions for your specific protocol.
