package types

type TextUnmarshaler interface {
	UnmarshalText(data string) error
}

type TextMarshaler interface {
	MarshalText() (string, error)
}
