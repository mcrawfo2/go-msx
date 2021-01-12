package types

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

var globalError = fmt.Errorf("Global Error")
var globalStackError = errors.New("Global Stack Error")

func testError3() error {
	return fmt.Errorf("Built-In Error")
}

func testError2() error {
	return errors.Wrap(testError3(), "Wrapped Error")
}

func testError1() error {
	return errors.Wrap(testError2(), "Re-wrapped Error")
}

func testError4() error {
	return errors.New("Stack Error")
}

func testError5() error {
	return errors.Wrap(testError4(), "Wrapped Stack Error")
}

func testError6() error {
	return errors.Wrap(testError5(), "Re-wrapped Stack Error")
}

func TestBackTraceFromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want BackTrace
	}{
		{
			name: "BuiltIn",
			args: args{err: testError3()},
			want: BackTrace{
				{
					Message: "Built-In Error",
					Frames:  nil,
				},
			},
		},
		{
			name: "Wrapped",
			args: args{err: testError2()},
			want: BackTrace{
				{
					Message: "Built-In Error",
					Frames:  nil,
				},
				{
					Message: "Wrapped Error",
					Frames: []BackTraceFrame{
						{
							FullMethod: "types.testError2",
							FullFile:   "backtrace_test.go",
						},
						{
							FullMethod: "types.TestBackTraceFromError",
							FullFile:   "backtrace_test.go",
						},
						{
							FullMethod: "testing.tRunner",
							FullFile:   "proc.go",
						},
						{
							FullMethod: "runtime.goexit",
							FullFile:   "asm_amd64.s",
						},
					},
				},
			},
		},
		{
			name: "Re-wrapped",
			args: args{err: testError1()},
			want: BackTrace{
				{
					Message: "Built-In Error",
					Frames:  nil,
				},
				{
					Message: "Wrapped Error",
					Frames: []BackTraceFrame{

						{
							FullMethod: "types.testError2",
							FullFile:   "backtrace_test.go",
						},
						{
							FullMethod: "types.testError1",
							FullFile:   "backtrace_test.go",
						},
					},
				},
				{
					Message: "Re-wrapped Error",
					Frames: []BackTraceFrame{
						{
							FullMethod: "types.testError1",
							FullFile:   "backtrace_test.go",
						},
						{
							FullMethod: "types.TestBackTraceFromError",
							FullFile:   "backtrace_test.go",
						},
						{
							FullMethod: "testing.tRunner",
							FullFile:   "proc.go",
						},
						{
							FullMethod: "runtime.goexit",
							FullFile:   "asm_amd64.s",
						},
					},
				},
			},
		},
		{
			name: "Global",
			args: args{err: globalError},
			want: BackTrace{
				{
					Message: "Global Error",
					Frames:  nil,
				},
			},
		},
		{
			name: "GlobalStackError",
			args: args{err: globalStackError},
			want: BackTrace{
				{
					Message: "Global Stack Error",
					Frames: []BackTraceFrame{
						{
							FullMethod: "types.init",
							FullFile:   "backtrace_test.go",
						},
						{
							FullMethod: "runtime.doInit",
							FullFile:   "proc.go",
						},
						{
							FullMethod: "runtime.doInit",
							FullFile:   "proc.go",
						},
						{
							FullMethod: "runtime.main",
							FullFile:   "proc.go",
						},
						{
							FullMethod: "runtime.goexit",
							FullFile:   "asm_amd64.s",
						},
					},
				},
			},
		},
		{
			name: "StackWrapped",
			args:args{err: testError5()},
			want: BackTrace{
				{
					Message: "Stack Error",
					Frames: []BackTraceFrame{

						{
							FullMethod: "types.testError4",
							FullFile:   "backtrace_test.go",
						},
						{
							FullMethod: "types.testError5",
							FullFile: "backtrace_test.go",
						},
					},
				},
				{
					Message: "Wrapped Stack Error",
					Frames: []BackTraceFrame{
						{
							FullMethod: "types.testError5",
							FullFile: "backtrace_test.go",
						},
						{
							FullMethod: "types.TestBackTraceFromError",
							FullFile:   "backtrace_test.go",
						},
						{
							FullMethod: "testing.tRunner",
							FullFile:   "proc.go",
						},
						{
							FullMethod: "runtime.goexit",
							FullFile:   "asm_amd64.s",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BackTraceFromError(tt.args.err); !got.Equal(tt.want) {
				t.Errorf("BackTraceFromError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackTrace_Stanza(t *testing.T) {
	t.Run("ReWrappedStack", func(t *testing.T) {
		stanza := BackTraceFromError(testError6()).Stanza()
		fmt.Println(stanza)
	})

	t.Run("ReWrappedGlobal", func(t *testing.T) {
		stanza := BackTraceFromError(testError1()).Stanza()
		fmt.Println(stanza)
	})
}
