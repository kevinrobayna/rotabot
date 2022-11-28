package shell

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerReturnsPanicWhenContextIsNil(t *testing.T) {
	assert.Panics(
		t,
		func() {
			Logger(nil) //nolint:staticcheck // This is a test, we want to panic
		},
		"nil context passed to Logger",
	)
}

func TestContextContainsLoggerWithValue(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	ctx := WithLogger(
		context.Background(),
		zap.New(observedZapCore).
			With(zap.String("key", "value")),
	)

	l := Logger(ctx)

	l.Info("Send Random Test")

	loggedEntry := observedLogs.AllUntimed()[0]
	assert.Equal(t, 1, len(loggedEntry.Context))
	assert.Equal(t, "key", loggedEntry.Context[0].Key)
	assert.Equal(t, "value", loggedEntry.Context[0].String)
}

func TestDefaultLoggerIsReturned(t *testing.T) {
	l := Logger(context.Background())

	assert.NotNil(t, l)
	assert.True(t, l.Core().Enabled(zapcore.InfoLevel))
}
