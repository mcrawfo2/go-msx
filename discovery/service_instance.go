package discovery

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ServiceInstances []*ServiceInstance

type ServiceInstancePredicate func(*ServiceInstance) bool

func (c ServiceInstances) Where(predicate ServiceInstancePredicate) ServiceInstances {
	result := ServiceInstances{}
	for _, v := range c {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func HasTagPredicate(tag string) ServiceInstancePredicate {
	return func(instance *ServiceInstance) bool {
		return instance.HasTag(tag)
	}
}

func (c ServiceInstances) SelectRandom() *ServiceInstance {
	//return a random member of the result set
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	idx := r1.Intn(len(c))
	return c[idx]
}

type ServiceInstance struct {
	ID   string
	Name string
	Host string
	Tags []string
	Meta map[string]string
	Port int
}

func (i ServiceInstance) Address() string {
	return fmt.Sprintf("%s:%d", i.Host, i.Port)
}

func (i ServiceInstance) HasTag(tag string) bool {
	for _, v := range i.Tags {
		if v == tag {
			return true
		}
	}
	return false
}

func (i ServiceInstance) ContextPath() string {
	for _, v := range i.Tags {
		if strings.HasPrefix(v, "contextPath=") {
			return strings.TrimPrefix(v, "contextPath=")
		}
	}
	return ""
}
