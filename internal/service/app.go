package service

import (
	"context"
	"errors"
	"github.com/kevinrobayna/rotabot/gen/rotabot"
	"github.com/kevinrobayna/rotabot/internal/zapctx"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"
	stdlog "log"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
)

type AppDeps struct {
	Component string
	Service   rotabot.Service
	Context   context.Context
	Listener  net.Listener
}

type App struct {
	group   *run.Group
	ctx     context.Context
	running int32

	Server *http.Server
}

func NewApp(options *AppDeps) *App {
	if strings.TrimSpace(options.Component) == "" {
		panic("No Component name provided to server app")
	}

	if options.Service == nil {
		panic("No Service provided to server app")
	}

	if options.Context == nil {
		panic("No root Context provided to server app")
	}

	if options.Listener == nil {
		panic("No Listener provided to server app")
	}

	var group run.Group
	app := &App{
		group: &group,
		ctx:   options.Context,
	}

	app.Server = initHttp(options, &group)

	return app
}

func (a *App) Run() error {
	swapped := atomic.CompareAndSwapInt32(&a.running, 0, 1)
	if !swapped {
		return errors.New("app is already running")
	}

	a.group.Add(func() error {
		<-a.ctx.Done()
		return nil
	}, func(error) {
		// noop
	})

	return a.group.Run()
}

func initHttp(appOpts *AppDeps, rg *run.Group) *http.Server {
	ctx, cancel := context.WithCancel(appOpts.Context)
	logger := zapctx.Logger(ctx).With(zap.String("component", appOpts.Component))
	ctx = zapctx.WithLogger(ctx, logger)

	srv := newHTTPServer(appOpts.Component, appOpts.Service, &httpOptions{logger: logger})
	srv.BaseContext = func(listener net.Listener) context.Context {
		return ctx
	}
	srv.ErrorLog = stdlog.New(&zapio.Writer{Log: logger, Level: zapcore.ErrorLevel}, "", 0)

	rg.Add(func() error {
		logger.Info("starting_server", zap.Stringer("address", appOpts.Listener.Addr()))
		return srv.Serve(appOpts.Listener)
	}, func(error) {
		logger.Info("stopping_server")
		cancel()

		if err := srv.Close(); err != nil {
			logger.Error("failed_to_close_server", zap.Error(err))
		}
	})

	return srv
}
