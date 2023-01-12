# Ports

_Ports_ describe the format of an incoming or outgoing stream or HTTP endpoint.  
They are used to automatically serialize and deserialize communications to a convenient
structure for your application.

The term "Ports" is taken from Hexagonal Architecture, where it used to describe 

> ... dedicated interfaces to communicate with the outside world. 
> They allow the entry or exiting of data to and from the application.
>
> -- [Hexagonal Architecture](https://medium.com/idealo-tech-blog/hexagonal-ports-adapters-architecture-e3617bcf00a0), Medium

go-msx uses two types of ports, with many overlapping options:

- **Input Port**: describes an HTTP request or incoming Stream Message
- **Output Port**: describes an HTTP response or outgoing Stream Message

### Declaration

Ports are defined using go structures, consisting of a series of fields.
Each field consists of three parts:

1. **Name**: The name by which you can access the struct member in go code.
2. **Type**: The type of the field to which the data will be converted.
   These types fall into one of a few categories, to simplify conversion:

    - Scalar: Any simple single-valued (eg string, int, uuid, bool)
    - Array: A sequence of scalars
    - Object: A dictionary of scalars with string keys
    - File: An uploaded file
    - FileArray: A sequence of uploaded files
3. **Tags**: A set of annotations of the struct field, describing attributes
   like source/destination, index, validation, optionality, etc. 

### Example

```go
type outputs struct {
    Code   int          `resp:"code"`
    Body   api.Response `resp:"body"`
    Error  api.Error    `resp:"body" error:"true"`
}
```
