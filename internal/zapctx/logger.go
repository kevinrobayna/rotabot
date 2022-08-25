package zapctx

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// NewLogger returns a new zap logger. This logger is configured to log at InfoLevel and above
// We are using this logger instead of zap.NewProduction because by default it's configured to sample the logs
func NewLogger() *zap.Logger {
	config := NewLoggerConfig()
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			os.Stdout,
			config.Level,
		),
	)
}

func NewLoggerConfig() zap.Config {
	zapconf := zap.NewProductionConfig()
	zapconf.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	zapconf.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	return zapconf
}

// This is required to make the logger available in the ctx or fetch it from the ctx
var loggerKey struct{}

// WithLogger attaches a zap logger to a given context, returning the new updated context
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// Logger extracts a zap logger from a context, providing a default logger if one is not present
func Logger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		panic("nil context passed to Logger")
	}

	if logger, _ := ctx.Value(loggerKey).(*zap.Logger); logger != nil {
		return logger
	}

	return NewLogger()
}
