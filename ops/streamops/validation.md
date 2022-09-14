# Validation

Stream operations use struct tags to declare JSON Schema
constraints for fields within the port and data transfer objects.

```go
type driftCheckResponseInput struct {
    EventType string                 `in:"header" const:"DriftCheck"`
    Payload   api.DriftCheckResponse `in:"body"`
}
```

For example, in the `driftCheckRequestInput` port, the `EventType` field
specifies it must contain the constant value `DriftCheck` through the `const` tag.

When generating AsyncApi documentation, validation constraints specified in the port
struct or the data transfer object will automatically be included in the documentation.

## JSON Schema Struct Tags

[JSON Schema](https://json-schema.org/) is a domain-specific language used to describe constraints on values
that may be expressed using the JSON type system (null, object, array, number, string).

The following struct tags may be used to specify JSON Schema constraints and validators
within a port struct or a data transfer object:

* `optional`, `required`

  Boolean entries to override the base field type optionality.  By default, pointer/slice/map and `types.Optional`
  types are optional, and other types are non-optional.  Must be "true" or "false".

    ```go
    type inPort struct {
        Expiry     types.Duration  `in="header" optional="true"`
        BestBefore *types.Duration `in="header" required="true"`
    }
    ```

* `deprecated`

  Boolean to indicate the field should not be used and will be removed. Must be "true" or "false".

    ```go
    type inPort struct {
        VmsTenant types.UUID `in="header" deprecated="true"`  
    }
    ```

* `title`, `description`

  Expository language to help users understand the purpose of the field.

    ```go
    type outPort struct {
        ContentType string `in="header" title="Content Type" description="MIME type for the message body" example="application/xml"`
    }
    ```

* `const`

  Value that the field must contain to be valid.  Equivalent to a single-valued `enum`.
  Must be a scalar convertible to a valid value of the field data type.
  See "Port Field Constraint Handling", below for more details.

    ```go
    type MyRequest struct {
        Answer int `in="header" const="42"`  
    }
    ```

* `default`

  Default value that the field will _behave as having_ if not explicitly specified.
  Must be a scalar convertible to a valid value of the field data type.
  See "Port Field Constraint Handling", below for more details.

    ```go
    type MyRequest struct {
        Pi float64 `in="header" default="3.14"`  
    }
    ```

* `example`

  Example value to be presented in the schema document.
  Must be a scalar convertible to a valid value of the field data type.

    ```go
    type MyRequest struct {
        Hour int   `in="header" example="12"`
        Minute int `in="header" example="30"`  
    }
    ```

* `enum`

  Comma-separated list of possible values for the field.
  Only these values will be accepted by the field during validation.
  Must be scalars convertible to the field data type.

    ```go
    type MyResponseOutput struct {
        Code int `out="code" enum="200,400,401,403,404"` 
    }
    ```

* `minimum`, `maximum`

  Range constraints for possible values of the field.
  Only values >= minimum (if specified) are valid.
  Only values <= maximum (if specified) are valid.
  Must be scalars convertible to the field data type.

    ```go
    type MyResponseOutput struct {
        Radians float64 `out="header" minimum="0" maximum="6.28"` 
    }
    ```

* `minLength`, `maxLength`

  Length constraints for value of the field.  Applies
  to `string` fields.  Must be integers if specified.

  ```go
  type MyResponse struct {
       ServiceType string `in="header" minLength="4" maxLength="16"`
  }
  ```

* `maxProperties`, `minProperties`

  Length constraints for value of the field.  Applies
  to `object` fields.  Must be integers if specified.

  ```go
  type MyResponse struct {
       ServiceType map[string]string `in="header" minProperties="1"`
  }
  ```

* `pattern`

  Regular expression that values must match to be valid.  Applies
  to `string` fields.

  ```go
  type MyResponse struct {
       DeviceId string `in="header" pattern="^CPE-.*$"`
  }
  ```

* `format`

  String identifier of pre-defined formats.  Applies
  to `string` fields.  Normally will be automatic based on
  the underlying field type.

  ```go
  type MyResponse struct {
       When types.Time `in="header" format="date"`
  }
  ```

* `minItems`, `maxItems`

  Length constraints for value of the field. Applies
  to `array` fields (slices).  Must be integers if specified.

  ```go
  type MyRequest struct {
       DeviceIds []types.UUID `in="body" minItems="1"`
       TenantIds []types.UUID `in="body" maxItems="1"`
  }
  ```

The underlying jsonschema-go library provides a few more constraints,
which you can view at the package [GoDoc](https://pkg.go.dev/github.com/swaggest/jsonschema-go#readme-field-tags)

### Constraints on Named Types

Fields in Port structures and DTOs with simple and anonymous types may be augmented using the JSON schema
tags above.  However, named types are shared across many fields and therefore cannot be augmented in-place.
For example:

```go
type DriftCheckRequest struct {
  Action                string          `json:"action" const:"checkDrift"`
  GroupId               types.UUID      `json:"groupId,omitempty"`
  Timestamp             types.Time      `json:"timestamp" minimum:"2022-01-01T00:00:00Z"`
  EntityLevelCompliance string          `json:"entityLevelCompliance" enum:"full,partial"`
  Standards             []ConfigPayload `json:"standards,omitempty" minItems:"1" required:"true"`
  ...
}
```

Fields with Named Types include:
- `GroupId` : `types.UUID`
- `Timestamp` : `types.Time`

These fields will ignore any schema constraints declared in the struct tag,
such as the `minimum` tag on `Timestamp`.

Fields with Simple or Anonymous Types include:
- `Action` : `string`
- `EntityLevelCompliance` : `string`
- `Standards`: `[]ConfigPayload`

Each of these fields has schema constraints declared which will be honoured.
`Standards` is an array of DTOs `[]ConfigPayload` and therefore is of an anonymous
type.

### Constraints on DTO structs

To configure a parent DTO struct using struct tags, include an anonymous 
field `_` with the desired constraints.  For example:

```go
type RemediateRequest struct {
	...
    _ struct{} `additionalProperties:"false" description:"RemediateRequest contains a remediation request."`
}
```

This will add `description` and `additionalProperties` schema constraints
to the `RemediateRequest` struct in the schema.

### Custom Schema Generation for Named Types

The underlying jsonschema-go library provides a number of interfaces to customize
or replace the JSON schema generated for your Named Type:

- `NamedEnum` - Provides a list of name/value pairs for your enumerable type.
- `Enum` - Provides a list of values for your enumerable type.
- `Preparer` - Intercepts the reflected JSON Schema and allows alteration.
- `Exposer` - Provides a complete parsed JSON Schema for your type.
- `RawExposer` - Provides a complete unparsed JSON Schema for your type.
- `OneOfExposer` - Provides a list of `oneOf` elements for your type.
- `AnyOfExposer` - Provides a list of `anyOf` elements for your type.
- `AllOfExposer` - Provides a list of `allOf` elements for your type.
- `NotExposer` - Provides a `not` element for your type.
- `IfExposer` - Provides an `if` element for your type.
- `ThenExposer` - Provides a `then` element for your type.
- `ElseExposer` - Provides an `else` element for your type.

You can find more details about these interfaces on the package
[GoDoc](https://pkg.go.dev/github.com/swaggest/jsonschema-go#readme-implementing-interfaces-on-a-type).

### Port Field Constraint Handling

To ease development burden, when using `const` or `default` on a Port Field, the value will be applied
during input population (subscriber) or output population (publisher).  Note that this only applies to
scalars (e.g. headers), and only those directly contained in the Port structure.  In particular, it does
not apply to the request/response body or its sub-fields.

From the example at the beginning of the chapter:

  ```go
  type driftCheckRequestOutput struct {
      EventType string            `out:"header" const:"DriftCheck"`
      Payload   DriftCheckRequest `out:"body"`
  }
  
  type DriftCheckRequest struct {
      Action string `json:"action" const:"checkDrift"`
	  ...
  }
  ```

The `EventType` field of `driftCheckRequestOutput` will be filled with `DriftCheck` if not supplied
by the publisher, since it is a scalar, and directly contained within the Port structure.  If another
value is supplied, the schema validation will fail, so it is best to simply not supply the value.

The `Action` field of `DriftCheckRequest` will _not_ be filled with `checkDrift` since it is not
directly contained within the Port structure.  It will be validated during schema validation to ensure
only that value is supplied.
