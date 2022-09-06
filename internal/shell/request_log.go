package shell

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	goa "goa.design/goa/v3/pkg"
	"time"
)

// EndpointRequestLogMiddleware waits for the request to complete logging the outcome of it.
func EndpointRequestLogMiddleware() EndpointMiddleware {
	return func(next goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			start := time.Now()
			l := Logger(ctx)

			// Run the next handler
			l.Info("endpoint.start")
			response, err := next(ctx, request)
			outcome := "success"
			if err != nil {
				outcome = "failure"
			}
			l.Info("endpoint.finish",
				zap.String("outcome", outcome),
				zap.Float64("duration", time.Since(start).Seconds()),
			)

			return response, err
		}
	}
}

// EndpointLoggerInjectionMiddleware Injects the logger into the context with shared info of the request
// i.e. the request ID
func EndpointLoggerInjectionMiddleware() EndpointMiddleware {
	return func(next goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			rawService := ctx.Value(goa.ServiceKey)
			rawMethod := ctx.Value(goa.MethodKey)
			val := ctx.Value(RequestIdKey)

			l := Logger(ctx).With(
				zap.String("handler", fmt.Sprintf("%s:%s", rawService, rawMethod)),
				zap.String("request_id", fmt.Sprintf("%s", val)))

			return next(WithLogger(ctx, l), req)
		}
	}
}
