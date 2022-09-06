package shell

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	goa "goa.design/goa/v3/pkg"
)

func EndpointRecoverMiddleware() EndpointMiddleware {
	return func(next goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, req interface{}) (result interface{}, err error) {
			l := Logger(ctx)

			defer func(logger *zap.Logger) {
				if rawErr := recover(); rawErr != nil {
					internalError := NewInternalError()
					err = internalError

					var err error
					switch e := rawErr.(type) {
					case error:
						err = e
					default:
						err = fmt.Errorf("panic: %v", rawErr)
					}

					l.Error("request_panic", zap.Stack("stacktrace"), zap.Error(err))
					return
				}
			}(l)

			result, err = next(ctx, req)
			return
		}
	}
}
