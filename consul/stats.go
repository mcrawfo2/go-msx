package consul

import (
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"time"
)

const (
	statsSubsystemConsul               = "consul"
	statsHistogramConsulCallTime       = "call_time"
	statsGaugeConsulCalls              = "calls"
	statsCounterConsulCallErrors       = "call_errors"
	statsGaugeConsulRegisteredServices = "registrations"

	statsApiListKeyValuePairs      = "list-kv-pairs"
	statsApiGetKeyValue            = "get-kv"
	statsApiSetKeyValue            = "set-kv"
	statsApiGetServiceInstances    = "get-service-instances"
	statsApiGetAllServiceInstances = "get-all-service-instances"
	statsApiRegisterService        = "register-service"
	statsApiDeregisterService      = "deregister-service"
	statsApiNodeHealth             = "node-health"
)

var (
	histVecConsulCallTime    = stats.NewHistogramVec(statsSubsystemConsul, statsHistogramConsulCallTime, nil, "api", "param")
	gaugeVecConsulCalls      = stats.NewGaugeVec(statsSubsystemConsul, statsGaugeConsulCalls, "api", "param")
	countVecConsulCallErrors = stats.NewCounterVec(statsSubsystemConsul, statsCounterConsulCallErrors, "api", "param")
	gaugeConsulRegistrations = stats.NewGauge(statsSubsystemConsul, statsGaugeConsulRegisteredServices)
)

type queryFunc func() error

type statsObserver struct{}

func (o *statsObserver) Observe(api, param string, queryFunc queryFunc) (err error) {
	start := time.Now()
	gaugeVecConsulCalls.WithLabelValues(api, param).Inc()

	defer func() {
		gaugeVecConsulCalls.WithLabelValues(api, param).Dec()
		histVecConsulCallTime.WithLabelValues(api, param).Observe(float64(time.Since(start)) / float64(time.Millisecond))
		if err != nil {
			countVecConsulCallErrors.WithLabelValues(api, param).Inc()
		}

		switch api {
		case statsApiRegisterService:
			gaugeConsulRegistrations.Inc()

		case statsApiDeregisterService:
			gaugeConsulRegistrations.Dec()
		}
	}()

	err = queryFunc()
	return err
}
