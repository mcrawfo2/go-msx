package auditlog

import "context"

type contextKey int

const contextKeyRequestAudit contextKey = iota

func ContextWithRequestDetails(ctx context.Context, requestAudit *RequestDetails) context.Context {
	return context.WithValue(ctx, contextKeyRequestAudit, requestAudit)
}

func RequestAuditFromContext(ctx context.Context) *RequestDetails {
	requestAudit, _ := ctx.Value(contextKeyRequestAudit).(*RequestDetails)
	return requestAudit
}
