package types

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime/debug"
	"strconv"
	"testing"
)

var globalError = fmt.Errorf("Global Error")
var globalStackError = errors.New("Global Stack Error")
var _, parseError = strconv.ParseBool("trueq")

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

func testError7() error {
	return errors.Wrap(parseError, "Wrapped Parse Error")
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
			args: args{err: testError5()},
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
							FullFile:   "backtrace_test.go",
						},
					},
				},
				{
					Message: "Wrapped Stack Error",
					Frames: []BackTraceFrame{
						{
							FullMethod: "types.testError5",
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
			name: "ParseError",
			args: args{err: testError7()},
			want: BackTrace{
				{
					Message: `strconv.ParseBool: parsing "trueq": invalid syntax`,
					Frames:  []BackTraceFrame{},
				},
				{
					Message: "Wrapped Parse Error",
					Frames: []BackTraceFrame{
						{
							FullMethod: "types.testError7",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BackTraceFromError(tt.args.err)
			t.Logf("Got: \n%s", got.Stanza())

			if !got.Equal(tt.want) {
				t.Logf("Wanted: \n%s", tt.want.Stanza())
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

func TestBackTraceErrorFromDebugStackTrace(t *testing.T) {
	type args struct {
		stackTraceBytes []byte
	}
	tests := []struct {
		name string
		args args
		want BackTraceError
	}{
		{
			name: "Standard",
			args: args{
				stackTraceBytes: []byte(`goroutine 19 [running]:
runtime/debug.Stack(0x0, 0xc000036500, 0x37e)
	/usr/local/go/src/runtime/debug/stack.go:24 +0x9d
runtime/debug.PrintStack()
	/usr/local/go/src/runtime/debug/stack.go:16 +0x22
cto-github.cisco.com/NFV-BU/go-msx/types.TestPanic.func1()
	/Users/mcrawfo2/vms-3.1/go-msx/types/backtrace_test.go:314 +0x42
panic(0x12f1920, 0x13d1c30)
	/usr/local/go/src/runtime/panic.go:967 +0x15d
cto-github.cisco.com/NFV-BU/go-msx/types.TestPanic(0xc0000d47e0)
	/Users/mcrawfo2/vms-3.1/go-msx/types/backtrace_test.go:317 +0x5b
testing.tRunner(0xc0000d47e0, 0x137ec48)
	/usr/local/go/src/testing/testing.go:992 +0xdc
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:1043 +0x357`),
			},
			want: BackTraceError{
				Message: "panic",
				Frames: []BackTraceFrame{
					{
						Method:     "TestPanic",
						FullMethod: "cto-github.cisco.com/NFV-BU/go-msx/types.TestPanic",
						File:       "backtrace_test.go",
					},
					{
						Method:     "tRunner",
						FullMethod: "testing.tRunner",
						File:       "testing.go",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BackTraceErrorFromDebugStackTrace(tt.args.stackTraceBytes); !tt.want.Equal(got) {
				t.Errorf("BackTraceErrorFromDebugStackTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func GeneratePanicBacktrace(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
		}
	}()
	panic("implement me")
}
