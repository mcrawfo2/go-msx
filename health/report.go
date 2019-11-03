package health

import "context"

type CheckResult struct {
	Status  Status                 `json:"status"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type Report struct {
	Status  Status                 `json:"status"`
	Details map[string]CheckResult `json:"details,omitempty"`
}

func (r *Report) Down() int {
	count := 0
	for _, v := range r.Details {
		if v.Status == StatusDown {
			count++
		}
	}
	return count
}

func (r *Report) Up() int {
	count := 0
	for _, v := range r.Details {
		if v.Status == StatusUp {
			count++
		}
	}
	return count
}

func (r *Report) Unknown() int {
	count := 0
	for _, v := range r.Details {
		if v.Status == StatusUnknown {
			count++
		}
	}
	return count
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

	sendReportStats(report)

	return report
}

func GenerateSummary(ctx context.Context) *Report {
	report := GenerateReport(ctx)
	return &Report{
		Status:  report.Status,
		Details: nil,
	}
}
