// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type TestWriter struct {
	w io.Writer
}

func (t TestWriter) Close() error {
	return nil
}

func (t TestWriter) Write(p []byte) (n int, err error) {
	if t.w == nil {
		return 0, errors.New("no writer")
	}

	c := make([]byte, len(p))
	copy(c, p)
	lc := len(c)
	for i := 0; i < lc/2; i++ {
		j := lc - i - 1
		c[i], c[j] = c[j], c[i]
	}
	return t.w.Write(c)
}

type TestReader struct {
	r io.ReadCloser
}

func (t TestReader) Read(p []byte) (n int, err error) {
	if t.r == nil {
		return 0, errors.New("no reader")
	}

	c := make([]byte, len(p))
	_, err = t.r.Read(c)
	if err != nil {
		return
	}

	lc := len(c)
	for i := 0; i < lc/2; i++ {
		j := lc - i - 1
		c[i], c[j] = c[j], c[i]
	}
	copy(p, c)
	return
}

func (t TestReader) Close() error {
	return t.r.Close()
}

type TestEncoder struct{}

func (e TestEncoder) Reader(r io.ReadCloser) (io.ReadCloser, error) {
	if r == nil {
		return nil, errors.New("no reader")
	}

	return TestReader{r: r}, nil
}

func (e TestEncoder) Writer(w io.Writer) (io.WriteCloser, error) {
	if w == nil {
		return nil, errors.New("no writer")
	}

	return TestWriter{w: w}, nil
}

func TestNewEncoder(t *testing.T) {
	RegisterEncoder("test", TestEncoder{})

	type args struct {
		encoder string
	}
	tests := []struct {
		name    string
		args    args
		want    Encoder
		wantErr bool
	}{
		{
			name: "Test",
			args: args{
				encoder: "test",
			},
			want:    encoders["test"],
			wantErr: false,
		},
		{
			name: "Unknown",
			args: args{
				encoder: "unknown",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEncoder(tt.args.encoder)
			if !(tt.wantErr == (err != nil)) {
				assert.Failf(t, "NewEncoder(%v)", tt.args.encoder)
			}
			if !tt.wantErr {
				assert.Equalf(t, tt.want, got, "%s", testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestEncoding_Writer(t *testing.T) {
	RegisterEncoder("test", TestEncoder{})

	tests := []struct {
		name         string
		c            Encoding
		source       string
		noTarget     bool
		wantTarget   string
		wantWriter   io.Writer
		wantWriteErr bool
		wantErr      bool
	}{
		{
			name:       "Empty",
			c:          Encoding{},
			source:     "ABC",
			wantTarget: "ABC",
			wantWriter: types.CloseableByteBuffer{Buffer: &bytes.Buffer{}},
			wantErr:    false,
		},
		{
			name:       "Test",
			c:          Encoding{"test"},
			source:     "ABC",
			wantTarget: "CBA",
			wantWriter: TestWriter{w: types.CloseableByteBuffer{Buffer: &bytes.Buffer{}}},
			wantErr:    false,
		},
		{
			name:       "DoubleTest",
			c:          Encoding{"test", "test"},
			source:     "ABC",
			wantTarget: "ABC",
			wantWriter: TestWriter{w: TestWriter{w: types.CloseableByteBuffer{Buffer: &bytes.Buffer{}}}},
			wantErr:    false,
		},
		{
			name:    "MissingEncoder",
			c:       Encoding{"test2"},
			wantErr: true,
		},
		{
			name:       "NoWriter",
			c:          Encoding{"test"},
			noTarget:   true,
			wantWriter: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target io.WriteCloser
			var buffer = new(bytes.Buffer)
			if !tt.noTarget {
				target = types.CloseableByteBuffer{Buffer: buffer}
			}

			gotWriter, err := tt.c.Writer(target)

			if !(tt.wantErr == (err != nil)) {
				assert.Failf(t, "Unexpected Writer() error", "error = %v", err)
			}
			if tt.wantErr == false {
				assert.Equalf(t, tt.wantWriter, gotWriter, "Unexpected writer: %v", target)

				_, err = gotWriter.Write([]byte(tt.source))
				if !(tt.wantWriteErr == (err != nil)) {
					assert.Failf(t, "Unexpected Write() error", "error = %v", err)
				}
				if !tt.wantWriteErr && !tt.noTarget {
					assert.Equalf(t, tt.wantTarget, buffer.String(), "Writer(%v)", target)
				}
			}
		})
	}
}

func TestEncoding_Reader(t *testing.T) {
	RegisterEncoder("test", TestEncoder{})

	tests := []struct {
		name        string
		c           Encoding
		source      string
		noSource    bool
		wantTarget  string
		wantReader  io.ReadCloser
		wantErr     bool
		wantReadErr bool
	}{

		{
			name:       "Empty",
			c:          Encoding{},
			source:     "ABC",
			wantTarget: "ABC",
			wantReader: io.NopCloser(bytes.NewBufferString("ABC")),
			wantErr:    false,
		},
		{
			name:       "Test",
			c:          Encoding{"test"},
			source:     "ABC",
			wantTarget: "CBA",
			wantReader: TestReader{r: io.NopCloser(bytes.NewBufferString("ABC"))},
			wantErr:    false,
		},
		{
			name:       "DoubleTest",
			c:          Encoding{"test", "test"},
			source:     "ABC",
			wantTarget: "ABC",
			wantReader: TestReader{r: TestReader{r: io.NopCloser(bytes.NewBufferString("ABC"))}},
			wantErr:    false,
		},
		{
			name:    "MissingEncoder",
			c:       Encoding{"test2"},
			wantErr: true,
		},
		{
			name:       "NoWriter",
			c:          Encoding{"test"},
			noSource:   true,
			wantReader: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var source io.ReadCloser
			var buffer = new(bytes.Buffer)
			if !tt.noSource {
				buffer.Write([]byte(tt.source))
				source = io.NopCloser(buffer)
			}

			gotReader, err := tt.c.Reader(source)
			if !(tt.wantErr == (err != nil)) {
				assert.Failf(t, "Unexpected Reader() error", "error = %v", err)
			}
			if tt.wantErr == false {
				assert.Equalf(t, tt.wantReader, gotReader, "Unexpected reader: %v", source)

				target := make([]byte, buffer.Len())
				_, err = gotReader.Read(target)
				if !(tt.wantReadErr == (err != nil)) {
					assert.Failf(t, "Unexpected Read() error", "error = %v", err)
				}
				if !tt.wantReadErr && !tt.noSource {
					assert.Equalf(t, tt.wantTarget, string(target), "%s", testhelpers.Diff(tt.wantTarget, string(target)))
				}
			}
		})
	}
}
