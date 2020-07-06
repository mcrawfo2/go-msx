package discovery

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterDiscoveryProvider(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		discoveryProvider = nil
		mockDiscoveryProvider := new(MockDiscoveryProvider)
		RegisterDiscoveryProvider(mockDiscoveryProvider)
		assert.Equal(t, mockDiscoveryProvider, discoveryProvider)
	})

	t.Run("NoRegister", func(t *testing.T) {
		mockDiscoveryProvider := new(MockDiscoveryProvider)
		discoveryProvider = mockDiscoveryProvider
		RegisterDiscoveryProvider(nil)
		assert.Equal(t, mockDiscoveryProvider, discoveryProvider)
	})
}

func TestIsDiscoveryProviderRegistered(t *testing.T) {
	t.Run("Registered", func(t *testing.T) {
		discoveryProvider = nil
		assert.False(t, IsDiscoveryProviderRegistered())
	})

	t.Run("NotRegistered", func(t *testing.T) {
		discoveryProvider = new(MockDiscoveryProvider)
		assert.True(t, IsDiscoveryProviderRegistered())
	})
}

func TestDiscover(t *testing.T) {
	ctx := context.Background()

	t.Run("ProviderNotDefined", func(t *testing.T) {
		discoveryProvider = nil
		instances, err := Discover(ctx, "managedservice", true)
		assert.Error(t, err)
		assert.Equal(t, ErrDiscoveryProviderNotDefined, err)
		assert.Len(t, instances, 0)
	})

	t.Run("Discovered", func(t *testing.T) {
		serviceInstances := mockServiceInstances()
		mockDiscoveryProvider := new(MockDiscoveryProvider)
		mockDiscoveryProvider.
			On("Discover", ctx, "managedservice-1", true).
			Return(
				serviceInstances.Where(func(instance *ServiceInstance) bool {
					return instance.Name == "managedservice-1"
				}),
				nil)
		discoveryProvider = mockDiscoveryProvider

		actualServiceInstances, err := Discover(ctx, "managedservice-1", true)
		assert.NoError(t, err)
		assert.Len(t, actualServiceInstances, 1)
		assert.Equal(t, "managedservice-1", actualServiceInstances[0].Name)
	})
}

func TestDiscoverAll(t *testing.T) {
	ctx := context.Background()

	t.Run("ProviderNotDefined", func(t *testing.T) {
		discoveryProvider = nil
		instances, err := DiscoverAll(ctx, true)
		assert.Error(t, err)
		assert.Equal(t, ErrDiscoveryProviderNotDefined, err)
		assert.Len(t, instances, 0)
	})

	t.Run("Discovered", func(t *testing.T) {
		serviceInstances := mockServiceInstances()
		mockDiscoveryProvider := new(MockDiscoveryProvider)
		mockDiscoveryProvider.
			On("DiscoverAll", ctx, true).
			Return(serviceInstances, nil)
		discoveryProvider = mockDiscoveryProvider

		actualServiceInstances, err := DiscoverAll(ctx, true)
		assert.NoError(t, err)
		assert.Len(t, actualServiceInstances, 3)
	})
}

func TestRegisterRegistrationProvider(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		registrationProvider = nil
		mockRegistrationProvider := new(MockRegistrationProvider)
		RegisterRegistrationProvider(mockRegistrationProvider)
		assert.Equal(t, mockRegistrationProvider, registrationProvider)
	})

	t.Run("NoRegister", func(t *testing.T) {
		mockRegistrationProvider := new(MockRegistrationProvider)
		registrationProvider = mockRegistrationProvider
		RegisterRegistrationProvider(nil)
		assert.Equal(t, mockRegistrationProvider, registrationProvider)
	})
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	t.Run("ProviderNotDefined", func(t *testing.T) {
		registrationProvider = nil
		err := Register(ctx)
		assert.Error(t, err)
		assert.Equal(t, ErrRegistrationProviderNotDefined, err)
	})

	t.Run("Registered", func(t *testing.T) {
		mockRegistrationProvider := new(MockRegistrationProvider)
		mockRegistrationProvider.
			On("Register", ctx).
			Return(nil)

		registrationProvider = mockRegistrationProvider

		err := Register(ctx)
		assert.NoError(t, err)
	})

	t.Run("RegistrationError", func(t *testing.T) {
		mockError := errors.New("Custom error")
		mockRegistrationProvider := new(MockRegistrationProvider)
		mockRegistrationProvider.
			On("Register", ctx).
			Return(mockError)

		registrationProvider = mockRegistrationProvider

		err := Register(ctx)
		assert.Error(t, err)
		assert.Equal(t, mockError, err)
	})
}

func TestDeregister(t *testing.T) {
	ctx := context.Background()
	t.Run("ProviderNotDefined", func(t *testing.T) {
		registrationProvider = nil
		err := Deregister(ctx)
		assert.Error(t, err)
		assert.Equal(t, ErrRegistrationProviderNotDefined, err)
	})

	t.Run("Deregistered", func(t *testing.T) {
		mockRegistrationProvider := new(MockRegistrationProvider)
		mockRegistrationProvider.
			On("Deregister", ctx).
			Return(nil)

		registrationProvider = mockRegistrationProvider

		err := Deregister(ctx)
		assert.NoError(t, err)
	})

	t.Run("DeregistrationError", func(t *testing.T) {
		mockError := errors.New("Custom error")
		mockRegistrationProvider := new(MockRegistrationProvider)
		mockRegistrationProvider.
			On("Deregister", ctx).
			Return(mockError)

		registrationProvider = mockRegistrationProvider

		err := Deregister(ctx)
		assert.Error(t, err)
		assert.Equal(t, mockError, err)
	})
}
