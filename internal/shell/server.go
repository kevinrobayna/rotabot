package shell

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"
	goa "goa.design/goa/v3/http"
	stdlog "log"
	"net"
	"net/http"
)

// NewHttpServer creates a new HTTP server with some default middlewares like HttpRequestIdMiddleware.
// You can of course update some config on the http.Server like timeouts or even replace the handler all together.
// This method expects that a http.Muxer with already mounted endpoints and endpoint middlewares.
//
// By default, the server will use zap.Logger through the log.Logger wrapper to log errors.
func NewHttpServer(baseContext context.Context, mux goa.Muxer) *http.Server {
	handler := http.Handler(mux)
	handler = HttpRequestIdMiddleware(handler)

	return &http.Server{
		Handler: handler,
		BaseContext: func(listener net.Listener) context.Context {
			return baseContext
		},
		ErrorLog: NewStdErrorLogger(Logger(baseContext)),
	}
}

func NewStdErrorLogger(l *zap.Logger) *stdlog.Logger {
	return stdlog.New(&zapio.Writer{Log: l, Level: zapcore.ErrorLevel}, "", 0)
}
