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

func (s OptionalString) OrEmpty() string {
	return s.OrElse("")
}

func (s OptionalString) String() string {
	return s.OrElse("<nil>")
}

func (s OptionalString) Ptr() *string {
	return s.Value
}

func (s OptionalString) Equals(other OptionalString) bool {
	if s.Value == other.Value {
		return true
	} else if s.Value == nil || other.Value == nil {
		return false
	} else {
		return *s.Value == *other.Value
	}
}

func NewOptionalString(value *string) OptionalString {
	return OptionalString{Value: value}
}

func NewOptionalStringFromString(value string) OptionalString {
	return OptionalString{Value: &value}
}

func NewStringPtr(value string) *string {
	return &value
}

type Optional struct {
	Value interface{}
}

func (o Optional) IfPresent(fn func(v interface{})) Optional {
	if o.Value != nil {
		fn(o.Value)
	}
	return o
}

func (o Optional) IfNotPresent(fn func()) Optional {
	if o.Value == nil {
		fn()
	}
	return o
}

func (o Optional) OrElse(v interface{}) interface{} {
	if o.Value != nil {
		return o.Value
	}
	return v
}

func NewOptional(value interface{}) Optional {
	return Optional{Value: value}
}
