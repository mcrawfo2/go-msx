package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

type cmyk struct {
	C int
	M int
	Y int
	K int
}

func TestOpenApiRequestPopulator_PopulatePortStruct_Query(t *testing.T) {
	type portStruct struct {
		A1    int            `req:"query"`
		A2    *int           `req:"query" required:"true"`
		A3    *int           `req:"query"`
		A4    int            `req:"query" optional:"true"`
		B1    bool           `req:"query"`
		B2    *bool          `req:"query,required"`
		B3    *bool          `req:"query"`
		B4    bool           `req:"query,optional"`
		C1    string         `req:"query"`
		C2    *string        `req:"query,required"`
		C3    *string        `req:"query"`
		C4    string         `req:"query,optional"`
		D1    types.Time     `req:"query"`
		D2    *types.Time    `req:"query,required"`
		D3    *types.Time    `req:"query"`
		D4    types.Time     `req:"query,optional"`
		E1    []byte         `req:"query"`
		E2    *[]byte        `req:"query,required"`
		E3    *[]byte        `req:"query"`
		E4    []byte         `req:"query,optional"`
		F1    []string       `req:"query"`
		F2    []string       `req:"query,required" explode:"true"`
		F3    []string       `req:"query"`
		F4    []string       `req:"query,optional"`
		G1    []int          `req:"query"`
		G2    []int          `req:"query,required" explode:"true"`
		G3    []int          `req:"query"`
		G4    []int          `req:"query,optional"`
		H1    []types.Time   `req:"query"`
		H2    []types.Time   `req:"query,required" explode:"true"`
		H3    []types.Time   `req:"query"`
		H4    []types.Time   `req:"query,optional"`
		RGB   map[string]int `req:"query=rgb" style:"deepObject" explode:"true"`
		CMYK1 cmyk           `req:"query=cmyk" style:"deepObject" explode:"true"`
		CMYK2 *cmyk          `req:"query=cmyk2,required" style:"deepObject" explode:"true"`
		CMYK3 *cmyk          `req:"query" style:"deepObject" explode:"true"`
		CMYK4 *cmyk          `req:"query,optional" style:"deepObject" explode:"true"`
	}

	p := OpenApiRequestPopulator{
		Decoder: &OpenApiRequestDecoder{
			DataSource: MockRequestDataSource{
				query: url.Values{
					"a1":       {"10"},
					"a2":       {"20"},
					"b1":       {"true"},
					"b2":       {"false"},
					"c1":       {"100"},
					"c2":       {"110"},
					"d1":       {"2019-10-12T07:20:50.520273482Z"},
					"d2":       {"2019-11-12T07:20:50.520273482Z"},
					"e1":       {"abc"},
					"e2":       {"def"},
					"f1":       {"red,green,yellow", "black", "brown"},
					"f2":       {"blue", "black", "brown"},
					"g1":       {"210,220,230,240"},
					"g2":       {"250", "260", "270", "280"},
					"h1":       {"2019-08-12T07:20:50.520273482Z,2019-09-12T07:20:50.520273482Z"},
					"h2":       {"2019-07-12T07:20:50.520273482Z", "2019-06-12T07:20:50.520273482Z"},
					"rgb[r]":   {"100"},
					"rgb[g]":   {"200"},
					"rgb[b]":   {"150"},
					"cmyk[c]":  {"75"},
					"cmyk[m]":  {"30"},
					"cmyk[y]":  {"77"},
					"cmyk[k]":  {"11"},
					"cmyk2[c]": {"77"},
					"cmyk2[m]": {"32"},
					"cmyk2[y]": {"79"},
					"cmyk2[k]": {"13"},
				},
			},
		},
	}

	mustParseTime := func(ts string) types.Time {
		tv, err := types.ParseTime(ts)
		assert.NoError(t, err)
		return tv
	}

	want := &portStruct{
		A1: 10,
		A2: types.NewIntPtr(20),
		B1: true,
		B2: types.NewBoolPtr(false),
		C1: "100",
		C2: types.NewStringPtr("110"),
		C4: "",
		D1: mustParseTime("2019-10-12T07:20:50.520273482Z"),
		D2: types.NewTimePtr(mustParseTime("2019-11-12T07:20:50.520273482Z")),
		E1: []byte("abc"),
		E2: types.NewByteSlicePtr([]byte("def")),
		F1: []string{"red", "green", "yellow"},
		F2: []string{"blue", "black", "brown"},
		G1: []int{210, 220, 230, 240},
		G2: []int{250, 260, 270, 280},
		H1: []types.Time{
			mustParseTime("2019-08-12T07:20:50.520273482Z"),
			mustParseTime("2019-09-12T07:20:50.520273482Z"),
		},
		H2: []types.Time{
			mustParseTime("2019-07-12T07:20:50.520273482Z"),
			mustParseTime("2019-06-12T07:20:50.520273482Z"),
		},
		RGB: map[string]int{
			"r": 100,
			"g": 200,
			"b": 150,
		},
		CMYK1: cmyk{
			C: 75,
			M: 30,
			Y: 77,
			K: 11,
		},
		CMYK2: &cmyk{
			C: 77,
			M: 32,
			Y: 79,
			K: 13,
		},
	}

	request := NewEndpointRequest().WithPortStruct(portStruct{})
	got, err := p.PopulatePortStruct(request)
	assert.NoError(t, err, testhelpers.Diff(nil, err))
	assert.Equalf(t, want, got, testhelpers.Diff(want, got))
}

func TestOpenApiRequestPopulator_PopulatePortStruct_Form(t *testing.T) {
	type portStruct struct {
		A1    int            `req:"form"`
		A2    *int           `req:"form" required:"true"`
		A3    *int           `req:"form"`
		A4    int            `req:"form" optional:"true"`
		B1    bool           `req:"form"`
		B2    *bool          `req:"form,required"`
		B3    *bool          `req:"form"`
		B4    bool           `req:"form,optional"`
		C1    string         `req:"form"`
		C2    *string        `req:"form,required"`
		C3    *string        `req:"form"`
		C4    string         `req:"form,optional"`
		D1    types.Time     `req:"form"`
		D2    *types.Time    `req:"form,required"`
		D3    *types.Time    `req:"form"`
		D4    types.Time     `req:"form,optional"`
		E1    []byte         `req:"form"`
		E2    *[]byte        `req:"form,required"`
		E3    *[]byte        `req:"form"`
		E4    []byte         `req:"form,optional"`
		F1    []string       `req:"form"`
		F2    []string       `req:"form,required"`
		F3    []string       `req:"form"`
		F4    []string       `req:"form,optional"`
		G1    []int          `req:"form"`
		G2    []int          `req:"form,required"`
		G3    []int          `req:"form"`
		G4    []int          `req:"form,optional"`
		H1    []types.Time   `req:"form"`
		H2    []types.Time   `req:"form,required"`
		H3    []types.Time   `req:"form"`
		H4    []types.Time   `req:"form,optional"`
		RGB   map[string]int `req:"form=rgb"`
		CMYK  []cmyk         `req:"form"`
		CMYK1 cmyk           `req:"form=cmyk1"`
		CMYK2 *cmyk          `req:"form=cmyk2,required"`
		CMYK3 *cmyk          `req:"form"`
		CMYK4 *cmyk          `req:"form,optional"`
		// TODO: Single,Multi Files
	}

	p := OpenApiRequestPopulator{
		Decoder: &OpenApiRequestDecoder{
			DataSource: MockRequestDataSource{
				form: url.Values{
					"a1":    {"10"},
					"a2":    {"20"},
					"b1":    {"true"},
					"b2":    {"false"},
					"c1":    {"100"},
					"c2":    {"110"},
					"d1":    {"2019-10-12T07:20:50.520273482Z"},
					"d2":    {"2019-11-12T07:20:50.520273482Z"},
					"e1":    {"abc"},
					"e2":    {"def"},
					"f1":    {"red,green,yellow,black,brown"},
					"f2":    {"blue,black,brown"},
					"g1":    {"210"},
					"g2":    {"250,260,270,280"},
					"h1":    {"2019-08-12T07:20:50.520273482Z,2019-09-12T07:20:50.520273482Z"},
					"h2":    {"2019-07-12T07:20:50.520273482Z,2019-06-12T07:20:50.520273482Z"},
					"rgb":   {`{"r": 100,"g": 200,"b": 150}`},
					"cmyk1": {`{"c":75,"m":30,"y":77,"k":11}`},
					"cmyk2": {`{"c":77,"m":32,"y":79,"k":13}`},
				},
			},
		},
	}

	mustParseTime := func(ts string) types.Time {
		tv, err := types.ParseTime(ts)
		assert.NoError(t, err)
		return tv
	}

	want := &portStruct{
		A1: 10,
		A2: types.NewIntPtr(20),
		B1: true,
		B2: types.NewBoolPtr(false),
		C1: "100",
		C2: types.NewStringPtr("110"),
		C4: "",
		D1: mustParseTime("2019-10-12T07:20:50.520273482Z"),
		D2: types.NewTimePtr(mustParseTime("2019-11-12T07:20:50.520273482Z")),
		E1: []byte("abc"),
		E2: types.NewByteSlicePtr([]byte("def")),
		F1: []string{"red", "green", "yellow", "black", "brown"},
		F2: []string{"blue", "black", "brown"},
		G1: []int{210},
		G2: []int{250, 260, 270, 280},
		H1: []types.Time{
			mustParseTime("2019-08-12T07:20:50.520273482Z"),
			mustParseTime("2019-09-12T07:20:50.520273482Z"),
		},
		H2: []types.Time{
			mustParseTime("2019-07-12T07:20:50.520273482Z"),
			mustParseTime("2019-06-12T07:20:50.520273482Z"),
		},
		RGB: map[string]int{
			"r": 100,
			"g": 200,
			"b": 150,
		},
		CMYK1: cmyk{
			C: 75,
			M: 30,
			Y: 77,
			K: 11,
		},
		CMYK2: &cmyk{
			C: 77,
			M: 32,
			Y: 79,
			K: 13,
		},
	}

	request := NewEndpointRequest().WithPortStruct(portStruct{})
	got, err := p.PopulatePortStruct(request)
	assert.NoError(t, err, testhelpers.Diff(nil, err))
	assert.Equalf(t, want, got, testhelpers.Diff(want, got))
}

func TestOpenApiRequestPopulator_PopulatePortStruct_Header(t *testing.T) {}

func TestOpenApiRequestPopulator_PopulatePortStruct_Cookie(t *testing.T) {}

func TestOpenApiRequestPopulator_PopulatePortStruct_Path(t *testing.T) {}

func TestOpenApiRequestPopulator_PopulatePortStruct_Body(t *testing.T) {}
