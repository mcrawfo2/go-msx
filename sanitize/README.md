# Sanitize

MSX Sanitize allows request data to be pre-processed (before validation) to ensure potentially dangerous content is
removed. For example XSS and arbitrary HTML can be removed from plain-text strings.

## Usage

To sanitize a tree of data, including maps, slices, structs in-place:

```go
if err := sanitize.Input(&mydata, sanitize.NewOptions("xss")); err != nil {
	return err
}
```

After returning, mydata will be sanitized based on the supplied Options.

### Options

Options are available for each of the sanitizers from

    github.com/kennygrant/sanitize

including:
- Accents 
- BaseName 
- Xss 
- Name 
- Path

### Struct Tags

To specify these options on a struct field, use the `san:"..."` tag, for example:

```go
type MyRequest struct {
	Name 		string `json:"name" san:"xss"`
	Description string `json:"description" san:"xss"`
	Ignored 	string `json:"ignored" san:"-"`
}
```

In this struct, `Name` and `Description` fields indicate they must be sanitized for XSS/HTML content (`xss`),
and `Ignored` should not be sanitized at all (`-`).

NOTE: If a struct field does not have the `san` tag, it will inherit from its ancestors, up to the options passed
into the `sanitize.Input` call.
