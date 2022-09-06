package internal

import (
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "Starts the rotabot server and its dependencies",
		Action: func(c *cli.Context) error {
			logger := shell.Logger(c.Context).With(zap.String("component", "rotabot.server"))
			ctx := shell.WithLogger(c.Context, logger)

			fx.New(Module(ctx), shell.FxEvent()).Run()

			return nil
		},
	}
}
