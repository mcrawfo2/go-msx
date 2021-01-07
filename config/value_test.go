package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Implementation(t *testing.T) {
	var _ Values = new(SnapshotValues)
	var _ Values = new(Snapshot)
	var _ Values = new(Config)
}

func Test_Value_StringPtr(t *testing.T) {
	var v = "bravo"

	tests := []struct {
		name string
		val  string
		want interface{}
	}{
		{
			name: "String",
			val:  v,
			want: &v,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Value(tt.val)
			got := v.StringPtr()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Value_Int(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    int64
		wantErr bool
	}{
		{
			name: "Valid",
			val:  "42",
			want: 42,
		},
		{
			name:    "Invalid",
			val:     "err",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Value(tt.val)
			got, err := v.Int()
			if tt.wantErr != (err != nil) {
				t.Errorf("Value.Int() got err = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Value_Uint(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    uint64
		wantErr bool
	}{
		{
			name: "Valid",
			val:  "42",
			want: 42,
		},
		{
			name:    "Invalid",
			val:     "-40",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Value(tt.val)
			got, err := v.Uint()
			if tt.wantErr != (err != nil) {
				t.Errorf("Value.Uint() got err = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Value_Float(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    float64
		wantErr bool
	}{
		{
			name: "Valid",
			val:  "42.5",
			want: 42.5,
		},
		{
			name:    "Invalid",
			val:     "-err",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Value(tt.val)
			got, err := v.Float()
			if tt.wantErr != (err != nil) {
				t.Errorf("Value.Float() got err = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Value_Bool(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    bool
		wantErr bool
	}{
		{
			name: "ValidTrue",
			val:  "true",
			want: true,
		},
		{
			name: "ValidFalse",
			val:  "false",
			want: false,
		},
		{
			name:    "Invalid",
			val:     "-err",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Value(tt.val)
			got, err := v.Bool()
			if tt.wantErr != (err != nil) {
				t.Errorf("Value.Bool() got err = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Value_StringSlice(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		sep     string
		want    []string
	}{
		{
			name: "Zero",
			val: "",
			sep: ";",
			want: nil,
		},
		{
			name: "Single",
			val:  "alpha",
			sep:  ",",
			want: []string{"alpha"},
		},
		{
			name: "Multiple",
			val:  "alpha,bravo",
			sep:  ",",
			want: []string{"alpha", "bravo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Value(tt.val)
			got := v.StringSlice(tt.sep)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Value_Duration(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    time.Duration
		wantErr bool
	}{
		{
			name: "Valid",
			val:  "10s",
			want: 10 * time.Second,
		},
		{
			name:    "Invalid",
			val:     "err",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Value(tt.val)
			got, err := v.Duration()
			if tt.wantErr != (err != nil) {
				t.Errorf("Value.Bool() got err = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
