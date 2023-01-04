package internal

import (
	"github.com/kevinrobayna/rotabot/internal/config"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "Starts the rotabot server and its dependencies",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "server.addr",
				Usage:       "Port for the app to listen on",
				DefaultText: ":8080",
			},
		},
		Action: func(c *cli.Context) error {
			logger := shell.Logger(c.Context).With(zap.String("component", "rotabot.server"))
			ctx := shell.WithLogger(c.Context, logger)

			app := fx.New(
				Module(ctx),
				fx.Provide(func() *config.AppConfig { return config.NewAppConfig(c) }),
				shell.FxEvent(),
			)

			app.Run()

			return nil
		},
	}
}
