package types

import "encoding/base64"

type Base64Bytes []byte

func (b *Base64Bytes) UnmarshalText(data string) error {
	decodedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	*b = decodedBytes
	return nil
}

func (b *Base64Bytes) MarshalText() (string, error) {
	return base64.StdEncoding.EncodeToString(*b), nil
}

type Binary []byte

func (b *Binary) UnmarshalText(data string) error {
	*b = []byte(data)
	return nil
}

func (b *Binary) MarshalText() (string, error) {
	return string(*b), nil
}

func NewBinary(value []byte) Binary {
	return value
}

func NewBinaryFromString(value string) Binary {
	return []byte(value)
}

func NewBinaryPtr(value []byte) *Binary {
	v := Binary(value)
	return &v
}

func NewBinaryPtrFromString(value string) *Binary {
	return NewBinaryPtr([]byte(value))
}

type Unicode []rune

func (b *Unicode) UnmarshalText(data string) error {
	*b = []rune(data)
	return nil
}

func (b *Unicode) MarshalText() (string, error) {
	return string(*b), nil
}

func NewUnicode(value []rune) Unicode {
	return value
}

func NewUnicodeFromString(value string) Unicode {
	return []rune(value)
}

func NewUnicodePtr(value []rune) *Unicode {
	v := Unicode(value)
	return &v
}

func NewUnicodePtrFromString(value string) *Unicode {
	return NewUnicodePtr([]rune(value))
}
