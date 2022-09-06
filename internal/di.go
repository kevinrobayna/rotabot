package internal

import (
	"context"
	"errors"
	"github.com/kevinrobayna/rotabot/gen/http/rotabot/server"
	"github.com/kevinrobayna/rotabot/gen/rotabot"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"go.uber.org/fx"
	"go.uber.org/zap"
	goahttp "goa.design/goa/v3/http"
	"net"
	"net/http"
)

func provideListener(ctx context.Context) net.Listener {
	addr := ":8080"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		shell.Logger(ctx).Fatal("failed to listen on addr", zap.String("addr", addr), zap.Error(err))
	}
	return l
}

func provideServerMux(ctx context.Context, service rotabot.Service) goahttp.Muxer {
	mux := goahttp.NewMuxer()

	endpoints := rotabot.NewEndpoints(service)
	endpoints.Use(shell.EndpointRecoverMiddleware())
	endpoints.Use(shell.EndpointRequestLogMiddleware())
	endpoints.Use(shell.EndpointLoggerInjectionMiddleware())

	encoder := shell.ResponseEncoderWithLogs
	decoder := shell.RequestDecoderWithLogs
	rotabotServer := server.New(endpoints, mux, decoder, encoder, shell.ErrorHandler(), nil)
	server.Mount(mux, rotabotServer)

	l := shell.Logger(ctx)
	for _, m := range rotabotServer.Mounts {
		l.Info("mounts",
			zap.String("verb", m.Verb),
			zap.String("path", m.Pattern),
			zap.String("method", m.Method),
		)
	}

	return mux
}

func provideHttpServer(ctx context.Context, mux goahttp.Muxer) *http.Server {
	return shell.NewHttpServer(ctx, mux)
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
