package types

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
)

var anonymousFuncCount int32

// nameOfFunction returns the short name of the function f for documentation.
// It uses a runtime feature for debugging ; its value may change for later Go versions.
func ShortFunctionName(f interface{}) string {
	fun := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	tokenized := strings.Split(fun.Name(), ".")
	last := tokenized[len(tokenized)-1]
	if last == "func1" { // this could mean conflicts in API docs
		val := atomic.AddInt32(&anonymousFuncCount, 1)
		last = "anonymousFunction" + fmt.Sprintf("%d", val)
		atomic.StoreInt32(&anonymousFuncCount, val)
	}
	return last
}

func FullFunctionName(f interface{}) string {
	fun := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	tokenized := strings.Split(fun.Name(), ".")
	last := tokenized[len(tokenized)-1]
	if last == "func1" { // this could mean conflicts in API docs
		val := atomic.AddInt32(&anonymousFuncCount, 1)
		last = "anonymousFunction" + fmt.Sprintf("%d", val)
		atomic.StoreInt32(&anonymousFuncCount, val)
		tokenized[len(tokenized)-1] = last
	}
	return strings.Join(tokenized, ".")
}