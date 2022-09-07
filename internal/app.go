package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"
	stdlog "log"
	"net"
	"net/http"
	"strings"
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

func provideServerRouter(ctx context.Context) *httprouter.Router {
	r := httprouter.New()

	//TODO this needs to be moved out of here
	r.GET("/healthcheck", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		shell.Logger(r.Context()).Info("Healthcheck")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	return r
}

func provideHttpServer(ctx context.Context, r *httprouter.Router) *http.Server {
	handler := http.Handler(r)
	handler = shell.RecoveryHandler(handler)
	handler = shell.RequestLogHandler(handler)
	handler = shell.LoggerInjectionHandler(handler)
	handler = shell.RequestIdHandler(handler)

	//TODO: add requests time outs since we don't want to keep connections open forever
	return &http.Server{
		Handler: handler,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		ErrorLog: stdlog.New(&zapio.Writer{Log: shell.Logger(ctx), Level: zapcore.ErrorLevel}, "", 0),
	}
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
