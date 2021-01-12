package integration

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
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
	Cause      *Throwable             `json:"cause,omitempty"`
	StackTrace []types.BackTraceFrame `json:"stackTrace,omitempty"`
	Message    string                 `json:"message"`
}

func NewThrowable(err error) *Throwable {
	if err == nil {
		throwable := new(Throwable)
		throwable.Message = "Nil error"
		return throwable
	}

	var result *Throwable
	for _, bte := range types.BackTraceFromError(err) {
		throwable := new(Throwable)
		// OWASP: https://owasp.org/www-community/Improper_Error_Handling
		// throwable.StackTrace = bte.Frames
		throwable.Message = bte.Message
		throwable.Cause = result
		result = throwable
	}

	return result
}
