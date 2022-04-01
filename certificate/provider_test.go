// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package certificate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestMockProviderFactoryImplementation(t *testing.T) {
	var _ ProviderFactory = new(mockProviderFactory)
}

func TestMockProviderImplementation(t *testing.T) {
	var _ Provider = new(mockProvider)
}

func TestRegisterProviderFactory(t *testing.T) {
	factory := new(mockProviderFactory)
	factory.On("Name").Return("mock")

	RegisterProviderFactory(factory)

	actual, ok := factories["mock"]

	assert.True(t, ok)
	assert.Equal(t, factory, actual)
	mock.AssertExpectationsForObjects(t, factory)
}

func TestNewProvider_Success(t *testing.T) {
	provider := new(mockProvider)
	factory := new(mockProviderFactory)
	factory.On("Name").Return("mock")
	factory.On("New", mock.AnythingOfType("*context.valueCtx"), "certificate.source.test").Return(provider, nil)

	ctx := configtest.ContextWithNewInMemoryConfig(
		context.Background(),
		map[string]string{
			"certificate.source.test.provider": "mock",
		})

	// Clean factory list
	factories = make(map[string]ProviderFactory)
	RegisterProviderFactory(factory)

	actual, err := NewProvider(ctx, "test")
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, provider, actual)
	mock.AssertExpectationsForObjects(t, provider, factory)
}

func TestNewProvider_SourceNotConfigured(t *testing.T) {
	cfg := config.NewConfig()
	err := cfg.Load(context.Background())
	assert.NoError(t, err)

	ctx := config.ContextWithConfig(context.Background(), cfg)

	// Clean factory list
	factories = make(map[string]ProviderFactory)

	actual, err := NewProvider(ctx, "test")
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestNewProvider_NoSuchProvider(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(
		context.Background(),
		map[string]string{
			"certificate.source.test.provider": "mock",
		})

	// Clean factory list
	factories = make(map[string]ProviderFactory)

	actual, err := NewProvider(ctx, "test")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNoSuchProvider))
	assert.Nil(t, actual)
}
