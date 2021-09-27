package types

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
)

var anonymousFuncCount int32
var ErrSourceDirUnavailable = errors.New("Source directory could not be calculated")

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

func FindEntryPointDirFromStack() (string, error) {
	file, ok := getEntryPointFile(2)
	if !ok {
		return "", ErrSourceDirUnavailable
	}

	return filepath.Dir(file), nil
}

func FindSourceDirFromStack() (string, error) {
	file, ok := getEntryPointFile(2)
	if !ok {
		return "", ErrSourceDirUnavailable
	}

	thence := FindSourceDirFromFile(file)
	if thence == "" {
		return "", ErrSourceDirUnavailable
	}

	return thence, nil
}

var entryPointSuffixes = map[string]int{
	"go-msx/app.Run":  1,
	"testing.tRunner": -1,
	"main.main": 0,
}

// Hack when fs.sources is missing
func getEntryPointFile(skip int) (string, bool) {
	pcs := make([]uintptr, 32)
	frameCount := runtime.Callers(skip+1, pcs)
	frames := runtime.CallersFrames(pcs[:frameCount])
	var frameList []runtime.Frame
	var targetFrame = -1
	for {
		frame, more := frames.Next()
		frameList = append(frameList, frame)

		for suffix, offset := range entryPointSuffixes {
			if strings.HasSuffix(frame.Function, suffix) {
				targetFrame = len(frameList) - 1 + offset
				break
			}
		}
		if !more {
			break
		}
	}

	if targetFrame > -1 {
		if frameList[targetFrame].File != "" {
			return frameList[targetFrame].File, true
		}
	}

	return "", false
}

func FindSourceDirFromFile(whence string) string {
	for whence != "/" {
		whence = filepath.Dir(whence)
		gomod := filepath.Join(whence, "go.mod")
		_, err := os.Stat(gomod)
		if !os.IsNotExist(err) {
			return whence
		}
	}

	return ""
}
