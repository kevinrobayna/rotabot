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
		fmt.Printf(
			"Application: %s\nVersion: %v\nSha: %s\nGo Version: %v\nGo OS/Arch: %v/%v\nBuilt at: %v\n",
			AppName, Version, Sha, runtime.Version(), runtime.GOOS, runtime.GOARCH, Date,
		)
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print detailed version information",
	}

	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 AppName,
		Usage:                "SlackApp that makes team rotations easy",
		Version:              Version,
		Commands: []*cli.Command{
			internal.NewCommand(),
		},
		Before: func(ctx *cli.Context) error {
			zl, err := shell.DefaultLoggerConfig().Build(zap.AddCaller())
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
