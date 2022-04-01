// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
