package shell

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	goa "goa.design/goa/v3/pkg"
	"testing"
)

func TestPanicRecoveryMiddleware(t *testing.T) {
	ctx := context.Background()

	mockEndpoint := EndpointRecoverMiddleware()(
		func(ctx context.Context, req interface{}) (interface{}, error) {
			panic("test")
		},
	)

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	ctx = WithLogger(ctx, observedLogger)

	_, err := mockEndpoint(ctx, nil)
	assert.Error(t, err)

	var expectedError *goa.ServiceError
	assert.True(t, errors.As(err, &expectedError))
	assert.NotEmpty(t, expectedError.ID)

	assert.Equal(t, 1, observedLogs.Len())

	loggedEntry := observedLogs.AllUntimed()[0]
	assert.Equal(t, "request_panic", loggedEntry.Message)
	assert.Equal(t, 2, len(loggedEntry.Context))
	assert.Equal(t, "stacktrace", loggedEntry.Context[0].Key)
}

func TestPanicRecoveryMiddlewareWithError(t *testing.T) {
	ctx := context.Background()

	mockEndpoint := EndpointRecoverMiddleware()(
		func(ctx context.Context, req interface{}) (interface{}, error) {
			panic(errors.New("hello darkness my old friend"))
		},
	)

	observedZapCore, _ := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	ctx = WithLogger(ctx, observedLogger)

	_, err := mockEndpoint(ctx, nil)
	assert.Error(t, err)

	var expectedError *goa.ServiceError
	assert.True(t, errors.As(err, &expectedError))
	assert.NotEmpty(t, expectedError.ID)
	assert.Equal(t, "unexpected error occurred", expectedError.Message)
}
