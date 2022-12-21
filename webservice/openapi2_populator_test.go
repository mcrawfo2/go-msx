// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"reflect"
	"testing"
	"time"
)

func TestRouteParam_populateBody(t *testing.T) {
	type body struct {
		A int    `json:"a"`
		B string `json:"b"`
	}

	type params struct {
		Body body `req:"body"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "JsonStruct",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRouteTargetReturn(200).
				WithRequestHeader("Content-Type", MIME_JSON).
				WithRequestBodyString(`{"a":3,"b":"c"}`).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, 3, p.Body.A)
					assert.Equal(t, "c", p.Body.B)
				}),
		},
		{
			name: "NoBody",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.Error(t, err)
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRouteParam_populateHeader(t *testing.T) {
	type params struct {
		XCustomParameter *int   `req:"header"`
		ContentType      string `req:"header"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "AllValues",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("GET").
				WithRouteTargetReturn(200).
				WithRequestHeader("Content-Type", MIME_JSON).
				WithRequestHeader("X-Custom-Parameter", "22").
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, 22, *p.XCustomParameter)
					assert.Equal(t, MIME_JSON, p.ContentType)
				}),
		},
		{
			name: "OptionalValue",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("GET").
				WithRouteTargetReturn(200).
				WithRequestHeader("Content-Type", MIME_XML).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.True(t, nil == p.XCustomParameter)
					assert.Equal(t, MIME_XML, p.ContentType)
				}),
		},
		{
			name: "MissingValue",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("GET").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.Error(t, err)
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRouteParam_populatePath(t *testing.T) {
	type params struct {
		PathId string `req:"path"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "PathElement",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestPath("/populate/path/xyz").
				WithRoutePath("/populate/path/{pathId}").
				WithRouteParameter(restful.PathParameter("pathId", "")).
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, "xyz", p.PathId)
				}),
		},
		{
			name: "PathTail",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestPath("/populate/path/xyz/abc").
				WithRoutePath("/populate/path/{pathId:*}").
				WithRouteParameter(restful.PathParameter("pathId", "")).
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, "xyz/abc", p.PathId)
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRouteParam_populateQuery(t *testing.T) {
	type params struct {
		Required  string       `req:"query"`
		Optional  *int         `req:"query"`
		Uuid      types.UUID   `req:"query"`
		MultiUuid []types.UUID `req:"query,multi"`
		Csv       []string     `req:"query,csv"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "AllValues",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("GET").
				WithRouteTargetReturn(200).
				WithRequestQueryParameter("required", MIME_JSON).
				WithRequestQueryParameter("optional", "22").
				WithRequestQueryParameter("uuid", "703665b9-6d89-4bda-b786-3645ce75b699").
				WithRequestQueryParameter("multiUuid", "e9db0eee-dd1f-4ef6-b749-4ad810cb1842").
				WithRequestQueryParameter("multiUuid", "096c431d-94b7-4b2a-8257-9990eb841c1c").
				WithRequestQueryParameter("csv", "a,b,c").
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, MIME_JSON, p.Required)
					assert.Equal(t, 22, *p.Optional)
					assert.Equal(t, "703665b9-6d89-4bda-b786-3645ce75b699", p.Uuid.String())
					assert.Len(t, p.MultiUuid, 2)
					assert.Equal(t, "e9db0eee-dd1f-4ef6-b749-4ad810cb1842", p.MultiUuid[0].String())
					assert.Equal(t, "096c431d-94b7-4b2a-8257-9990eb841c1c", p.MultiUuid[1].String())
					assert.Equal(t, []string{"a", "b", "c"}, p.Csv)
				}),
		},
		{
			name: "OptionalValue",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("GET").
				WithRouteTargetReturn(200).
				WithRequestQueryParameter("required", MIME_XML).
				WithRequestQueryParameter("uuid", "703665b9-6d89-4bda-b786-3645ce75b699").
				WithRequestQueryParameter("multiUuid", "e9db0eee-dd1f-4ef6-b749-4ad810cb1842").
				WithRequestQueryParameter("multiUuid", "096c431d-94b7-4b2a-8257-9990eb841c1c").
				WithRequestQueryParameter("csv", "a,b,c").
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, MIME_XML, p.Required)
					assert.True(t, nil == p.Optional)
					assert.Equal(t, "703665b9-6d89-4bda-b786-3645ce75b699", p.Uuid.String())
					assert.Len(t, p.MultiUuid, 2)
					assert.Equal(t, "e9db0eee-dd1f-4ef6-b749-4ad810cb1842", p.MultiUuid[0].String())
					assert.Equal(t, "096c431d-94b7-4b2a-8257-9990eb841c1c", p.MultiUuid[1].String())
					assert.Equal(t, []string{"a", "b", "c"}, p.Csv)
				}),
		},
		{
			name: "MissingValue",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("GET").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.Error(t, err)
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRouteParam_populateForm_File(t *testing.T) {

	type params struct {
		Field string                `req:"form"`
		File  *multipart.FileHeader `req:"form,file"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "AllValues",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRouteTargetReturn(200).
				WithRequestFormFieldParameter("field", "def456").
				WithRequestFormFileParameter("file", "example.txt", "abc123").
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)

					assert.Equal(t, "def456", p.Field)

					assert.Equal(t, "example.txt", p.File.Filename)

					reader, err := p.File.Open()
					assert.NoError(t, err)

					contents, err := ioutil.ReadAll(reader)
					assert.NoError(t, err)

					assert.Equal(t, "abc123", string(contents))
				}),
		},
		{
			name: "BodyFile",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRouteTargetReturn(200).
				WithRequestQueryParameter("field", "def456").
				WithRequestBodyString("abc123").
				WithoutRequestHeader("Content-Type").
				WithRequestHeader("Content-Type", "application/json").
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)

					assert.Equal(t, "def456", p.Field)

					assert.Equal(t, "body", p.File.Filename)

					reader, err := p.File.Open()
					assert.NoError(t, err)

					contents, err := ioutil.ReadAll(reader)
					assert.NoError(t, err)

					assert.Equal(t, "abc123", string(contents))
				}),
		},
		{
			name: "MissingValue",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("GET").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.Error(t, err)
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRouteParam_populateScalar(t *testing.T) {
	type scalars struct {
		String  string  `req:"query"`
		Int     int     `req:"query"`
		Int8    int8    `req:"query"`
		Int16   int16   `req:"query"`
		Int32   int32   `req:"query"`
		Int64   int64   `req:"query"`
		Float32 float32 `req:"query"`
		Float64 float64 `req:"query"`
		Uint    uint    `req:"query"`
		Uint8   uint8   `req:"query"`
		Uint16  uint16  `req:"query"`
		Uint32  uint32  `req:"query"`
		Uint64  uint64  `req:"query"`
		Bool    bool    `req:"query"`
	}

	type pointers struct {
		String  *string  `req:"query"`
		Int     *int     `req:"query"`
		Int8    *int8    `req:"query"`
		Int16   *int16   `req:"query"`
		Int32   *int32   `req:"query"`
		Int64   *int64   `req:"query"`
		Float32 *float32 `req:"query"`
		Float64 *float64 `req:"query"`
		Uint    *uint    `req:"query"`
		Uint8   *uint8   `req:"query"`
		Uint16  *uint16  `req:"query"`
		Uint32  *uint32  `req:"query"`
		Uint64  *uint64  `req:"query"`
		Bool    *bool    `req:"query"`
	}

	type custom struct {
		Time *types.Time `req:"query"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Scalars",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestQueryParameter("string", "xyz").
				WithRequestQueryParameter("int", "1").
				WithRequestQueryParameter("int8", "2").
				WithRequestQueryParameter("int16", "3").
				WithRequestQueryParameter("int32", "4").
				WithRequestQueryParameter("int64", "5").
				WithRequestQueryParameter("float32", "6.0").
				WithRequestQueryParameter("float64", "7.0").
				WithRequestQueryParameter("uint", "8").
				WithRequestQueryParameter("uint8", "9").
				WithRequestQueryParameter("uint16", "10").
				WithRequestQueryParameter("uint32", "11").
				WithRequestQueryParameter("uint64", "12").
				WithRequestQueryParameter("bool", "true").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &scalars{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, "xyz", p.String)
					assert.Equal(t, int(1), p.Int)
					assert.Equal(t, int8(2), p.Int8)
					assert.Equal(t, int16(3), p.Int16)
					assert.Equal(t, int32(4), p.Int32)
					assert.Equal(t, int64(5), p.Int64)
					assert.Equal(t, float32(6), p.Float32)
					assert.Equal(t, float64(7), p.Float64)
					assert.Equal(t, uint(8), p.Uint)
					assert.Equal(t, uint8(9), p.Uint8)
					assert.Equal(t, uint16(10), p.Uint16)
					assert.Equal(t, uint32(11), p.Uint32)
					assert.Equal(t, uint64(12), p.Uint64)
					assert.Equal(t, true, p.Bool)
				}),
		},
		{
			name: "Pointers",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestQueryParameter("string", "xyz").
				WithRequestQueryParameter("int", "1").
				WithRequestQueryParameter("int8", "2").
				WithRequestQueryParameter("int16", "3").
				WithRequestQueryParameter("int32", "4").
				WithRequestQueryParameter("int64", "5").
				WithRequestQueryParameter("float32", "6.0").
				WithRequestQueryParameter("float64", "7.0").
				WithRequestQueryParameter("uint", "8").
				WithRequestQueryParameter("uint8", "9").
				WithRequestQueryParameter("uint16", "10").
				WithRequestQueryParameter("uint32", "11").
				WithRequestQueryParameter("uint64", "12").
				WithRequestQueryParameter("bool", "true").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &pointers{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, "xyz", *p.String)
					assert.Equal(t, int(1), *p.Int)
					assert.Equal(t, int8(2), *p.Int8)
					assert.Equal(t, int16(3), *p.Int16)
					assert.Equal(t, int32(4), *p.Int32)
					assert.Equal(t, int64(5), *p.Int64)
					assert.Equal(t, float32(6), *p.Float32)
					assert.Equal(t, float64(7), *p.Float64)
					assert.Equal(t, uint(8), *p.Uint)
					assert.Equal(t, uint8(9), *p.Uint8)
					assert.Equal(t, uint16(10), *p.Uint16)
					assert.Equal(t, uint32(11), *p.Uint32)
					assert.Equal(t, uint64(12), *p.Uint64)
					assert.Equal(t, true, *p.Bool)
				}),
		},
		{
			name: "Custom",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestQueryParameter("time", "2020-12-18T19:44:14Z").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &custom{}
					err := Populate(req, p)
					assert.NoError(t, err)

					parsedTime := p.Time.ToTimeTime()
					expectedTime := time.Unix(1608320654, 0)

					assert.True(t, expectedTime.Equal(parsedTime))
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRouteParam_populateScaler_sanitize(t *testing.T) {
	type XssBody struct {
		Name        string  `json:"name" san:"xss"`
		Description *string `json:"description"`
		Ignored     string  `json:"ignored" san:"-"`
	}

	type Scalars struct {
		String string `req:"query,san" san:"xss"`
	}

	type Bodies struct {
		Body XssBody `req:"body,san" san:"xss"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Scalar",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestQueryParameter("string", "x<a>yz</a>").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &Scalars{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, "xyz", p.String)
				}),
		},
		{
			name: "Body",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestHeader("Content-Type", MIME_JSON).
				WithRequestBodyString(`{"name": "a<b>c</b>d", "description": "d<c>b</c>a", "ignored": "<a>bcd</a>"}`).
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &Bodies{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, "acd", p.Body.Name)
					assert.Equal(t, types.NewStringPtr("dba"), p.Body.Description)
					assert.Equal(t, "<a>bcd</a>", p.Body.Ignored)
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRouteParam_populateSlice(t *testing.T) {
	testID1 := "56b36e70-72bc-11eb-a576-ef1d86082c19"
	testID2 := "e9db0eee-dd1f-4ef6-b749-4ad810cb1842"

	type stringSlice struct {
		Strings []string `req:"query,multi"`
	}

	type optionalStringSlice struct {
		Strings *[]string `req:"query,multi"`
	}

	type uuidSlice struct {
		Uuids []types.UUID `req:"query,multi"`
	}

	type optionalUuidSlice struct {
		Uuids *[]types.UUID `req:"query,multi"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Strings",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestQueryParameter("strings", "a").
				WithRequestQueryParameter("strings", "b").
				WithRequestQueryParameter("strings", "c").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &stringSlice{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, p.Strings, []string{"a", "b", "c"})
				}),
		},
		{
			name: "OptionalStrings",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestQueryParameter("strings", "a").
				WithRequestQueryParameter("strings", "b").
				WithRequestQueryParameter("strings", "c").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &optionalStringSlice{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.NotNil(t, p.Strings)
					assert.Equal(t, *p.Strings, []string{"a", "b", "c"})
				}),
		},
		{
			name: "MissingOptionalStrings",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &optionalStringSlice{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Nil(t, p.Strings)
				}),
		},
		{
			name: "UUIDs",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestQueryParameter("uuids", testID1).
				WithRequestQueryParameter("uuids", testID2).
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &uuidSlice{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Len(t, p.Uuids, 2)
					assert.Equal(t, p.Uuids[0].String(), testID1)
					assert.Equal(t, p.Uuids[1].String(), testID2)
				}),
		},
		{
			name: "OptionalUuids",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRequestQueryParameter("uuids", testID1).
				WithRequestQueryParameter("uuids", testID2).
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &optionalUuidSlice{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.NotNil(t, p.Uuids)

					uuids := *p.Uuids
					assert.Len(t, uuids, 2)
					assert.Equal(t, uuids[0].String(), testID1)
					assert.Equal(t, uuids[1].String(), testID2)
				}),
		},
		{
			name: "MissingOptionalUuids",
			test: new(webservicetest.RouteBuilderTest).
				WithRequestMethod("POST").
				WithRouteTargetReturn(200).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &optionalUuidSlice{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Nil(t, p.Uuids)
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRouteParam_populateSlice2(t *testing.T) {
	testID := "56b36e70-72bc-11eb-a576-ef1d86082c19"

	uuidType := reflect.TypeOf(types.UUID{})

	uuidSliceType := reflect.SliceOf(uuidType)

	type fields struct {
		Field     reflect.StructField
		Source    string
		Name      string
		Options   map[string]string
		Parameter restful.ParameterData
	}

	type args struct {
		fieldValue reflect.Value
		values     []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Uuid",
			fields: fields{
				Field: reflect.StructField{
					Name: "Id",
					Type: uuidType,
				},
				Source:    requestTagSourceQuery,
				Name:      "id",
				Parameter: restful.ParameterData{Name: "id", DataType: "uuid"},
				Options:   map[string]string{},
			},
			args: args{
				fieldValue: reflect.New(uuidType),
				values:     []string{testID},
			},
			wantErr: false,
		},
		{
			name: "MultiUuid",
			fields: fields{
				Field: reflect.StructField{
					Name: "tenantId",
					Type: uuidSliceType,
				},
				Source:    requestTagSourceQuery,
				Name:      "tenantId",
				Parameter: restful.ParameterData{Name: "tenantId", AllowMultiple: true, Required: true},
				Options: map[string]string{
					"multi": "true",
				},
			},
			args: args{
				fieldValue: reflect.New(uuidSliceType).Elem(),
				values:     []string{testID},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RouteParam{
				Field:     tt.fields.Field,
				Source:    tt.fields.Source,
				Name:      tt.fields.Name,
				Options:   tt.fields.Options,
				Parameter: tt.fields.Parameter,
			}
			if err := r.populateField(tt.args.fieldValue, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("populateField() error = %v, wantErr %v", err, tt.wantErr)
			}
			returned := fmt.Sprintf("%s", tt.args.fieldValue)
			if tt.args.fieldValue.Kind() == reflect.Slice {
				returned = fmt.Sprintf("%s", tt.args.fieldValue.Index(0))
			}

			if returned != testID {
				t.Errorf("populateField() got value = %v, wanted %v", returned, testID)
			}
		})
	}
}
