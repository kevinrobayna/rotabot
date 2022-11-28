package internal

import (
	"context"
	"errors"
	stdlog "log"
	"net"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kevinrobayna/rotabot/internal/config"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"
)

func Module(ctx context.Context) fx.Option {
	return fx.Module("rotabot",
		fx.Provide(providePort),
		fx.Provide(provideListener),
		fx.Provide(provideServerRouter),
		fx.Provide(provideHttpServer),
		fx.Provide(func() context.Context { return ctx }),

		fx.Invoke(invokeHttpServer),
	)
}

type Port string

func providePort(listener net.Listener) Port {
	addr := strings.TrimPrefix(listener.Addr().String(), "[::]:")
	return Port(addr)
}

func provideListener(ctx context.Context) net.Listener {
	addr := ":8080"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		shell.Logger(ctx).Fatal("failed to listen on addr", zap.String("addr", addr), zap.Error(err))
	}
	return l
}

func provideServerRouter(cfg *config.AppConfig) *httprouter.Router {
	r := httprouter.New()

	resource := &resource{
		cfg: cfg,
	}

	r.HandlerFunc(http.MethodGet, "/healthcheck", resource.HealthCheck())
	r.HandlerFunc(http.MethodPost, "/slack/commands", resource.HandleSlashCommand())

	return r
}

func provideHttpServer(ctx context.Context, r *httprouter.Router) *http.Server {
	// TODO: add requests time-outs since we don't want to keep connections open forever
	return &http.Server{
		Handler: wireUpMiddlewares(http.Handler(r)),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		ErrorLog: stdlog.New(&zapio.Writer{Log: shell.Logger(ctx), Level: zapcore.ErrorLevel}, "", 0),
	}
}

func wireUpMiddlewares(h http.Handler) http.Handler {
	h = shell.RecoveryHandler(h)
	h = shell.RequestAccessLogHandler(h)
	h = shell.LoggerInjectionHandler(h)
	h = shell.RequestIdHandler(h)
	return h
}

func invokeHttpServer(lc fx.Lifecycle, server *http.Server, listener net.Listener) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l := shell.Logger(ctx)
			l.Info("running server", zap.String("addr", listener.Addr().String()))
			go func() {
				if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
					l.Fatal("failed to serve from server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
