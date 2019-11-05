package integration

import "github.com/pkg/errors"

type MsxEnvelope struct {
	Command    string                 `json:"command"`
	Debug      map[string]interface{} `json:"debug,omitEmpty"`
	Errors     []string               `json:"errors,omitEmpty"`
	HttpStatus string                 `json:"httpStatus"`
	Message    string                 `json:"message"`
	Params     map[string]interface{} `json:"params"`
	Payload    interface{}            `json:"responseObject"`
	Success    bool                   `json:"success"`
	Throwable  interface{} 			  `json:"throwable,omitEmpty"`
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