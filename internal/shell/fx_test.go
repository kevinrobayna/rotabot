package shell

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestFxEvent(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	ctx := WithLogger(
		context.Background(),
		zap.New(observedZapCore),
	)

	app := fxtest.New(
		t,
		fx.Provide(func() context.Context { return ctx }),
		FxEvent(),
	)

	defer app.RequireStart().RequireStop()

	assert.Equal(t, 6, len(observedLogs.AllUntimed()))
}
