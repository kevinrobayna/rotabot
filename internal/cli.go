package internal

import (
	"context"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"go.uber.org/fx/fxevent"

	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "web",
		Usage:  "Starts the rotabot server and its dependencies",
		Action: serviceFunc(),
	}
}

func serviceFunc() cli.ActionFunc {
	return func(c *cli.Context) error {
		ctx := c.Context
		logger := shell.Logger(ctx).With(zap.String("component", "rotabot.server"))
		ctx = shell.WithLogger(ctx, logger)

		app := fx.New(
			fx.Provide(provideListener),
			fx.Provide(provideServerMux),
			fx.Provide(provideHttpServer),
			fx.Provide(NewRotabotService),
			fx.Provide(func() context.Context { return ctx }),

			fx.Invoke(invokeHttpServer),
			fx.WithLogger(func(ctx context.Context) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: shell.Logger(ctx)}
			}),
		)

		app.Run()
		return nil
	}
}
