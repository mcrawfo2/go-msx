package webservice

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewBadRequestError(t *testing.T) {
	err := NewBadRequestError(errors.New("some error"))
	assert.Error(t, err)
	assert.Equal(t, 400, err.(*StatusError).StatusCode())
}

func TestNewConflictError(t *testing.T) {
	err := NewBadRequestError(errors.New("some error"))
	assert.Error(t, err)
	assert.Equal(t, 409, err.(*StatusError).StatusCode())
}

func TestNewForbiddenError(t *testing.T) {
	err := NewBadRequestError(errors.New("some error"))
	assert.Error(t, err)
	assert.Equal(t, 403, err.(*StatusError).StatusCode())
}

func TestNewInternalError(t *testing.T) {
	err := NewBadRequestError(errors.New("some error"))
	assert.Error(t, err)
	assert.Equal(t, 500, err.(*StatusError).StatusCode())
}

func TestNewNotFoundError(t *testing.T) {
	err := NewBadRequestError(errors.New("some error"))
	assert.Error(t, err)
	assert.Equal(t, 404, err.(*StatusError).StatusCode())
}

func TestNewStatusCodeProvider(t *testing.T) {
	body := struct{}{}
	status := 200
	provider := NewStatusCodeProvider(body, status)
	assert.Equal(t, 200, provider.StatusCode())
}

func TestNewStatusError(t *testing.T) {
	type args struct {
		err error
		status int
	}
	type want struct {
		status int
		message string
	}
	var tests = []struct {
		name string
		args args
		want want
	} {
		{
			name: "Cause",
			args: args{
				err:    errors.New("some error"),
				status: 404,
			},
			want: want{
				status: 404,
				message: "some error",
			},
		},
		{
			name: "NoCause",
			args: args{
				err:    nil,
				status: 404,
			},
			want: want{
				status: 404,
				message: "Unknown status error: 404",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewStatusError(tt.args.err, tt.args.status)
			assert.Equal(t, tt.want.status, provider.(*StatusError).StatusCode())
			assert.Equal(t, tt.want.message, provider.(*StatusError).Error())

		})
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	err := NewBadRequestError(errors.New("some error"))
	assert.Error(t, err)
	assert.Equal(t, 401, err.(*StatusError).StatusCode())
}

func TestStatusError_Cause(t *testing.T) {
	someError := errors.New("some error")
	err := NewBadRequestError(someError)
	assert.Error(t, err)
	assert.Equal(t, someError, err.(*StatusError).Cause())
}

func TestStatusError_Error(t *testing.T) {
	errText := "some error"
	someError := errors.New(errText)
	err := NewBadRequestError(someError)
	assert.Error(t, err)
	assert.Equal(t, someError, err.(*StatusError).Error())
}

func TestStatusError_StatusCode(t *testing.T) {
	errText := "some error"
	someError := errors.New(errText)
	err := NewBadRequestError(someError)
	assert.Error(t, err)
	assert.Equal(t, 400, err.(*StatusError).StatusCode())
}

func TestStatusError_Unwrap(t *testing.T) {
	someError := errors.New("some error")
	err := NewBadRequestError(someError)
	assert.Error(t, err)
	assert.Equal(t, someError, err.(*StatusError).Unwrap())
}

func Test_statusCodeProviderImpl_MarshalJSON(t *testing.T) {
	type fields struct {
		body       interface{}
		statusCode int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusCodeProviderImpl{
				body:       tt.fields.body,
				statusCode: tt.fields.statusCode,
			}
			got, err := s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statusCodeProviderImpl_StatusCode(t *testing.T) {
	type fields struct {
		body       interface{}
		statusCode int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusCodeProviderImpl{
				body:       tt.fields.body,
				statusCode: tt.fields.statusCode,
			}
			if got := s.StatusCode(); got != tt.want {
				t.Errorf("StatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
