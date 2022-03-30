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
		E1    types.UUID     `req:"query"`
		E2    *types.UUID    `req:"query,required"`
		E3    *types.UUID    `req:"query"`
		E4    types.UUID     `req:"query,optional"`
		F1    types.Binary   `req:"query"`
		F2    *types.Binary  `req:"query,required"`
		F3    *types.Binary  `req:"query"`
		F4    types.Binary   `req:"query,optional"`
		G1    []byte         `req:"query"`
		G2    *[]byte        `req:"query,required" explode:"true"`
		G3    *[]byte        `req:"query"`
		G4    []byte         `req:"query,optional"`
		H1    []rune         `req:"query"`
		H2    *[]rune        `req:"query,required" explode:"true"`
		H3    *[]rune        `req:"query"`
		H4    []rune         `req:"query,optional"`
		N1    []types.Binary `req:"query"`
		N2    []types.Binary `req:"query,required" explode:"true"`
		N3    []types.Binary `req:"query"`
		N4    []types.Binary `req:"query,optional"`
		O1    []string       `req:"query"`
		O2    []string       `req:"query,required" explode:"true"`
		O3    []string       `req:"query"`
		O4    []string       `req:"query,optional"`
		P1    []int          `req:"query"`
		P2    []int          `req:"query,required" explode:"true"`
		P3    []int          `req:"query"`
		P4    []int          `req:"query,optional"`
		Q1    []types.Time   `req:"query"`
		Q2    []types.Time   `req:"query,required" explode:"true"`
		Q3    []types.Time   `req:"query"`
		Q4    []types.Time   `req:"query,optional"`
		R1    []types.UUID   `req:"query"`
		R2    []types.UUID   `req:"query,required" explode:"true"`
		R3    []types.UUID   `req:"query"`
		R4    []types.UUID   `req:"query,optional"`
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
					"e1":       {"24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839"},
					"e2":       {"248b0dd5-e4e7-4b1e-9279-c500ba761942"},
					"f1":       {"abc"},
					"f2":       {"def"},
					"g1":       {"abc"},
					"g2":       {"def"},
					"h1":       {"abc"},
					"h2":       {"def"},
					"n1":       {"abc,def"},
					"n2":       {"abc", "def"},
					"o1":       {"red,green,yellow", "black", "brown"},
					"o2":       {"blue", "black", "brown"},
					"p1":       {"210,220,230,240"},
					"p2":       {"250", "260", "270", "280"},
					"q1":       {"2019-08-12T07:20:50.520273482Z,2019-09-12T07:20:50.520273482Z"},
					"q2":       {"2019-07-12T07:20:50.520273482Z", "2019-06-12T07:20:50.520273482Z"},
					"r1":       {"24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839,248b0dd5-e4e7-4b1e-9279-c500ba761942"},
					"r2":       {"24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839", "248b0dd5-e4e7-4b1e-9279-c500ba761942"},
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
		E1: types.MustParseUUID("24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839"),
		E2: types.NewUUIDPtr(types.MustParseUUID("248b0dd5-e4e7-4b1e-9279-c500ba761942")),
		F1: types.NewBinary([]byte("abc")),
		F2: types.NewBinaryPtr([]byte("def")),
		G1: []byte("abc"),
		G2: types.NewByteSlicePtr([]byte("def")),
		H1: []rune("abc"),
		H2: types.NewRuneSlicePtr([]rune("def")),
		N1: []types.Binary{types.NewBinaryFromString("abc"), types.NewBinaryFromString("def")},
		N2: []types.Binary{types.NewBinaryFromString("abc"), types.NewBinaryFromString("def")},
		O1: []string{"red", "green", "yellow"},
		O2: []string{"blue", "black", "brown"},
		P1: []int{210, 220, 230, 240},
		P2: []int{250, 260, 270, 280},
		Q1: []types.Time{
			mustParseTime("2019-08-12T07:20:50.520273482Z"),
			mustParseTime("2019-09-12T07:20:50.520273482Z"),
		},
		Q2: []types.Time{
			mustParseTime("2019-07-12T07:20:50.520273482Z"),
			mustParseTime("2019-06-12T07:20:50.520273482Z"),
		},
		R1: []types.UUID{
			types.MustParseUUID("24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839"),
			types.MustParseUUID("248b0dd5-e4e7-4b1e-9279-c500ba761942"),
		},
		R2: []types.UUID{
			types.MustParseUUID("24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839"),
			types.MustParseUUID("248b0dd5-e4e7-4b1e-9279-c500ba761942"),
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
		E1    types.UUID     `req:"form"`
		E2    *types.UUID    `req:"form,required"`
		E3    *types.UUID    `req:"form"`
		E4    types.UUID     `req:"form,optional"`
		F1    types.Binary   `req:"form"`
		F2    *types.Binary  `req:"form,required"`
		F3    *types.Binary  `req:"form"`
		F4    types.Binary   `req:"form,optional"`
		G1    []byte         `req:"form"`
		G2    *[]byte        `req:"form,required" explode:"true"`
		G3    *[]byte        `req:"form"`
		G4    []byte         `req:"form,optional"`
		H1    []rune         `req:"form"`
		H2    *[]rune        `req:"form,required" explode:"true"`
		H3    *[]rune        `req:"form"`
		H4    []rune         `req:"form,optional"`
		N1    []types.Binary `req:"form"`
		N2    []types.Binary `req:"form,required" explode:"true"`
		N3    []types.Binary `req:"form"`
		N4    []types.Binary `req:"form,optional"`
		O1    []string       `req:"form"`
		O2    []string       `req:"form,required" explode:"true"`
		O3    []string       `req:"form"`
		O4    []string       `req:"form,optional"`
		P1    []int          `req:"form"`
		P2    []int          `req:"form,required" explode:"true"`
		P3    []int          `req:"form"`
		P4    []int          `req:"form,optional"`
		Q1    []types.Time   `req:"form"`
		Q2    []types.Time   `req:"form,required" explode:"true"`
		Q3    []types.Time   `req:"form"`
		Q4    []types.Time   `req:"form,optional"`
		R1    []types.UUID   `req:"form"`
		R2    []types.UUID   `req:"form,required" explode:"true"`
		R3    []types.UUID   `req:"form"`
		R4    []types.UUID   `req:"form,optional"`
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
					"e1":    {"24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839"},
					"e2":    {"248b0dd5-e4e7-4b1e-9279-c500ba761942"},
					"f1":    {"abc"},
					"f2":    {"def"},
					"g1":    {"abc"},
					"g2":    {"def"},
					"h1":    {"abc"},
					"h2":    {"def"},
					"n1":    {"abc,def"},
					"n2":    {"abc,def"},
					"o1":    {"red,green,yellow,black,brown"},
					"o2":    {"blue,black,brown"},
					"p1":    {"210,220,230,240"},
					"p2":    {"250,260,270,280"},
					"q1":    {"2019-08-12T07:20:50.520273482Z,2019-09-12T07:20:50.520273482Z"},
					"q2":    {"2019-07-12T07:20:50.520273482Z,2019-06-12T07:20:50.520273482Z"},
					"r1":    {"24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839,248b0dd5-e4e7-4b1e-9279-c500ba761942"},
					"r2":    {"24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839,248b0dd5-e4e7-4b1e-9279-c500ba761942"},
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
		E1: types.MustParseUUID("24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839"),
		E2: types.NewUUIDPtr(types.MustParseUUID("248b0dd5-e4e7-4b1e-9279-c500ba761942")),
		F1: types.NewBinary([]byte("abc")),
		F2: types.NewBinaryPtr([]byte("def")),
		G1: []byte("abc"),
		G2: types.NewByteSlicePtr([]byte("def")),
		H1: []rune("abc"),
		H2: types.NewRuneSlicePtr([]rune("def")),
		N1: []types.Binary{types.NewBinaryFromString("abc"), types.NewBinaryFromString("def")},
		N2: []types.Binary{types.NewBinaryFromString("abc"), types.NewBinaryFromString("def")},
		O1: []string{"red", "green", "yellow", "black", "brown"},
		O2: []string{"blue", "black", "brown"},
		P1: []int{210, 220, 230, 240},
		P2: []int{250, 260, 270, 280},
		Q1: []types.Time{
			mustParseTime("2019-08-12T07:20:50.520273482Z"),
			mustParseTime("2019-09-12T07:20:50.520273482Z"),
		},
		Q2: []types.Time{
			mustParseTime("2019-07-12T07:20:50.520273482Z"),
			mustParseTime("2019-06-12T07:20:50.520273482Z"),
		},
		R1: []types.UUID{
			types.MustParseUUID("24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839"),
			types.MustParseUUID("248b0dd5-e4e7-4b1e-9279-c500ba761942"),
		},
		R2: []types.UUID{
			types.MustParseUUID("24ba04f0-2f5b-4e1d-a7d0-ebc4f9e61839"),
			types.MustParseUUID("248b0dd5-e4e7-4b1e-9279-c500ba761942"),
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
