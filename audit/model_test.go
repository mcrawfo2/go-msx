package audit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testUserName = "monty_burns"

type testAuditable struct {
	CreatedOn time.Time
	CreatedBy string
	UpdatedOn time.Time
	UpdatedBy string
}

func (t *testAuditable) SetCreatedOn(createdOn time.Time) {
	t.CreatedOn = createdOn
}

func (t *testAuditable) SetCreatedBy(createdBy string) {
	t.CreatedBy = createdBy
}

func (t *testAuditable) SetUpdatedOn(updatedOn time.Time) {
	t.UpdatedOn = updatedOn
}

func (t *testAuditable) SetUpdatedBy(updatedBy string) {
	t.UpdatedBy = updatedBy
}

func testModelAuditorContext() context.Context {
	ctx := context.Background()
	ctx = security.ContextWithUserContext(ctx, &security.UserContext{
		UserName: testUserName,
	})
	return ctx
}

func TestModelAuditor_Created(t *testing.T) {
	ctx := testModelAuditorContext()
	tests := []struct {
		name      string
		auditable *testAuditable
	}{
		{
			name:      "Empty",
			auditable: new(testAuditable),
		},
		{
			name: "Existing",
			auditable: &testAuditable{
				CreatedOn: time.Time{},
				CreatedBy: "superuser",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auditor := NewModelAuditor(ctx)
			auditable := tt.auditable
			auditor.Created(auditable)

			assert.NotEqual(t, time.Time{}, auditable.CreatedOn)
			assert.Equal(t, testUserName, auditable.CreatedBy)
			assert.Equal(t, time.Time{}, auditable.UpdatedOn)
			assert.Equal(t, "", auditable.UpdatedBy)
		})
	}
}

func TestModelAuditor_CreatedUpdated(t *testing.T) {
	ctx := testModelAuditorContext()
	tests := []struct {
		name      string
		auditable *testAuditable
	}{
		{
			name:      "Empty",
			auditable: new(testAuditable),
		},
		{
			name: "Existing",
			auditable: &testAuditable{
				CreatedOn: time.Time{},
				CreatedBy: "superuser",
				UpdatedOn: time.Time{},
				UpdatedBy: "superuser",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auditor := NewModelAuditor(ctx)
			auditable := tt.auditable
			auditor.CreatedUpdated(auditable)

			assert.NotEqual(t, time.Time{}, auditable.CreatedOn)
			assert.Equal(t, testUserName, auditable.CreatedBy)
			assert.NotEqual(t, time.Time{}, auditable.UpdatedOn)
			assert.Equal(t, testUserName, auditable.UpdatedBy)
		})
	}
}

func TestModelAuditor_Updated(t *testing.T) {
	ctx := testModelAuditorContext()
	tests := []struct {
		name      string
		auditable *testAuditable
	}{
		{
			name:      "Empty",
			auditable: new(testAuditable),
		},
		{
			name: "Existing",
			auditable: &testAuditable{
				UpdatedOn: time.Time{},
				UpdatedBy: "superuser",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auditor := NewModelAuditor(ctx)
			auditable := tt.auditable
			auditor.Updated(auditable)

			assert.Equal(t, time.Time{}, auditable.CreatedOn)
			assert.Equal(t, "", auditable.CreatedBy)
			assert.NotEqual(t, time.Time{}, auditable.UpdatedOn)
			assert.Equal(t, testUserName, auditable.UpdatedBy)
		})
	}
}

func TestNewModelAuditor(t *testing.T) {
	ctx := testModelAuditorContext()
	auditor := NewModelAuditor(ctx)
	assert.NotNil(t, auditor)
}
