package consul

import (
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"time"
)

const (
	statsCounterConsulCalls            = "consul.calls"
	statsCounterConsulCallErrors       = "consul.callErrors"
	statsTimerConsulCallTime           = "consul.timer"
	statsGaugeConsulRegisteredServices = "consul.registration"

	statsApiListKeyValuePairs   = "list-kv-pairs"
	statsApiGetServiceInstances = "get-service-instances"
	statsApiRegisterService     = "register-service"
	statsApiDeregisterService   = "deregister-service"
	statsApiNodeHealth          = "node-health"
)

type queryFunc func() error

type statsObserver struct{}

func (o *statsObserver) Observe(api, param string, queryFunc queryFunc) (err error) {
	start := time.Now()
	defer func() {
		stats.Incr(stats.Name(statsCounterConsulCalls, api, param), 1)
		stats.PrecisionTiming(stats.Name(statsTimerConsulCallTime, api, param), time.Since(start))
		if err != nil {
			stats.Incr(statsCounterConsulCallErrors, 1)
		}

		switch api {
		case statsApiRegisterService:
			stats.GaugeDelta(statsGaugeConsulRegisteredServices, 1)

		case statsApiDeregisterService:
			stats.GaugeDelta(statsGaugeConsulRegisteredServices, -1)
		}
	}()

	err = queryFunc()
	return err
}
