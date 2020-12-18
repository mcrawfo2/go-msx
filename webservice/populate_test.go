package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_populateBody(t *testing.T) {
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
			test: new(RouteBuilderTest).
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
			test: new(RouteBuilderTest).
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

func Test_populateHeader(t *testing.T) {
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
			test: new(RouteBuilderTest).
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
			test: new(RouteBuilderTest).
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
			test: new(RouteBuilderTest).
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

func Test_populatePath(t *testing.T) {
	type params struct {
		PathId string `req:"path"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "PathElement",
			test: new(RouteBuilderTest).
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
			test: new(RouteBuilderTest).
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

func Test_populateQuery(t *testing.T) {
	type params struct {
		Required string `req:"query"`
		Optional *int   `req:"query"`
	}

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "AllValues",
			test: new(RouteBuilderTest).
				WithRequestMethod("GET").
				WithRouteTargetReturn(200).
				WithRequestQueryParameter("required", MIME_JSON).
				WithRequestQueryParameter("optional", "22").
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, MIME_JSON, p.Required)
					assert.Equal(t, 22, *p.Optional)
				}),
		},
		{
			name: "OptionalValue",
			test: new(RouteBuilderTest).
				WithRequestMethod("GET").
				WithRouteTargetReturn(200).
				WithRequestQueryParameter("required", MIME_XML).
				WithRequestVerifier(func(t *testing.T, req *restful.Request) {
					p := &params{}
					err := Populate(req, p)
					assert.NoError(t, err)
					assert.Equal(t, MIME_XML, p.Required)
					assert.True(t, nil == p.Optional)
				}),
		},
		{
			name: "MissingValue",
			test: new(RouteBuilderTest).
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

func Test_populateScalar(t *testing.T) {
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
			test: new(RouteBuilderTest).
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
			test: new(RouteBuilderTest).
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
			test: new(RouteBuilderTest).
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
