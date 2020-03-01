package auditlog

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"github.com/sirupsen/logrus"
)

type State string

const (
	StateInit    State = "init"
	StateFail    State = "fail"
	StateSuccess State = "success"
	StateAction  State = "action"

	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"

	FieldEntityId = "entityId"
	FieldResource = "resource"
	FieldAction   = "action"
	FieldState    = "state"
	FieldAudit    = "audit"
	FieldSource   = "source"
	FieldProtocol = "protocol"
	FieldHost     = "host"
	FieldPort     = "port"
)

func Entry(logger *log.Logger, ctx context.Context, resourceName, action string, state State) *logrus.Entry {
	entry := logger.
		WithField("resource", resourceName).
		WithField("action", action).
		WithField("state", state).
		WithField("audit", "true")

	requestAudit := RequestAuditFromContext(ctx)
	if requestAudit != nil {
		entry = entry.
			WithField("source", requestAudit.Source).
			WithField("protocol", requestAudit.Protocol).
			WithField("host", requestAudit.Host).
			WithField("port", requestAudit.Port)
	}

	userContext := security.UserContextFromContext(ctx)
	if userContext != nil {
		entry = entry.WithField("user", userContext.UserName)
	}

	return entry.WithContext(ctx)
}

func Init(logger *log.Logger, ctx context.Context, resourceName, action string) *logrus.Entry {
	return Entry(logger, ctx, resourceName, action, StateInit)
}

func Action(logger *log.Logger, ctx context.Context, resourceName, action string) *logrus.Entry {
	return Entry(logger, ctx, resourceName, action, StateAction)
}

func Success(logger *log.Logger, ctx context.Context, resourceName, action string) *logrus.Entry {
	return Entry(logger, ctx, resourceName, action, StateSuccess)
}

func Error(logger *log.Logger, ctx context.Context, resourceName, action string, err error) *logrus.Entry {
	return Failure(logger, ctx, resourceName, action).WithError(err)
}

func Failure(logger *log.Logger, ctx context.Context, resourceName, action string) *logrus.Entry {
	return Entry(logger, ctx, resourceName, action, StateFail)
}

func Result(logger *log.Logger, ctx context.Context, resourceName, action string, err error) *logrus.Entry {
	if err == nil {
		return Success(logger, ctx, resourceName, action)
	} else {
		return Error(logger, ctx, resourceName, action, err)
	}
}

func ResultOf(logger *log.Logger, ctx context.Context, resourceName, action string, fn func() error) *logrus.Entry {
	err := fn()
	return Result(logger, ctx, resourceName, action, err)
}

func Audit(logger *log.Logger, ctx context.Context, resourceName, action string, fn func() error) {
	Init(logger, ctx, resourceName, action).Info()
	ResultOf(logger, ctx, resourceName, action, fn).Info()
}
