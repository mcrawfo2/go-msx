package health

import "context"

type Check func(context.Context) CheckResult

var healthChecks = make(map[string]Check)

func RegisterCheck(name string, check Check) {
	healthChecks[name] = check
}
