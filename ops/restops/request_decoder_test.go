// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"reflect"
	"testing"
)

func TestEndpointRequestDecoder_DecodePrimitive(t *testing.T) {
	type fields struct {
		dataSource MockRequestDataSource
	}
	type args struct {
		in      string
		name    string
		style   string
		explode bool
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult types.Optional[string]
		wantErr    bool
	}{
		{
			name: "HeaderSimple",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"Example": {"blue"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantResult: types.OptionalOf("blue"),
			wantErr:    false,
		},
		{
			name: "HeaderSimpleExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"Example": {"blue"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "simple",
				explode: true,
			},
			wantResult: types.OptionalOf("blue"),
			wantErr:    false,
		},
		{
			name: "HeaderInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"example": {"blue", "black", "brown"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantErr: true,
		},
		{
			name: "PathSimple",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "blue",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "simple",
				explode: true,
			},
			wantResult: types.OptionalOf("blue"),
		},
		{
			name: "PathMatrixUnsupported",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": ";color=blue,black,brown",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "matrix",
				explode: true,
			},
			wantErr: true,
		},
		{
			name: "PathLabelUnsupported",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": ".blue.black.brown",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "label",
				explode: true,
			},
			wantErr: true,
		},
		{
			name: "PathInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "blue,black,brown",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantErr: true,
		},
		{
			name: "QueryForm",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "form",
				explode: false,
			},
			wantResult: types.OptionalOf("blue"),
		},
		{
			name: "QueryFormExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue", "black", "brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantResult: types.OptionalOf("blue"),
		},
		{
			name: "QuerySpaceDelimited",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue black brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "spaceDelimited",
				explode: false,
			},
			wantResult: types.OptionalOf("blue black brown"),
		},
		{
			name: "QueryPipeDelimited",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue|black|brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "pipeDelimited",
				explode: false,
			},
			wantResult: types.OptionalOf("blue|black|brown"),
		},
		{
			name: "QueryInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantErr: true,
		},
		{
			name: "CookieForm",
			fields: fields{
				dataSource: MockRequestDataSource{
					cookies: []*http.Cookie{
						{Name: "example", Value: "blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "cookie",
				name:    "example",
				style:   "form",
				explode: false,
			},
			wantResult: types.OptionalOf("blue,black,brown"),
		},
		{
			name: "CookieFormExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					cookies: []*http.Cookie{
						{Name: "example", Value: "blue"},
						{Name: "example", Value: "black"},
						{Name: "example", Value: "brown"},
					},
				},
			},
			args: args{
				in:      "cookie",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantResult: types.OptionalOf("blue"),
		},
		{
			name: "CookieInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "cookie",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantErr: true,
		},
		{
			name: "Form",
			fields: fields{
				dataSource: MockRequestDataSource{
					form: map[string][]string{
						"example": {"blue", "black", "brown"},
					},
				},
			},
			args: args{
				in:      "form",
				name:    "example",
				style:   "ignored",
				explode: true,
			},
			wantResult: types.OptionalOf("blue"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &EndpointRequestDecoder{
				DataSource: tt.fields.dataSource,
			}

			pf := ops.PortField{
				Group: tt.args.in,
				Peer:  tt.args.name,
				Options: map[string]string{
					"style":   tt.args.style,
					"explode": fmt.Sprintf("%v", tt.args.explode),
				},
			}

			gotResult, err := d.DecodePrimitive(&pf)
			if tt.wantErr != (err != nil) {
				if err == nil {
					assert.Errorf(t, err, "DecodePrimitive(%v, %v, %v, %v)", tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
				} else {
					assert.NoErrorf(t, err, "DecodePrimitive(%v, %v, %v, %v)", tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
				}
			}
			assert.Equalf(t, tt.wantResult, gotResult, "DecodePrimitive(%v, %v, %v, %v)", tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
		})
	}
}

func TestEndpointRequestDecoder_DecodeArray(t *testing.T) {
	type fields struct {
		dataSource MockRequestDataSource
	}
	type args struct {
		in      string
		name    string
		style   string
		explode bool
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult []string
		wantErr    bool
	}{
		{
			name: "HeaderSimple",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"Example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantResult: []string{"blue", "black", "brown"},
			wantErr:    false,
		},
		{
			name: "HeaderSimpleExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"Example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "simple",
				explode: true,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
		{
			name: "HeaderInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"Example": {"blue", "black", "brown"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "PathSimple",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "blue,black,brown",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "simple",
				explode: true,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
		{
			name: "PathMatrixUnsupported",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "blue,black,brown",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "matrix",
				explode: true,
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "PathLabelUnsupported",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "blue,black,brown",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "label",
				explode: true,
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "PathInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "blue,black,brown",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "QueryForm",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "form",
				explode: false,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
		{
			name: "QueryFormExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue", "black", "brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
		{
			name: "QuerySpaceDelimited",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue black brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "spaceDelimited",
				explode: false,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
		{
			name: "QueryPipeDelimited",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue|black|brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "pipeDelimited",
				explode: false,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
		{
			name: "QueryInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "CookieForm",
			fields: fields{
				dataSource: MockRequestDataSource{
					cookies: []*http.Cookie{
						{Name: "example", Value: "blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "cookie",
				name:    "example",
				style:   "form",
				explode: false,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
		{
			name: "CookieFormExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					cookies: []*http.Cookie{
						{Name: "example", Value: "blue"},
						{Name: "example", Value: "black"},
						{Name: "example", Value: "brown"},
					},
				},
			},
			args: args{
				in:      "cookie",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
		{
			name: "CookieInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "cookie",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "FormExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					form: map[string][]string{
						"example": {"blue", "black", "brown"},
					},
				},
			},
			args: args{
				in:      "form",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
		{
			name: "FormNoExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					form: map[string][]string{
						"example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "form",
				name:    "example",
				style:   "form",
				explode: false,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &EndpointRequestDecoder{
				DataSource: tt.fields.dataSource,
			}

			pf := ops.PortField{
				Group: tt.args.in,
				Peer:  tt.args.name,
				Options: map[string]string{
					"style":   tt.args.style,
					"explode": fmt.Sprintf("%v", tt.args.explode),
				},
			}

			gotResult, err := d.DecodeArray(&pf)
			if tt.wantErr != (err != nil) {
				if err == nil {
					assert.Errorf(t, err, "DecodeArray(%v, %v, %v, %v)", tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
				} else {
					assert.NoErrorf(t, err, "DecodeArray(%v, %v, %v, %v)", tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
				}
			}
			assert.Equalf(t, tt.wantResult, gotResult, "DecodeArray(%v, %v, %v, %v)", tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
		})
	}
}

func TestEndpointRequestDecoder_DecodeObject(t *testing.T) {
	type fields struct {
		dataSource MockRequestDataSource
	}
	type args struct {
		in      string
		name    string
		style   string
		explode bool
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult types.Pojo
		wantErr    bool
	}{
		{
			name: "HeaderSimple",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"Example": {"R,100,G,200,B,150"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
		{
			name: "HeaderSimpleExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"Example": {"R=100,G=200,B=150"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "simple",
				explode: true,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
			wantErr:    false,
		},
		{
			name: "HeaderInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"example": {"blue", "black", "brown"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantErr: true,
		},
		{
			name: "PathSimple",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "R,100,G,200,B,150",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
		{
			name: "PathSimpleExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "R=100,G=200,B=150",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "simple",
				explode: true,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
		{
			name: "PathMatrixUnsupported",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "R=100,G=200,B=150",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "matrix",
				explode: true,
			},
			wantErr: true,
		},
		{
			name: "PathLabelUnsupported",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": ".blue.black.brown",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "label",
				explode: true,
			},
			wantErr: true,
		},
		{
			name: "PathInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					pathParameters: map[string]string{
						"example": "blue,black,brown",
					},
				},
			},
			args: args{
				in:      "path",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantErr: true,
		},
		{
			name: "QueryForm",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"R,100,G,200,B,150"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "form",
				explode: false,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
		{
			name: "QueryFormExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"R": {"100"},
						"G": {"200"},
						"B": {"150"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
		{
			name: "QuerySpaceDelimited",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"R 100 G 200 B 150"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "spaceDelimited",
				explode: false,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
		{
			name: "QueryPipeDelimited",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"R|100|G|200|B|150"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "pipeDelimited",
				explode: false,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
		{
			name: "QueryDeepObject",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example[R]": {"100"},
						"example[G]": {"200"},
						"example[B]": {"150"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "deepObject",
				explode: true,
			},
			wantResult: types.Pojo{
				"R": "100",
				"G": "200",
				"B": "150",
			},
		},
		{
			name: "QueryInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "query",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantErr: true,
		},
		{
			name: "CookieForm",
			fields: fields{
				dataSource: MockRequestDataSource{
					cookies: []*http.Cookie{
						{Name: "example", Value: "R,100,G,200,B,150"},
					},
				},
			},
			args: args{
				in:      "cookie",
				name:    "example",
				style:   "form",
				explode: false,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
		{
			name: "CookieFormExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					cookies: []*http.Cookie{
						{Name: "R", Value: "100"},
						{Name: "G", Value: "200"},
						{Name: "B", Value: "150"},
					},
				},
			},
			args: args{
				in:      "cookie",
				name:    "example",
				style:   "form",
				explode: true,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
		{
			name: "CookieInvalid",
			fields: fields{
				dataSource: MockRequestDataSource{
					query: map[string][]string{
						"example": {"blue,black,brown"},
					},
				},
			},
			args: args{
				in:      "cookie",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantErr: true,
		},
		{
			name: "Form",
			fields: fields{
				dataSource: MockRequestDataSource{
					form: map[string][]string{
						"example": {`{"R":"100","G":"200","B":"150"}`},
					},
				},
			},
			args: args{
				in:      "form",
				name:    "example",
				style:   "ignored",
				explode: true,
			},
			wantResult: types.Pojo{"R": "100", "G": "200", "B": "150"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &EndpointRequestDecoder{
				DataSource: tt.fields.dataSource,
			}

			pf := ops.PortField{
				Group: tt.args.in,
				Peer:  tt.args.name,
				Options: map[string]string{
					"style":   tt.args.style,
					"explode": fmt.Sprintf("%v", tt.args.explode),
				},
			}

			gotResult, err := d.DecodeObject(&pf)
			if tt.wantErr != (err != nil) {
				if err == nil {
					assert.Errorf(t, err, "DecodeObject(%v, %v, %v, %v)", tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
				} else {
					assert.NoErrorf(t, err, "DecodeObject(%v, %v, %v, %v)", tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
				}
			}
			assert.Equalf(t, tt.wantResult, gotResult, "DecodeObject(%v, %v, %v, %v)", tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
		})
	}
}

func TestEndpointRequestDecoder_DecodeFile(t *testing.T) {
	tests := []struct {
		name       string
		formData   string
		boundary   string
		fieldName  string
		wantResult *multipart.FileHeader
		wantErr    bool
	}{
		{
			name:      "File",
			fieldName: "field2",
			formData: `--boundary
Content-Disposition: form-data; name="field1"

value1
--boundary
Content-Disposition: form-data; name="field2"; filename="example.txt"

value2
--boundary--`,
			boundary: "boundary",
			wantResult: &multipart.FileHeader{
				Filename: "example.txt",
				Header: textproto.MIMEHeader{
					HeaderContentDisposition: {`form-data; name="field2"; filename="example.txt"`},
				},
				Size: 6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multipartFormReader := multipart.NewReader(bytes.NewBufferString(tt.formData), tt.boundary)
			multipartForm, err := multipartFormReader.ReadForm(10 * 1024 * 1024)
			assert.NoError(t, err)

			d := &EndpointRequestDecoder{
				DataSource: MockRequestDataSource{
					multipartForm: multipartForm,
				},
			}

			pf := ops.PortField{
				Group: FieldGroupHttpForm,
				Peer:  tt.fieldName,
				Type: ops.PortFieldType{
					Shape:       ops.FieldShapeFile,
					Type:        reflect.TypeOf(tt.wantResult),
					HandlerType: reflect.TypeOf(tt.wantResult),
				},
			}

			gotResult, err := d.DecodeFile(&pf)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t,
					reflect.DeepEqual(gotResult.Filename, tt.wantResult.Filename),
					testhelpers.Diff(gotResult.Filename, tt.wantResult.Filename))
				assert.True(t,
					reflect.DeepEqual(gotResult.Header, tt.wantResult.Header),
					testhelpers.Diff(gotResult.Header, tt.wantResult.Header))
				assert.True(t,
					reflect.DeepEqual(gotResult.Size, tt.wantResult.Size),
					testhelpers.Diff(gotResult.Size, tt.wantResult.Size))
			}
		})
	}
}

func TestEndpointRequestDecoder_DecodeFileArray(t *testing.T) {
	tests := []struct {
		name        string
		formData    string
		boundary    string
		fieldName   string
		wantResults []*multipart.FileHeader
		wantErr     bool
	}{
		{
			name:      "FileArray",
			fieldName: "field1",
			formData: `--boundary
Content-Disposition: form-data; name="field1"; filename="example.txt"

value1
--boundary
Content-Disposition: form-data; name="field1"; filename="example2.txt"

value2
--boundary--`,
			boundary: "boundary",
			wantResults: []*multipart.FileHeader{
				{
					Filename: "example.txt",
					Header: textproto.MIMEHeader{
						HeaderContentDisposition: {`form-data; name="field1"; filename="example.txt"`},
					},
					Size: 6,
				},
				{
					Filename: "example2.txt",
					Header: textproto.MIMEHeader{
						HeaderContentDisposition: {`form-data; name="field1"; filename="example2.txt"`},
					},
					Size: 6,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multipartFormReader := multipart.NewReader(bytes.NewBufferString(tt.formData), tt.boundary)
			multipartForm, err := multipartFormReader.ReadForm(10 * 1024 * 1024)
			assert.NoError(t, err)

			d := NewRequestDecoder(MockRequestDataSource{
				multipartForm: multipartForm,
			})

			pf := ops.PortField{
				Group: FieldGroupHttpForm,
				Peer:  tt.fieldName,
				Type: ops.PortFieldType{
					Shape:       ops.FieldShapeFileArray,
					Type:        reflect.TypeOf(tt.wantResults),
					HandlerType: reflect.TypeOf(tt.wantResults),
				},
			}

			gotResults, err := d.DecodeFileArray(&pf)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				for i, wantResult := range tt.wantResults {
					gotResult := gotResults[i]
					assert.True(t,
						reflect.DeepEqual(gotResult.Filename, wantResult.Filename),
						testhelpers.Diff(gotResult.Filename, wantResult.Filename))
					assert.True(t,
						reflect.DeepEqual(gotResult.Header, wantResult.Header),
						testhelpers.Diff(gotResult.Header, wantResult.Header))
					assert.True(t,
						reflect.DeepEqual(gotResult.Size, wantResult.Size),
						testhelpers.Diff(gotResult.Size, wantResult.Size))
				}
			}
		})
	}
}
