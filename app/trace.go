package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, activateTracing)
	OnEvent(EventFinal, PhaseAfter, deactivateTracing)
}

func activateTracing(ctx context.Context) error {
	logger.Info("Activating tracing")
	return trace.ConfigureTracer(ctx)
}

func deactivateTracing(ctx context.Context) error {
	logger.Info("Deactivating tracing")
	return trace.ShutdownTracer(ctx)
}

var anonymousFuncCount int32

// nameOfFunction returns the short name of the function f for documentation.
// It uses a runtime feature for debugging ; its value may change for later Go versions.
func nameOfFunction(f interface{}) string {
	fun := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	tokenized := strings.Split(fun.Name(), ".")
	last := tokenized[len(tokenized)-1]
	if last == "func1" { // this could mean conflicts in API docs
		val := atomic.AddInt32(&anonymousFuncCount, 1)
		last = "func" + fmt.Sprintf("%d", val)
		atomic.StoreInt32(&anonymousFuncCount, val)
	}
	return last
}
