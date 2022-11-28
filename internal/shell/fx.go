package shell

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func FxEvent() fx.Option {
	return fx.WithLogger(func(ctx context.Context) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: Logger(ctx)}
	})
}
