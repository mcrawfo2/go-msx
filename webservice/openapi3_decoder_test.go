package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
)

type MockRequestDataSource struct {
	cookies        []*http.Cookie
	headers        http.Header
	form           url.Values
	multipartForm  *multipart.Form
	query          url.Values
	pathParameters map[string]string
	body           []byte
}

func (m MockRequestDataSource) Body() ([]byte, error) {
	return m.body, nil
}

func (m MockRequestDataSource) Cookies() []*http.Cookie {
	return m.cookies
}

func (m MockRequestDataSource) Headers() http.Header {
	return m.headers
}

func (m MockRequestDataSource) Form() (url.Values, *multipart.Form, error) {
	return m.form, m.multipartForm, nil
}

func (m MockRequestDataSource) Query() url.Values {
	return m.query
}

func (m MockRequestDataSource) PathParameters() map[string]string {
	return m.pathParameters
}

func (m MockRequestDataSource) ReadEntity(e interface{}) (err error) {
	return json.Unmarshal(m.body, e)
}

func TestOpenApiRequestDecoder_DecodePrimitive(t *testing.T) {
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
		wantResult types.OptionalString
		wantErr    bool
	}{
		{
			name: "HeaderSimple",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"example": {"blue"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "simple",
				explode: false,
			},
			wantResult: types.NewOptionalStringFromString("blue"),
			wantErr:    false,
		},
		{
			name: "HeaderSimpleExplode",
			fields: fields{
				dataSource: MockRequestDataSource{
					headers: map[string][]string{
						"example": {"blue"},
					},
				},
			},
			args: args{
				in:      "header",
				name:    "example",
				style:   "simple",
				explode: true,
			},
			wantResult: types.NewOptionalStringFromString("blue"),
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
			wantResult: types.NewOptionalStringFromString("blue"),
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
			wantResult: types.NewOptionalStringFromString("blue"),
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
			wantResult: types.NewOptionalStringFromString("blue"),
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
			wantResult: types.NewOptionalStringFromString("blue black brown"),
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
			wantResult: types.NewOptionalStringFromString("blue|black|brown"),
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
			wantResult: types.NewOptionalStringFromString("blue,black,brown"),
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
			wantResult: types.NewOptionalStringFromString("blue"),
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
			wantResult: types.NewOptionalStringFromString("blue"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &OpenApiRequestDecoder{
				DataSource: tt.fields.dataSource,
			}
			gotResult, err := d.DecodePrimitive(tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
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

func TestOpenApiRequestDecoder_DecodeArray(t *testing.T) {
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
						"example": {"blue,black,brown"},
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
						"example": {"blue,black,brown"},
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
				style:   "form",
				explode: false,
			},
			wantResult: []string{"blue", "black", "brown"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &OpenApiRequestDecoder{
				DataSource: tt.fields.dataSource,
			}
			gotResult, err := d.DecodeArray(tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
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

func TestOpenApiRequestDecoder_DecodeObject(t *testing.T) {
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
						"example": {"R,100,G,200,B,150"},
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
						"example": {"R=100,G=200,B=150"},
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
						"example[R]": {"100"},
						"example[G]": {"200"},
						"example[B]": {"150"},
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
			d := &OpenApiRequestDecoder{
				DataSource: tt.fields.dataSource,
			}
			gotResult, err := d.DecodeObject(tt.args.in, tt.args.name, tt.args.style, tt.args.explode)
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
