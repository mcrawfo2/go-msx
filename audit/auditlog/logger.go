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

func Failure(logger *log.Logger, ctx context.Context, resourceName, action string) *logrus.Entry {
	return Entry(logger, ctx, resourceName, action, StateFail)
}

func Result(logger *log.Logger, ctx context.Context, resourceName, action string, err error) *logrus.Entry {
	if err == nil {
		return Success(logger, ctx, resourceName, action)
	} else {
		return Failure(logger, ctx, resourceName, action).WithError(err)
	}
}

func ResultOf(logger *log.Logger, ctx context.Context, resourceName, action string, fn func () error) (*logrus.Entry, error) {
	err := fn()
	return Result(logger, ctx, resourceName, action, err), err
}
