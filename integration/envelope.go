package integration

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type MsxEnvelope struct {
	Command    string                 `json:"command"`
	Debug      map[string]interface{} `json:"debug,omitempty"`
	Errors     []string               `json:"errors,omitempty"`
	HttpStatus string                 `json:"httpStatus"`
	Message    string                 `json:"message"`
	Params     map[string]interface{} `json:"params"`
	Payload    interface{}            `json:"responseObject"`
	Success    bool                   `json:"success"`
	Throwable  *Throwable             `json:"throwable,omitempty"`
}

func (e *MsxEnvelope) Error() error {
	return errors.New(e.Message)
}

func (e *MsxEnvelope) IsError() bool {
	return !e.Success && e.Message != ""
}

func (e *MsxEnvelope) StatusCode() int {
	return GetSpringStatusCodeForName(e.HttpStatus)
}

func NewEnvelope(payload interface{}) *MsxEnvelope {
	return &MsxEnvelope{
		Payload: payload,
	}
}

type Throwable struct {
	Cause      *Throwable   `json:"cause,omitempty"`
	StackTrace []StackFrame `json:"stackTrace,omitempty"`
	Message    string       `json:"message"`
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

func NewThrowable(err error) *Throwable {
	throwable := new(Throwable)

	if err == nil {
		throwable.Message = "Nil error"
		return throwable
	}

	// Parse message
	errMessage := err.Error()
	lines := strings.Split(errMessage, "\n")
	parts := strings.Split(lines[0], ": ")
	throwable.Message = parts[0]

	// Parse stack trace
	if errWithStack, ok := err.(stackTracer); ok {
		for _, frame := range errWithStack.StackTrace() {
			stackFrame := make(StackFrame)
			stackFrame.SetLineNumber(fmt.Sprintf("%d", frame))
			stackFrame.SetFileName(fmt.Sprintf("%s", frame))
			stackFrame.SetMethodName(fmt.Sprintf("%n", frame))

			extendedLocation := fmt.Sprintf("%+s", frame)
			extendedParts := strings.Split(extendedLocation, "\n")
			stackFrame.SetFullMethodName(extendedParts[0])
			stackFrame.SetFullFileName(extendedParts[1][1:])

			throwable.StackTrace = append(throwable.StackTrace, stackFrame)
		}
	}

	// Recurse
	if errWithCause, ok := err.(causer); ok {
		if cause := errWithCause.Cause(); cause != nil {
			throwableCause := NewThrowable(cause)
			// Skip over errors.Wrap artifacts
			if throwableCause.Message == throwable.Message && throwableCause.StackTrace == nil {
				throwableCause = throwableCause.Cause
			}
			throwable.Cause = throwableCause
		}
	}

	return throwable
}

type StackFrame map[string]interface{}

func (f StackFrame) SetLineNumber(lineNumber string) {
	lineNumberInt, _ := strconv.Atoi(lineNumber)
	f["lineNumber"] = lineNumberInt
}

func (f StackFrame) SetFileName(fileName string) {
	f["fileName"] = fileName
}

func (f StackFrame) SetMethodName(methodName string) {
	f["methodName"] = methodName
}

func (f StackFrame) SetFullFileName(extendedFileName string) {
	f["fullFileName"] = extendedFileName
}

func (f StackFrame) SetFullMethodName(extendedMethodName string) {
	f["fullMethodName"] = extendedMethodName
}
