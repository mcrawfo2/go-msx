# REST Input Ports

As described in [Input Ports](../../docs/input-ports.md), an Input Port is a go structure used to describe
a source payload to be parsed, such as an HTTP REST API Request.

Each field in the Input Port structure is expected to have a `req` struct tag. Any fields
missing this tag will be ignored by the input populator.

The structure includes any required or optional parameters (Cookies,
Form, Headers, Path, Query), along with any expected body content.

### Example

The following example shows a simple Create API input port definition:

```go
type createEntityInputs struct {
    ControlPlaneId types.UUID              `req:"path"`
    Payload        api.CreateEntityRequest `req:"body"`
}
```

In this example, `ControlPlaneId` is expected to be found in the path (with the default
path parameter style, `controlPlaneId`). The body is expected to contain a JSON-serialized
instance of `api.CreateEntityRequest`.  

### Field Groups

The possible field groups used by the `req` struct tag are:

- **method**: The HTTP method
- **header**: An HTTP header
- **cookie**: An sub-entry from the Cookie header 
- **path**: A segment of the path
- **query**: A query parameter
- **form**: A form field
- **body**: The body content

### Field Index

Each field will typically have a group (source) and index (key).
You will recall the format of the `req` tag:

    req:"<fieldGroup>[=<fieldIndex>]"

Most field indices default to the lowerCamelCase inflection of the field name.
The only exception is for headers, which are the Upper-Kebab-Case inflection
of the field name by default.

Non-indexed fields such as `method` and `body` do not accept a field index, and
they will be ignored if specified; there is only one of each of these in any request.
