package discovery

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockServiceInstance() ServiceInstance {
	return ServiceInstance{
		ID:   types.EmptyUUID().String(),
		Name: "managedservice",
		Host: "127.0.0.1",
		Port: 8080,
		Tags: []string{"tag1", "tag2"},
		Meta: map[string]string{"key1": "value1", "key2": "value2"},
	}
}

func TestServiceInstance_Address(t *testing.T) {
	address := mockServiceInstance().Address()
	assert.Equal(t, "127.0.0.1:8080", address)
}

func TestServiceInstance_HasTag(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		assert.True(t, mockServiceInstance().HasTag("tag1"))
	})
	t.Run("False", func(t *testing.T) {
		assert.False(t, mockServiceInstance().HasTag("tag3"))
	})
}

func TestServiceInstance_ContextPath(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		serviceInstance := mockServiceInstance()
		serviceInstance.Tags = append(serviceInstance.Tags, "contextPath=/managed")
		assert.Equal(t, "/managed", serviceInstance.ContextPath())
	})
	t.Run("Not Exists", func(t *testing.T) {
		serviceInstance := mockServiceInstance()
		assert.Equal(t, "", serviceInstance.ContextPath())
	})
}

func mockServiceInstances() ServiceInstances {
	serviceInstance1 := mockServiceInstance()
	serviceInstance1.Tags = append(serviceInstance1.Tags, "tag3")
	serviceInstance1.Name = "managedservice-1"
	serviceInstance2 := mockServiceInstance()
	serviceInstance2.Name = "managedservice-2"
	serviceInstance3 := mockServiceInstance()
	serviceInstance3.Name = "managedservice-3"

	return ServiceInstances{
		&serviceInstance1,
		&serviceInstance2,
		&serviceInstance3,
	}
}

func TestServiceInstances_Where(t *testing.T) {
	serviceInstances := mockServiceInstances()

	t.Run("Matching", func(t *testing.T) {
		predicate := func(s *ServiceInstance) bool {
			return s.Name == "managedservice-3"
		}

		whereServiceInstances := serviceInstances.Where(predicate)
		assert.Len(t, whereServiceInstances, 1)
		assert.Equal(t, whereServiceInstances[0].Name, "managedservice-3")
	})

	t.Run("Non-Matching", func(t *testing.T) {
		predicateNever := func(s *ServiceInstance) bool {
			return false
		}

		whereServiceInstances := serviceInstances.Where(predicateNever)
		assert.Len(t, whereServiceInstances, 0)
	})
}

func TestServiceInstances_SelectRandom(t *testing.T) {
	serviceInstances := mockServiceInstances()
	serviceInstance := serviceInstances.SelectRandom()
	assert.NotNil(t, serviceInstance)
}

func TestHasTagPredicate(t *testing.T) {
	serviceInstances := mockServiceInstances()
	whereServiceInstances := serviceInstances.Where(HasTagPredicate("tag3"))
	assert.Len(t, whereServiceInstances, 1)
	assert.Equal(t, whereServiceInstances[0].Name, "managedservice-1")
}
