package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kevinrobayna/rotabot/gen/http/rotabot/server"
	"github.com/kevinrobayna/rotabot/gen/rotabot"
	"github.com/kevinrobayna/rotabot/internal/zapctx"
	"go.uber.org/zap"
	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/middleware"
	goa "goa.design/goa/v3/pkg"
	"net/http"
	"time"
)

type httpOptions struct {
	logger *zap.Logger
}

type decoder struct {
	decoder goahttp.Decoder
	logger  *zap.Logger
	kind    string
}

func (d decoder) Decode(v interface{}) error {
	err := d.decoder.Decode(v)
	if err != nil {
		d.logger.Error(fmt.Sprintf("failed to decode %s", d.kind), zap.Error(err))
	}

	return err
}

type encoder struct {
	encoder goahttp.Encoder
	logger  *zap.Logger
	kind    string
}

func (e encoder) Encode(v interface{}) error {
	err := e.encoder.Encode(v)
	if err != nil {
		e.logger.Error(fmt.Sprintf("failed to decode %s", e.kind), zap.Error(err))
	}

	return err
}

func newHTTPServer(component string, service rotabot.Service, options *httpOptions) *http.Server {
	mux := goahttp.NewMuxer()

	decoder := func(r *http.Request) goahttp.Decoder {
		return &decoder{
			decoder: goahttp.RequestDecoder(r),
			logger:  options.logger.With(zap.String("component", fmt.Sprintf("%s.http", component))),
			kind:    "request",
		}
	}

	encoder := func(ctx context.Context, res http.ResponseWriter) goahttp.Encoder {
		return &encoder{
			encoder: goahttp.ResponseEncoder(ctx, res),
			logger:  zapctx.Logger(ctx),
			kind:    "response",
		}
	}

	endpoints := rotabot.NewEndpoints(service)

	endpoints.Use(func(e goa.Endpoint) goa.Endpoint {
		// Logs endpoint response details
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			timeStart := time.Now()
			l := zapctx.Logger(ctx)
			l.Info("", zap.String("event", "endpoint.started"))

			result, err := e(ctx, req)

			outcome := "success"
			if err != nil {
				outcome = "failure"
			}

			duration := time.Since(timeStart).Seconds()
			l.Info("", zap.String("event", "endpoint.finished"),
				zap.String("outcome", outcome),
				zap.Float64("duration", duration))

			return result, err
		}
	})
	endpoints.Use(func(e goa.Endpoint) goa.Endpoint {
		// Injects the logger into the context with shared info of the request i.e. the request ID
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			rawService := ctx.Value(goa.ServiceKey)
			rawMethod := ctx.Value(goa.MethodKey)
			val := ctx.Value(middleware.RequestIDKey)
			if val == nil {
				val = "unknown"
			}

			id, _ := val.(string)

			l := zapctx.Logger(ctx).With(
				zap.String("handler", fmt.Sprintf("%s:%s", rawService, rawMethod)),
				zap.String("request_id", id))

			return e(zapctx.WithLogger(ctx, l), req)
		}
	})

	errHandler := func(ctx context.Context, w http.ResponseWriter, err error) {
		logger := zapctx.Logger(ctx)
		logger.Error("error handling request", zap.Error(err))
		http.Error(w, fmt.Sprintf("%s: Unexpected error occurred", "http"), http.StatusInternalServerError)
	}

	rotabotServer := server.New(endpoints, mux, decoder, encoder, errHandler, nil)
	server.Mount(mux, rotabotServer)

	for _, m := range rotabotServer.Mounts {
		options.logger.Info("mounts",
			zap.String("verb", m.Verb),
			zap.String("path", m.Pattern),
			zap.String("method", m.Method),
		)
	}

	handler := http.Handler(mux)
	handler = RequestId(handler)

	return &http.Server{
		Handler: handler,
	}
}

func RequestId(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = uuid.New().String()
		}
		ctx = context.WithValue(ctx, middleware.RequestIDKey, id)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
