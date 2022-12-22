// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime"
)

type Encoding []string

func (c Encoding) Writer(target io.WriteCloser) (writer io.WriteCloser, err error) {
	writer = target
	for _, encoding := range c {
		var encoder Encoder
		if encoder, err = NewEncoder(encoding); err != nil {
			return nil, err
		}
		if writer, err = encoder.Writer(writer); err != nil {
			return nil, err
		}
	}
	return writer, nil
}

func (c Encoding) Reader(source io.ReadCloser) (reader io.ReadCloser, err error) {
	reader = source
	for _, encoding := range c {
		var encoder Encoder
		if encoder, err = NewEncoder(encoding); err != nil {
			return nil, err
		}
		if reader, err = encoder.Reader(reader); err != nil {
			return nil, err
		}
	}
	return reader, nil
}

var ErrUnknownEncoder = errors.New("Unknown encoder")

type Encoder interface {
	Reader(io.ReadCloser) (io.ReadCloser, error)
	Writer(io.Writer) (io.WriteCloser, error)
}

type GzipEncoder struct{}

func (g GzipEncoder) Reader(r io.ReadCloser) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

func (g GzipEncoder) Writer(w io.Writer) (io.WriteCloser, error) {
	return gzip.NewWriterLevel(w, gzip.BestCompression)
}

type Base64Encoder struct{}

func (g Base64Encoder) Reader(r io.ReadCloser) (io.ReadCloser, error) {
	return io.NopCloser(base64.NewDecoder(base64.StdEncoding, r)), nil
}

func (g Base64Encoder) Writer(w io.Writer) (io.WriteCloser, error) {
	return base64.NewEncoder(base64.StdEncoding, w), nil
}

const (
	EncoderGzip   = "gzip"
	EncoderBase64 = "base64"
)

var encoders = map[string]Encoder{
	EncoderGzip:   GzipEncoder{},
	EncoderBase64: Base64Encoder{},
}

func RegisterEncoder(name string, m Encoder) {
	encoders[name] = m
}

func NewEncoder(encoder string) (Encoder, error) {
	m, ok := encoders[encoder]

	if !ok {
		return nil, errors.Wrap(ErrUnknownEncoder, encoder)
	}

	return m, nil
}

var ErrUnknownMarshaler = errors.New("Unknown marshaler")

type Marshaler interface {
	WriteEntity(w io.Writer, value interface{}) error
	ReadEntity(r io.ReadCloser, target interface{}) error
}

type JsonMarshaler struct{}

func (m JsonMarshaler) ReadEntity(r io.ReadCloser, target interface{}) error {
	return json.NewDecoder(r).Decode(target)
}

func (m JsonMarshaler) WriteEntity(w io.Writer, value interface{}) error {
	return json.NewEncoder(w).Encode(value)
}

type XmlMarshaler struct{}

func (m XmlMarshaler) ReadEntity(r io.ReadCloser, target interface{}) error {
	return xml.NewDecoder(r).Decode(target)
}

func (m XmlMarshaler) WriteEntity(w io.Writer, value interface{}) error {
	return xml.NewEncoder(w).Encode(value)
}

type BinaryMarshaler struct{}

func (b BinaryMarshaler) WriteEntity(w io.Writer, value interface{}) (err error) {
	switch tv := value.(type) {
	case []byte:
		_, err = w.Write(tv)
	case []rune:
		_, err = w.Write([]byte(string(tv)))
	case string:
		_, err = w.Write([]byte(tv))
	case io.Reader:
		_, err = io.Copy(w, tv)
	default:
		err = errors.Errorf("Could not encode %T to binary", tv)
	}
	return
}

func (b BinaryMarshaler) ReadEntity(r io.ReadCloser, target interface{}) (err error) {
	var tp []byte
	switch tv := target.(type) {
	case *[]byte:
		tp, err = io.ReadAll(r)
		if err == nil {
			*tv = tp
		}
	case *[]rune:
		tp, err = io.ReadAll(r)
		if err == nil {
			*tv = []rune(string(tp))
		}
	case *string:
		tp, err = io.ReadAll(r)
		if err == nil {
			*tv = string(tp)
		}
	case io.Writer:
		_, err = io.Copy(tv, r)
	default:
		err = errors.Errorf("Could not decode %T from binary", tv)
	}
	return
}

const (
	MarshalerJson   = "application/json"
	MarshalerXml    = "application/xml"
	MarshalerBinary = "application/octet-stream"
)

var marshalers = map[string]Marshaler{
	MarshalerJson:   JsonMarshaler{},
	MarshalerXml:    XmlMarshaler{},
	MarshalerBinary: BinaryMarshaler{},
}

func RegisterMarshaler(name string, m Marshaler) {
	marshalers[name] = m
}

func NewMarshaler(marshaler string) (Marshaler, error) {
	m, ok := marshalers[marshaler]

	if !ok {
		return nil, errors.Wrap(ErrUnknownMarshaler, marshaler)
	}

	return m, nil
}

type ContentOptions struct {
	MimeType string
	Encoding
}

func (c *ContentOptions) WithEncoding(encodings ...string) *ContentOptions {
	c.Encoding = append(c.Encoding, encodings...)
	return c
}

func (c ContentOptions) WriteEntity(w io.WriteCloser, value interface{}) (err error) {
	w, err = c.Encoding.Writer(w)
	if err != nil {
		return
	}
	defer func() {
		err = w.Close()
	}()

	mimeType, err := c.BaseMediaType()
	if err != nil {
		return err
	}

	m, err := NewMarshaler(mimeType)
	if err != nil {
		return
	}

	return m.WriteEntity(w, value)
}

func (c ContentOptions) ReadEntity(r io.ReadCloser, target interface{}) (err error) {
	r, err = c.Encoding.Reader(r)
	if err != nil {
		return err
	}

	mimeType, err := c.BaseMediaType()
	if err != nil {
		return err
	}

	m, err := NewMarshaler(mimeType)
	if err != nil {
		return
	}

	return m.ReadEntity(r, target)
}

func (c ContentOptions) ReadBytes(r io.ReadCloser) (data []byte, err error) {
	r, err = c.Encoding.Reader(r)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}

func (c *ContentOptions) BaseMediaType() (string, error) {
	mimeType, _, err := mime.ParseMediaType(c.MimeType)
	return mimeType, err
}

func NewContentOptions(mimeType string) ContentOptions {
	return ContentOptions{
		MimeType: mimeType,
	}
}

type Content struct {
	present bool
	options ContentOptions
	source  io.ReadCloser
}

func (s Content) IsPresent() bool {
	return s.present
}

func (s Content) BaseMediaType() (string, error) {
	return s.options.BaseMediaType()
}

func (s Content) MimeType() string {
	return s.options.MimeType
}

func (s Content) ReadEntity(target interface{}) (err error) {
	return s.options.ReadEntity(s.source, target)
}

func (s Content) ReadBytes() (data []byte, err error) {
	return s.options.ReadBytes(s.source)
}

func (s Content) Reader() (r io.ReadCloser, err error) {
	if !s.present {
		return nil, errors.Wrap(io.ErrUnexpectedEOF, "Content source not present")
	}
	return s.options.Reader(s.source)
}

func NewContentFromBytes(content ContentOptions, source []byte) Content {
	return Content{
		present: source != nil,
		options: content,
		source:  io.NopCloser(bytes.NewBuffer(source)),
	}
}

func NewContentFromReadCloser(content ContentOptions, source io.ReadCloser) Content {
	return Content{
		present: source != nil,
		options: content,
		source:  source,
	}
}
