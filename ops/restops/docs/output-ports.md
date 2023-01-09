# REST Output Ports

As described in [Output Ports](../../docs/output-ports.md), an Output Port is a go structure 
used to describe a destination payload to be generated, such as an HTTP REST API Response.
An Output port is the Response equivalent of an Input Port, and follows many of the same
rules and patterns.

Each field in the Output Port structure is expected to have a `resp` struct tag. Any fields
missing this tag will be ignored by the output populator.

The structure includes any required or optional results, including status code, headers,
paging, and success and/or error bodies.

### Example

The following example shows a simple Create API output port definition:

```go
type createEntityInputs struct {
    Payload        api.CreateEntityResponse `resp:"body"`
}
```

The response body will be populated with a JSON-serialized instance of `api.CreateEntityResponse`.

### Field Groups

The possible field groups used by the `resp` struct tag are:

- **code**: The HTTP status code for the response
- **header**: 
- **paging**: An envelope wrapping the body containing the paging response
- **body**: The primary payload of the response (excluding any envelopes/paging).
  You may also specify `success:"true"` or `error:"true"` to define multiple
  potential bodies.

### Field Index

Each header field will have a group (source) and index (key).
You will recall the format of the `resp` tag:

    resp:"<fieldGroup>[=<fieldIndex>]"

Header field indices default to the Upper-Kebab-Case inflection of the field name 
by default.

Non-indexed fields such as `code` and `body` do not accept a field index, and
they will be ignored if specified; there is only one of each of these in any 
generated response.
