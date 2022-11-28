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
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "This will enable all possible logging, this means that sensitive information will be logged, so  avoid using this in production",
				Value: false,
			},
			&cli.StringFlag{
				Name:     "slack.signing_secret",
				Usage:    "Secret that ensures the requests from slack are real",
				Required: true,
				EnvVars:  []string{"SLACK_SIGNING_SECRET"},
			},
			&cli.StringFlag{
				Name:     "slack.client_secret",
				Usage:    "Secret that allows app to access the slack API, format: `xoxb-*`",
				Required: true,
				EnvVars:  []string{"SLACK_CLIENT_SECRET"},
			},
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
