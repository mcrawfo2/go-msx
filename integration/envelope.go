package integration

import (
	"strconv"

	"github.com/pkg/errors"
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
