// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"net/http"
	"time"
)

func tracingFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	now := time.Now()
	ctx := req.Request.Context()
	operationName := RouteOperationFromContext(ctx)

	var opts []trace.StartSpanOption

	// Grab the incoming trace
	wireContext, err := trace.HttpHeadersCarrier(req.Request.Header).Extract()
	if err == nil {
		opts = append(opts,
			trace.StartWithTag(trace.FieldSpanKind, trace.SpanKindServer),
			trace.StartWithRelated(trace.RefChildOf, wireContext))
	}
	opts = append(opts,
		trace.StartWithTag(trace.FieldOperation, operationName),
		trace.StartWithTag(trace.FieldHttpMethod, req.Request.Method),
		trace.StartWithTag(trace.FieldHttpUrl, req.Request.URL.Path),
		trace.StartWithTag(trace.FieldSpanType, "web"))

	ctx, span := trace.NewSpan(ctx, operationName, opts...)
	defer span.Finish()
	req.Request = req.Request.WithContext(ctx)

	chain.ProcessFilter(req, resp)

	logContext := log.LogContext{
		"operation": operationName,
		"method":    req.Request.Method,
		"path":      req.Request.URL.Path,
		"code":      resp.StatusCode(),
		"period":    time.Now().Sub(now).String(),
	}

	traceContext, _ := log.LogContextFromContext(req.Request.Context())
	for k, v := range traceContext {
		logContext[k] = v
	}
	if traceContext == nil {
		traceContext = make(log.LogContext)
	}

	span.LogFields(trace.Int(trace.FieldHttpCode, resp.StatusCode()))

	err = resp.Error()
	if err == nil {
		errInterface := req.Attribute(AttributeError)
		if errInterface != nil {
			err = errInterface.(error)
		}
	}

	if err != nil {
		span.LogFields(trace.Error(err))

		bt := types.BackTraceFromError(err)
		logger.
			WithLogContext(logContext).
			WithError(err).
			WithFields(bt.LogFields()).
			Errorf("Incoming request failed: %s: %s", http.StatusText(resp.StatusCode()), err.Error())
		log.Stack(logger, ctx, bt)
	} else if resp.StatusCode() < 399 {
		var silenced = false
		silencedAttribute := req.Attribute(AttributeSilenceLog)
		if silencedAttributeValue, ok := silencedAttribute.(bool); ok {
			silenced = silencedAttributeValue
		}
		if !silenced {
			logger.WithLogContext(logContext).Infof("Incoming request succeeded: %s", http.StatusText(resp.StatusCode()))
		}
	} else {
		logger.WithLogContext(logContext).Errorf("Incoming request failed: %s", http.StatusText(resp.StatusCode()))
	}
}
