package cassandra

import "testing"

func Test_createBatchKey(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "0", args: args{i: 0}, want: "a"},
		{name: "1", args: args{i: 1}, want: "b"},
		{name: "15", args: args{i: 15}, want: "p"},
		{name: "16", args: args{i: 16}, want: "aa"},
		{name: "17", args: args{i: 17}, want: "ab"},
		{name: "255", args: args{i: 255}, want: "op"},
		{name: "256", args: args{i: 256}, want: "aaa"},
		{name: "257", args: args{i: 257}, want: "aab"},
		{name: "4096", args: args{i: 4096}, want: "aaaa"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createBatchKey(tt.args.i); got != tt.want {
				t.Errorf("createBatchKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
