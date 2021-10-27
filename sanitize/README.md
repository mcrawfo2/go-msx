# Sanitize

MSX Sanitize allows request data to be pre-processed (before validation) to ensure potentially dangerous content is
removed. For example XSS and arbitrary HTML can be removed from plain-text strings.  MSX Sanitize also auto-sanitizes
log messages.  

## Sanitizing Input

To explicitly sanitize a tree of data, including maps, slices, structs in-place:

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
- Accents (`accents`) 
- BaseName (`basename`)
- Xss (`xss`) 
- Name (`name`)
- Path (`path`)

Custom sanitizers provided by MSX Sanitize include:
- Secret (`secret`)

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

## Sanitizing Logs

Logs are auto-sanitized using some base rules.  These can be augmented by the microservice using the 
`sanitize.secrets` configuration:

```yaml
sanitize.secrets:
  keys:
    - status
  custom:
    enabled: true
    patterns:
        - from: "\\[userviceconfiguration/\\w+\\]"
          to: "[userviceconfiguration/...]"
        - from: "\\[secret/\\w+\\]"
          to: "[secret/...]"
```

Within `sanitize.secrets` you can configure the following options:

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `enabled` | true | Optional | Enable secret replacement |
| `keys` | - | Optional | A set of XML/JSON/ToString attributes and objects to flag as sensitive |
| `custom.*` | - | Optional | Custom go regex replacement.  Does not use `keys`. | 
| `json.*` | - | Optional | JSON replacement.  Replaces once per entry in `keys`. | 
| `xml.*` | - | Optional | XML replacement.  Replaces once per entry in `keys`. |
| `to-string.*` | - | Optional | Stringer replacement.  Replaces once per entry in `keys`. |

For `custom`, specify a list of regexes and replacements in `custom.patterns`, as above.

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `custom.patterns[*].from` | - | Required | Regex to match |
| `custom.patterns[*].to` | - | Required | Replacement (including variables) |

For `json`, `xml`, `tostring`, specify a list of regexes to match, including the named capture groups
`prefix` and `postfix`:

| Key                             | Default   | Required | Description |
|---------------------------------|-----------|----------|-------------|
| `.enabled` | true | Optional | Enable this set of patterns (`json`, `xml`, `to-string`) |
| `.patterns[*].from` | - | Required | Regex to match |
| `.patterns[*].to` | `${prefix}*****${postfix}` | Optional | Replacement (including regex variables) |

