package types

import (
	"reflect"
	"testing"
)

func TestGetTypeName(t *testing.T) {
	type args struct {
		instance interface{}
		root bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "UUID",
			args: args{ instance: new(UUID), root: true },
			want: "types.UUID",
		},
		{
			name: "[]UUID",
			args: args{ instance: []UUID{}, root: false },
			want: "List«types.UUID»",
		},
		{
			name: "Time",
			args: args{ instance: new(Time), root: true },
			want: "types.Time",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTypeName(reflect.TypeOf(tt.args.instance), tt.args.root); got != tt.want {
				t.Errorf("GetTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}
