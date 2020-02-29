package types

type OptionalString struct {
	Value *string
}

func (s OptionalString) IsPresent() bool {
	return s.Value != nil
}

func (s OptionalString) OrOptional(value *string) OptionalString {
	if s.Value == nil {
		return OptionalString{value}
	}
	return s
}

func (s OptionalString) OrElse(value string) string {
	if s.Value != nil {
		return *s.Value
	}
	return value
}

func (s OptionalString) String() string {
	return s.OrElse("<nil>")
}

func NewOptionalString(value *string) OptionalString {
	return OptionalString{Value:value}
}