package background

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockImplements(t *testing.T) {
	// Compile-time test
	var _ ErrorReporter = new(MockErrorReporter)
}

func TestContextWithErrorReporter(t *testing.T) {
	errorReporter := new(MockErrorReporter)
	ctx := ContextWithErrorReporter(context.Background(), errorReporter)
	assert.NotNil(t, ctx)
	assert.Equal(t, errorReporter, ctx.Value(contextKeyErrorReporter))
}

func TestErrorReporterFromContext(t *testing.T) {
	errorReporter := new(MockErrorReporter)
	ctx := ContextWithErrorReporter(context.Background(), errorReporter)
	assert.NotNil(t, ctx)
	actualErrorReporter := ErrorReporterFromContext(ctx)
	assert.Equal(t, errorReporter, actualErrorReporter)
}

func TestErrorReporterFromContext_Nil(t *testing.T) {
	actualErrorReporter := ErrorReporterFromContext(context.Background())
	assert.Nil(t, actualErrorReporter)
}
