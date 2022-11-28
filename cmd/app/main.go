package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/kevinrobayna/rotabot/internal"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var (
	AppName = "rotabot"
	Version = "unknown"
	Sha     = "unknown"
	Date    = "unknown"
)

func main() {
	var logger *zap.Logger

	cli.VersionPrinter = func(c *cli.Context) {
		_, _ = fmt.Printf(
			"Application: %s\nVersion: %v\nSha: %s\nGo Version: %v\nGo OS/Arch: %v/%v\nBuilt at: %v\n",
			AppName, Version, Sha, runtime.Version(), runtime.GOOS, runtime.GOARCH, Date,
		)
	}

	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 AppName,
		Usage:                "SlackApp that makes team rotations easy",
		Version:              Version,
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
		},
		Commands: []*cli.Command{
			internal.NewCommand(),
		},
		Before: func(ctx *cli.Context) error {
			cfg := shell.DefaultLoggerConfig()

			if ctx.Bool("verbose") {
				cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
			}

			zl, err := cfg.Build(zap.AddCaller())
			if err != nil {
				return err
			}

			zl = zl.With(
				zap.String("application", AppName),
				zap.String("sha", Sha),
			)

			zap.RedirectStdLog(zl)
			zap.ReplaceGlobals(zl)

			logger = zl
			logger.Info("Logger initialized", zap.String("level", cfg.Level.String()))
			ctx.Context = shell.WithLogger(ctx.Context, logger)
			return nil
		},
	}

	rootCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()

	if err := app.RunContext(rootCtx, os.Args); err != nil {
		if logger != nil {
			logger.Fatal("app errored", zap.Error(err))
		}

		log.Fatalln("app errored: ", err.Error())
	}
}
