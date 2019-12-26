package audit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"time"
)

const (
	userNameDefault = "system"
)

type CreateAuditable interface {
	SetCreatedOn(time.Time)
	SetCreatedBy(string)
}

type UpdateAuditable interface {
	SetUpdatedOn(time.Time)
	SetUpdatedBy(string)
}

type CreateUpdateAuditable interface {
	CreateAuditable
	UpdateAuditable
}

type ModelAuditor struct {
	ctx context.Context
}

func (a ModelAuditor) Created(auditable CreateAuditable) {
	auditable.SetCreatedBy(a.userName())
	auditable.SetCreatedOn(time.Now().In(time.UTC))
}

func (a ModelAuditor) Updated(auditable UpdateAuditable) {
	auditable.SetUpdatedBy(a.userName())
	auditable.SetUpdatedOn(time.Now().In(time.UTC))
}

func (a ModelAuditor) CreatedUpdated(auditable CreateUpdateAuditable) {
	t := time.Now()
	u := a.userName()
	auditable.SetCreatedBy(u)
	auditable.SetCreatedOn(t)
	auditable.SetUpdatedBy(u)
	auditable.SetUpdatedOn(t)
}

func (a ModelAuditor) userName() string {
	user := security.UserContextFromContext(a.ctx)
	userName := userNameDefault
	if user != nil {
		userName = user.UserName
	}
	return userName
}

func NewModelAuditor(ctx context.Context) ModelAuditor {
	return ModelAuditor{ctx: ctx}
}
