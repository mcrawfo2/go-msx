package config

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type intStruct struct {
	A int
	B int16
	C *int   `config:"optional"`
	D *int32 `config:"default=${a}${b}"`
	E int64  `config:"optional"`
}

type stringStruct struct {
	A string
	B *string `config:"optional"`
	C *string `config:"default=${a}${b}"`
	D string  `config:"optional"`
}

type mapStructSubStruct struct {
	A string
	B int32
	C string `config:"optional"`
	D int    `config:"default=42"`
}

type mapStruct struct {
	A map[string]string
	B map[string]*string
	C map[string]mapStructSubStruct
	D map[string]*mapStructSubStruct
}

type sliceStruct struct {
	A []int `config:"default=1;2;3"`
	B []bool `config:"default=true;false;true"`
	C []*string `config:"default=alpha;bravo;charlie"`
	D []*float64 `config:"default=${f};${g}"`
	E []mapStructSubStruct
}

type durationStruct struct {
	A time.Duration `config:"default=10m"`
	B *time.Duration `config:"optional"`
}

type structSubStruct struct {
	A string `config:"optional"`
	B string
}

type structStruct struct {
	Env struct {
		All    map[string]string
		Linux  map[string]string
		Darwin map[string]string
	}
	Sub *structSubStruct
}

type validateStruct struct {
	A time.Duration `config:"default=10m"`
}

func (v *validateStruct) Validate() error {
	return types.ErrorMap{
		"a": validation.Validate(&v.A, validation.Required, validation.Min(time.Minute), validation.Max(30 * time.Minute)),
	}
}

func TestPopulateStruct(t *testing.T) {
	D1 := int32(12)
	C3 := "alpha"
	B5 := "charlie"
	C6 := "bravo"
	D6 := float64(1.2)
	D7 := float64(2.3)
	B8 := 600 * time.Second

	tests := []struct {
		name    string
		values  map[string]string
		prefix  string
		arg     interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "IntStruct",
			values: map[string]string{
				"a": "1",
				"b": "2",
			},
			arg: &intStruct{
				E: int64(5),
			},
			want: &intStruct{
				A: 1,
				B: 2,
				C: nil,
				D: &D1,
				E: int64(5),
			},
		},
		{
			name: "IntStructMissingKey",
			values: map[string]string{
				"a": "1",
			},
			arg:     &intStruct{},
			wantErr: true,
		},
		{
			name: "StringStruct",
			values: map[string]string{
				"a": "alpha",
			},
			arg: &stringStruct{
				D: "delta",
			},
			want: &stringStruct{
				A: "alpha",
				B: nil,
				C: &C3,
				D: "delta",
			},
		},
		{
			name: "MapStruct",
			values: map[string]string{
				"a.a": "alpha",
				"a.b": "beta",
				"b.a": "charlie",
				"c.d.a": "charlie",
				"c.d.b": "33",
				"c.d.d": "34",
				"d.e.a": "delta",
				"d.e.b": "21",
			},
			arg: &mapStruct{},
			want: &mapStruct{
				A: map[string]string{
					"a": "alpha",
					"b": "beta",
				},
				B: map[string]*string{
					"a": &B5,
				},
				C: map[string]mapStructSubStruct{
					"d": {
						A: "charlie",
						B: 33,
						D: 34,
					},
				},
				D: map[string]*mapStructSubStruct{
					"e": {
						A: "delta",
						B: 21,
						D: 42,
										},
				},
			},
		},
		{
			name: "SliceStruct",
			values: map[string]string{
				"a[0]": "1",
				"a[2]": "2",
				"a[4]": "3",
				"b[0]": "true",
				"b[1]": "false",
				"b[2]": "true",
				"c[0]": "alpha",
				"c[1]": "bravo",
				"c[3]": "charlie",
				"d[0]": "1.2",
				"e[0].a": "foxtrot",
				"e[0].b": "21",
				"e[2].a": "golf",
				"e[2].b": "20",
				"e[2].c": "hotel",
				"e[2].d": "19",
			},
			arg: &sliceStruct{},
			want: &sliceStruct{
				A: []int{1,2,3},
				B: []bool{true,false,true},
				C: []*string{&C3,&C6,&B5},
				D: []*float64{&D6},
				E: []mapStructSubStruct{
					{
						A: "foxtrot",
						B: 21,
						C: "",
						D: 42,
					},
					{
						A: "golf",
						B: 20,
						C: "hotel",
						D: 19,
					},
				},
			},
		},
		{
			name: "SliceStructDefaults",
			values: map[string]string{
				"f": "1.2",
				"g": "2.3",
			},
			arg: &sliceStruct{},
			want: &sliceStruct{
				A: []int{1,2,3},
				B: []bool{true,false,true},
				C: []*string{&C3,&C6,&B5},
				D: []*float64{&D6,&D7},
			},
		},
		{
			name:    "DurationStruct",
			values:  map[string]string{
				"b": "600s",
			},
			arg:     &durationStruct{},
			want:    &durationStruct{
				A: 10 * time.Minute,
				B: &B8,
			},
			wantErr: false,
		},
		{
			name: "StructStruct",
			values: map[string]string{
				"env.all.GOROOT": "/usr/local/go",
				"env.linux.GOPATH": "/home/ubuntu/go",
			},
			arg: &structStruct{},
			want: &structStruct{
				Env: struct {
					All    map[string]string
					Linux  map[string]string
					Darwin map[string]string
				}{
					All:    map[string]string{"goroot": "/usr/local/go"},
					Linux:  map[string]string{"gopath": "/home/ubuntu/go"},
					Darwin: map[string]string{},
				},
			},
		},
		{
			name: "ValidatableStruct",
			values: map[string]string{},
			arg: &validateStruct{},
			want: &validateStruct{
				A: time.Minute * 10,
			},
			wantErr: false,
		},
		{
			name: "ValidatableStructError",
			values: map[string]string{
				"a": "1h",
			},
			arg: &validateStruct{},
			want: &validateStruct{
				A: time.Hour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inMem := NewInMemoryProvider(tt.name, tt.values)
			cfg := NewConfig(inMem)
			err := cfg.Load(context.Background())
			assert.NoError(t, err)

			err = Populate(tt.arg, tt.prefix, cfg.OriginalValues())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.arg)
			}
		})
	}
}

func TestPopulateSlice(t *testing.T) {
	var intSlice []int
	var intSliceResult = []int{1,2,3}
	tests := []struct {
		name    string
		values  map[string]string
		prefix  string
		arg     interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "IntSlice",
			values: map[string]string{
				"a[0]": "1",
				"a[1]": "2",
				"a[3]": "3",
			},
			prefix: "a",
			arg: &intSlice,
			want: &intSliceResult,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inMem := NewInMemoryProvider(tt.name, tt.values)
			cfg := NewConfig(inMem)
			err := cfg.Load(context.Background())
			assert.NoError(t, err)

			err = Populate(tt.arg, tt.prefix, cfg.OriginalValues())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.arg)
			}
		})
	}

}

