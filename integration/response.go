package integration

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type MsxResponse struct {
	StatusCode int          `json:"statusCode"`
	Status     string       `json:"status"`
	Headers    http.Header  `json:"headers"`
	Envelope   *MsxEnvelope `json:"envelope"`
	Payload    interface{}  `json:"payload"`
	Body       []byte       `json:"-"`
	BodyString string       `json:"body"`
}

type OAuthErrorDTO struct {
	ErrorCode   string `json:"error"`
	Description string `json:"error_description"`
}

func (e *OAuthErrorDTO) Error() error {
	return errors.New(e.Description)
}

func (e *OAuthErrorDTO) IsError() bool {
	return e.Description != ""
}

type ErrorDTO struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Path       string `json:"path"`
	HttpStatus string `json:"httpStatus"`
	Timestamp  string `json:"timestamp"`
}

func (e *ErrorDTO) Error() error {
	return errors.New(e.Message)
}

func (e *ErrorDTO) IsError() bool {
	return e.Message != ""
}

func (e *ErrorDTO) SetError(code int, err error, path string) {
	e.Code = fmt.Sprintf("%d", code)
	if springHttpStatus, ok := springHttpStatusByCode[code]; ok {
		e.HttpStatus = springHttpStatus.Name()
	}
	e.Message = err.Error()
	e.Path = path
	e.Timestamp = time.Now().Format(time.RFC3339)
}

type ErrorDTO2 struct {
	Status    int    `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Path      string `json:"path"`
	Message   string `json:"message"`
}

func (e *ErrorDTO2) Error() error {
	return errors.New(e.Message)
}

func (e *ErrorDTO2) IsError() bool {
	return e.Message != ""
}

type ErrorDTO3 struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	ErrorName string `json:"error"`
	Message   string `json:"message"`
	Path      string `json:"path"`
}

func (e *ErrorDTO3) Error() error {
	return errors.New(e.Message)
}

func (e *ErrorDTO3) IsError() bool {
	return e.Message != ""
}

type Pojo map[string]interface{}

type PojoArray []map[string]interface{}

type HealthDTO struct {
	Status string `json:"status"`
}

type HealthResult struct {
	Response *MsxResponse
	Payload  *HealthDTO
}
