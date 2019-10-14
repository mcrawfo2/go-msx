package health

import "context"

type CheckResult struct {
	Status  Status                 `json:"status"`
	Details map[string]interface{} `json:"details"`
}

type Report struct {
	Status  Status                 `json:"status"`
	Details map[string]CheckResult `json:"details"`
}

func GenerateReport(ctx context.Context) *Report {
	report := &Report{
		Status:  StatusUp,
		Details: make(map[string]CheckResult),
	}

	for name, healthCheck := range healthChecks {
		result := healthCheck(ctx)
		report.Details[name] = result
		report.Status = report.Status.Aggregate(result.Status)
	}

	return report
}

func GenerateSummary(ctx context.Context) *Report {
	report := GenerateReport(ctx)
	return &Report{
		Status:  report.Status,
		Details: nil,
	}
}
