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
