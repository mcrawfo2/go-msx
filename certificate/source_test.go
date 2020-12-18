package certificate

import (
	"context"
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/background"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"time"
)

func TestSource_Certificate(t *testing.T) {
	cert := &tls.Certificate{}

	tests := []struct {
		name string
		src  *Source
		want *tls.Certificate
	}{
		{
			name: "HasCertificate",
			src: &Source{
				certificate: cert,
			},
			want: cert,
		},
		{
			name: "NoCertificate",
			src:  &Source{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.src.Certificate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Certificate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_TlsCertificate(t *testing.T) {
	cert := &tls.Certificate{}

	tests := []struct {
		name    string
		src     *Source
		want    *tls.Certificate
		wantErr error
	}{
		{
			name: "HasCertificate",
			src: &Source{
				certificate: cert,
			},
			want:    cert,
			wantErr: nil,
		},
		{
			name:    "NoCertificate",
			src:     &Source{},
			want:    nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, gotErr := tt.src.TlsCertificate(nil); gotErr != tt.wantErr {
				t.Errorf("TlsCertificate() err = %v, want %v", gotErr, tt.wantErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TlsCertificate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_setCertificate(t *testing.T) {
	cert := &tls.Certificate{}
	type args struct {
		certificate *tls.Certificate
		src         *Source
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "HasCertificate",
			args: args{
				certificate: cert,
				src:         &Source{},
			},
		},
		{
			name: "NoCertificate",
			args: args{
				src: &Source{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.src
			c.setCertificate(tt.args.certificate)
			assert.Equal(t, tt.args.certificate, c.certificate)
		})
	}
}

func TestSource_period(t *testing.T) {
	clock := types.NewMockClock()

	tests := []struct {
		name    string
		src     *Source
		wantMin time.Duration
		wantMax time.Duration
		wantErr error
	}{
		{
			name: "ShortCertificate",
			src: &Source{
				certificate: generateCertificate(t, clock, 10*time.Minute),
				clock:       clock,
			},
			wantMin: 5*time.Minute + 150*time.Second,
			wantMax: 5*time.Minute + 225*time.Second,
			wantErr: nil,
		},
		{
			name: "LongCertificate",
			src: &Source{
				certificate: generateCertificate(t, clock, 30*time.Minute),
				clock:       clock,
			},
			wantMin: 15*time.Minute,
			wantMax: 30*time.Minute,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.src
			got, gotErr := c.period()
			if gotErr != tt.wantErr {
				t.Errorf("Expected error %v; got %v", tt.wantErr, gotErr)
			}
			if tt.wantMin > got {
				t.Errorf("Wanted min %f seconds; got %f seconds", tt.wantMin.Seconds(), got.Seconds())
			}
			if tt.wantMax < got {
				t.Errorf("Wanted max %f seconds; got %f seconds", tt.wantMax.Seconds(), got.Seconds())
			}
		})
	}
}

func TestSource_renewOnce_Success(t *testing.T) {
	clock := types.NewMockClock()
	cert := generateCertificate(t, clock, 10*time.Minute)
	provider := new(mockProvider)
	provider.On("GetCertificate", mock.AnythingOfType("*context.emptyCtx")).Return(cert, nil)

	src := &Source{
		provider: provider,
		clock:    clock,
	}

	err := src.renewOnce(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, cert, src.certificate)
}

func TestSource_renewOnce_Error(t *testing.T) {
	clock := types.NewMockClock()
	errResult := errors.New("some error")
	provider := new(mockProvider)
	provider.On("GetCertificate", mock.AnythingOfType("*context.emptyCtx")).Return(nil, errResult)

	src := &Source{
		provider: provider,
		clock:    clock,
	}

	err := src.renewOnce(context.Background())
	assert.Error(t, err)
	assert.Nil(t, src.certificate)
}

func TestSource_renew(t *testing.T) {
	clock := types.NewMockClock()
	certificate := generateCertificate(t, clock, 10*time.Minute)

	deadlineCtx, cancelDeadlineCtx := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer func() {
		if deadlineCtx.Err() != nil {
			t.Error("GetCertificate was not called")
		}
		cancelDeadlineCtx()
	}()

	ctx, cancelCtx := context.WithCancel(deadlineCtx)
	ctx = types.ContextWithClock(ctx, clock)

	provider := new(mockProvider)
	provider.
		On("GetCertificate", mock.AnythingOfType("*context.valueCtx")).
		Return(certificate, nil).
		Run(func(args mock.Arguments) {
			// After retrieving the certificate, cancel the renew loop
			cancelCtx()
		})

	src := &Source{
		certificate: certificate,
		provider:    provider,
		clock:       clock,
	}

	go func() {
		// Trigger the timer
		clock.Advance(9 * time.Minute)
		clock.Trigger(renewTimerId)
	}()

	src.renew(ctx)
}

func TestSource_renewFailure(t *testing.T) {
	clock := types.NewMockClock()
	certificate := generateCertificate(t, clock, 10*time.Minute)

	deadlineCtx, cancelDeadlineCtx := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer func() {
		if deadlineCtx.Err() != nil {
			t.Error("GetCertificate was not called")
		}
		cancelDeadlineCtx()
	}()

	ctx, cancelCtx := context.WithCancel(deadlineCtx)
	ctx = types.ContextWithClock(ctx, clock)

	errorReporter := new(background.MockErrorReporter)
	errorReporter.On("Fatal", mock.AnythingOfType("*errors.withStack")).Return()
	ctx = background.ContextWithErrorReporter(ctx, errorReporter)

	errRenewalFailed := errors.New("renewal failure")

	provider := new(mockProvider)
	provider.
		On("GetCertificate", mock.AnythingOfType("*context.valueCtx")).
		Return(nil, errRenewalFailed).
		Run(func(args mock.Arguments) {
			// After retrieving the certificate, cancel the renew loop
			cancelCtx()
		})

	src := &Source{
		certificate: certificate,
		provider:    provider,
		clock:       clock,
	}

	go func() {
		// Trigger the timer
		clock.Advance(9 * time.Minute)
		clock.Trigger(renewTimerId)
	}()

	src.renew(ctx)

	mock.AssertExpectationsForObjects(t, errorReporter, provider)
}

func TestNewSource(t *testing.T) {
	sources = make(map[string]*Source)
	factories = make(map[string]ProviderFactory)

	ctx := configtest.ContextWithNewStaticConfig(
		context.Background(),
		map[string]string{
			"certificate.source.test.provider": "mock",
		})

	clock := types.NewMockClock()
	ctx = types.ContextWithClock(ctx, clock)

	certificate := generateCertificate(t, clock, 10*time.Minute)

	provider := new(mockProvider)
	provider.
		On("GetCertificate", mock.AnythingOfType("*context.valueCtx")).
		Return(certificate, nil)
	provider.On("Renewable").Return(false)

	factory := new(mockProviderFactory)
	factory.On("Name").Return("mock")
	factory.On("New", mock.AnythingOfType("*context.valueCtx"), "certificate.source.test").
		Return(provider, nil)
	RegisterProviderFactory(factory)

	source, err := NewSource(ctx, "test")
	assert.NoError(t, err)
	assert.NotNil(t, source)
}
