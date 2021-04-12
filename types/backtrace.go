package types

import (
	"bytes"
	"github.com/pkg/errors"
	"runtime"
	"strconv"
	"strings"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

type BackTraceFrame struct {
	pc         uintptr `json:"-"`
	Method     string `json:"methodName"`
	FullMethod string `json:"fullMethodName"`
	FullFile   string `json:"fullFileName"`
	File       string `json:"fileName"`
	Line       int    `json:"lineNumber"`
}

func (b BackTraceFrame) Equals(other BackTraceFrame) bool {
	return b.pc == other.pc
}

func BackTraceFrameFromFrame(frame errors.Frame) BackTraceFrame {
	pc := uintptr(frame) - 1
	fn := runtime.FuncForPC(pc)

	fullMethod := "unknown"
	method := "unknown"
	fullFile := ""
	file := ""
	line := 0

	if fn != nil {
		fullFile, line = fn.FileLine(pc)
		fullMethod = fn.Name()

		// Short function name
		i := strings.LastIndex(fullMethod, "/")
		method = fullMethod[i+1:]
		i = strings.Index(method, ".")
		method = method[i+1:]

		// Short file name
		i = strings.LastIndex(fullFile, "/")
		file = fullFile[i+1:]
	}

	return BackTraceFrame{
		pc:         pc,
		FullMethod: fullMethod,
		Method:     method,
		FullFile:   fullFile,
		File:       file,
		Line:       line,
	}
}

type BackTraceError struct {
	Message string
	Frames  []BackTraceFrame
}

func (b BackTraceError) Equal(other BackTraceError) bool {
	if b.Message != other.Message {
		return false
	}
	if len(b.Frames) != len(other.Frames) {
		return false
	}
	return true
}

func (b BackTraceError) TrimFrames(other BackTraceError) BackTraceError {
	myCount := len(b.Frames)
	otherCount := len(other.Frames)

	if myCount < otherCount {
		return b
	}

	if otherCount == 0 {
		return b
	}

	// Trim matching stack traces
	idx := 0
	for b.Frames[myCount-idx-1].Equals(other.Frames[otherCount-idx-1]) {
		idx++
	}

	return BackTraceError{
		Message: b.Message,
		Frames:  b.Frames[:myCount-idx],
	}
}

func BackTraceErrorFromError(err error) BackTraceError {
	errMessage := err.Error()

	// Default message is the entire first line
	lines := strings.Split(errMessage, "\n")
	message := lines[0]

	// If we have a cause, trim message after the first colon if it
	// matches the cause's message
	if cause := errors.Unwrap(err); cause != err && cause != nil {
		parts := strings.SplitN(lines[0], ": ", 2)
		causerMessage := cause.Error()
		if (len(parts) == 2 && parts[1] == causerMessage) || (errMessage == causerMessage) {
			message = parts[0]
		}
	}

	var frames []BackTraceFrame
	if errStackTracer, ok := err.(stackTracer); ok {
		for _, errorFrame := range errStackTracer.StackTrace() {
			btf := BackTraceFrameFromFrame(errorFrame)
			frames = append(frames, btf)
		}
	}

	return BackTraceError{
		Message: message,
		Frames:  frames,
	}
}

type BackTrace []BackTraceError

func (b BackTrace) Stanza() string {
	var buf bytes.Buffer
	for i, bte := range b {
		if i == 0 {
			buf.WriteString("Root Cause: ")
		} else {
			buf.WriteString("Caused: ")
		}

		buf.WriteString(bte.Message)
		buf.WriteRune('\n')

		for _, frame := range bte.Frames {
			buf.WriteString("  ")
			buf.WriteString(frame.FullMethod)
			buf.WriteRune('\n')

			buf.WriteString("    ")
			buf.WriteString(frame.FullFile)
			buf.WriteRune(':')
			buf.WriteString(strconv.Itoa(frame.Line))
			buf.WriteRune('\n')
		}
	}
	return buf.String()
}

func (b BackTrace) Causes() []string {
	var buf []string
	for i, bte := range b {
		if i == 0 {
			buf = append(buf, "Root Cause: "+bte.Message)
		} else {
			buf = append(buf, "Caused: "+bte.Message)
		}
	}

	return buf
}

func (b BackTrace) Lines() []string {
	var buf []string
	for i, bte := range b {
		if i == 0 {
			buf = append(buf, "Root Cause: "+bte.Message)
		} else {
			buf = append(buf, "Caused: "+bte.Message)
		}

		for _, frame := range bte.Frames {
			buf = append(buf, "  "+frame.FullMethod)
			buf = append(buf, "    "+frame.FullFile+":"+strconv.Itoa(frame.Line))
		}
	}

	return buf
}

func (b BackTrace) Reverse() BackTrace {
	result := make(BackTrace, len(b))
	for i, bt := range b {
		result[len(b)-i-1] = bt
	}
	return result
}

func (b BackTrace) Equal(other BackTrace) bool {
	if len(b) != len(other) {
		return false
	}
	for i := range b {
		if !b[i].Equal(other[i]) {
			return false
		}
	}

	return true
}

func BackTraceFromError(err error) BackTrace {
	var result BackTrace
	var prev BackTraceError
	for err != nil {
		cur := BackTraceErrorFromError(err)

		// De-duplicate
		if len(result) == 0 || cur.Frames != nil || cur.Message != prev.Message {
			result = append(BackTrace{cur}, result...)
			prev = cur
		}

		if errCauser, ok := err.(causer); ok {
			err = errCauser.Cause()
		} else {
			err = nil
		}
	}

	for i := 0; i < len(result)-1; i++ {
		// Squash stack lines
		result[i] = result[i].TrimFrames(result[i+1])
	}

	return result
}
