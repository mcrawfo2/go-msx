package retry

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPermanentErrorDecorator(t *testing.T) {
	fnOk := func() error { return nil }
	fnErr := func() error { return errors.New("some error") }

	tests := []struct {
		name      string
		retryable Retryable
		wantErr   bool
	}{
		{
			name:      "Success",
			retryable: fnOk,
			wantErr:   false,
		},
		{
			name:      "PermanentError",
			retryable: fnErr,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fnDecorated := PermanentErrorDecorator(tt.retryable)
			gotErr := fnDecorated()
			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("No error returned")
				} else {
					failErr, ok := gotErr.(failure)
					if !ok {
						t.Errorf("Returned error does not implement failure interface: %v", gotErr)
					}
					isPermanent := failErr.IsPermanent()
					if !isPermanent {
						t.Errorf("Returned error is not permanent: %v", gotErr)
					}
				}
			} else if gotErr != nil {
				t.Errorf("Nil expected, returned error: %v", gotErr)
			}
		})
	}
}

func TestPermanentErrorInterceptor(t *testing.T) {
	type args struct {
		fn Retryable
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PermanentErrorInterceptor(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("PermanentErrorInterceptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermanentError_Error(t *testing.T) {
	originalErr := errors.New("original error")
	err := &PermanentError{
		Cause: originalErr,
	}
	assert.Equal(t, originalErr.Error(), err.Error())
}

func TestPermanentError_IsPermanent(t *testing.T) {
	err := &PermanentError{}
	assert.True(t, err.IsPermanent())
}

func TestPermanentError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := &PermanentError{
		Cause: originalErr,
	}
	assert.Equal(t, originalErr, err.Unwrap())
}

func TestTransientErrorDecorator(t *testing.T) {
	fnOk := func() error { return nil }
	fnErr := func() error { return errors.New("some error") }

	tests := []struct {
		name      string
		retryable Retryable
		wantErr   bool
	}{
		{
			name:      "Success",
			retryable: fnOk,
			wantErr:   false,
		},
		{
			name:      "TransientError",
			retryable: fnErr,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fnDecorated := TransientErrorDecorator(tt.retryable)
			gotErr := fnDecorated()
			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("No error returned")
				} else {
					failErr, ok := gotErr.(failure)
					if !ok {
						t.Errorf("Returned error does not implement failure interface: %v", gotErr)
					}
					isTransient := !failErr.IsPermanent()
					if !isTransient {
						t.Errorf("Returned error is not transient: %v", gotErr)
					}
				}
			} else if gotErr != nil {
				t.Errorf("Nil expected, returned error: %v", gotErr)
			}
		})
	}
}

func TestTransientErrorInterceptor(t *testing.T) {
	type args struct {
		fn Retryable
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TransientErrorInterceptor(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("TransientErrorInterceptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransientError_Error(t *testing.T) {
	originalErr := errors.New("original error")
	err := &TransientError{
		Cause: originalErr,
	}
	assert.Equal(t, originalErr.Error(), err.Error())
}

func TestTransientError_IsPermanent(t *testing.T) {
	err := &TransientError{}
	assert.False(t, err.IsPermanent())
}

func TestTransientError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := &TransientError{
		Cause: originalErr,
	}
	assert.Equal(t, originalErr, err.Unwrap())
}

func TestImplementations(t *testing.T) {
	var _ failure = new(TransientError)
	var _ failure = new(PermanentError)
}
