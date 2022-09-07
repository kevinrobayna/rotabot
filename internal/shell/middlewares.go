package shell

import (
	"context"
	"fmt"
	uuidGen "github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RecoveryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := Logger(r.Context())

		defer func(logger *zap.Logger) {
			if rawErr := recover(); rawErr != nil {
				var err error
				switch e := rawErr.(type) {
				case error:
					err = e
				default:
					err = fmt.Errorf("panic: %v", rawErr)
				}

				w.WriteHeader(http.StatusInternalServerError)
				l.Error("request_panic", zap.Stack("stacktrace"), zap.Error(err))
				return
			}
		}(l)
		next.ServeHTTP(w, r)
	})
}

// RequestIdHandler is a middleware that adds a unique request ID to the request context.
func RequestIdHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get(string(RequestIdKey))
		if id == "" {
			id = uuid()
		}
		w.Header().Set(string(RequestIdKey), id)
		ctx = context.WithValue(ctx, RequestIdKey, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func uuid() string {
	return uuidGen.New().String()
}

// RequestLogHandler waits for the request to complete logging the outcome of it.
func RequestLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()
		l := Logger(ctx)

		// Run the next handler
		l.Info("endpoint.start")
		next.ServeHTTP(w, r)
		l.Info("endpoint.finish",
			zap.Float64("duration", time.Since(start).Seconds()),
		)
	})
}

// LoggerInjectionHandler Injects the logger into the context with shared info of the request
// i.e. the request ID
func LoggerInjectionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := Logger(ctx).With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.EscapedPath()),
			zap.String("request_id", fmt.Sprintf("%s", ctx.Value(RequestIdKey))),
		)
		next.ServeHTTP(w, r.WithContext(WithLogger(ctx, l)))
	})
}
