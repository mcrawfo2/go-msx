# Ports

_Ports_ describe the various components of an incoming or outgoing stream or HTTP endpoint.  
They are used to automatically serialize and deserialize communications to a convenient
structure for your application.

The term "Ports" is taken from Hexagonal Architecture, where it used to describe an
entry point (interface) into the user application ("core logic").

go-msx uses two types of ports, with many overlapping options:
- **Input Port**: describes an HTTP request or incoming Stream Message
- **Output Port**: describes an HTTP response or outgoing Stream Message

Input and Output ports are described in the following sections.
