package shell

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	uuidGen "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestMiddlewareAddsUUIDToRequestContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

	handlerToTest := RequestIdHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NotEmpty(t, w.Header().Get(string(RequestIdKey)))
			_, err := uuidGen.Parse(w.Header().Get(string(RequestIdKey)))
			assert.NoError(t, err)
			assert.NotEmpty(t, r.Context().Value(RequestIdKey))

			assert.Equal(t, w.Header().Get(string(RequestIdKey)), r.Context().Value(RequestIdKey))
		}),
	)

	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}

func TestMiddlewareReusesUUIDIfPresentOnRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
	req.Header.Set(string(RequestIdKey), "123")

	handlerToTest := RequestIdHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NotEmpty(t, w.Header().Get(string(RequestIdKey)))
			_, err := uuidGen.Parse(w.Header().Get(string(RequestIdKey)))
			assert.Error(t, err)

			assert.Equal(t, "123", w.Header().Get(string(RequestIdKey)))

			assert.NotEmpty(t, r.Context().Value(RequestIdKey))
			assert.Equal(t, w.Header().Get(string(RequestIdKey)), r.Context().Value(RequestIdKey))
		}),
	)

	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}

func TestPanicRecoveryMiddleware(t *testing.T) {
	ctx := context.Background()

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	ctx = WithLogger(ctx, observedLogger)

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

	handlerToTest := RecoveryHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test")
		}),
	)

	res := httptest.NewRecorder()
	handlerToTest.ServeHTTP(res, req.WithContext(ctx))

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	assert.Equal(t, 1, observedLogs.Len())

	loggedEntry := observedLogs.AllUntimed()[0]
	assert.Equal(t, "request_panic", loggedEntry.Message)
	assert.Equal(t, 2, len(loggedEntry.Context))
	assert.Equal(t, "stacktrace", loggedEntry.Context[0].Key)
}

func TestPanicRecoveryMiddlewareWithError(t *testing.T) {
	ctx := context.Background()

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	ctx = WithLogger(ctx, observedLogger)

	req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

	handlerToTest := RecoveryHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("hello darkness my old friend"))
		}),
	)

	res := httptest.NewRecorder()
	handlerToTest.ServeHTTP(res, req.WithContext(ctx))

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	assert.Equal(t, 1, observedLogs.Len())

	loggedEntry := observedLogs.AllUntimed()[0]
	assert.Equal(t, "request_panic", loggedEntry.Message)
	assert.Equal(t, 2, len(loggedEntry.Context))
	assert.Equal(t, "stacktrace", loggedEntry.Context[0].Key)
}

func TestRequestLogHandler(t *testing.T) {
	ctx := context.Background()

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	ctx = WithLogger(ctx, observedLogger)

	req := httptest.NewRequest(http.MethodGet, "http://testing/foo", nil)

	handlerToTest := RequestAccessLogHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}),
	)

	handlerToTest.ServeHTTP(httptest.NewRecorder(), req.WithContext(ctx))

	assert.Equal(t, 2, observedLogs.Len())

	assert.Equal(t, "request.start", observedLogs.AllUntimed()[0].Message)
	assert.Equal(t, "request.finish", observedLogs.AllUntimed()[1].Message)

	loggedEntry := observedLogs.AllUntimed()[1]
	assert.Equal(t, 2, len(loggedEntry.Context))
	assert.Equal(t, "duration", loggedEntry.Context[0].Key)
	assert.Equal(t, "status", loggedEntry.Context[1].Key)
	assert.Equal(t, int64(http.StatusBadRequest), loggedEntry.Context[1].Integer)
}

func TestLoggerInjectionHandler(t *testing.T) {
	ctx := context.Background()

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	ctx = WithLogger(ctx, observedLogger)

	req := httptest.NewRequest(http.MethodGet, "http://testing/foo", nil)

	handlerToTest := LoggerInjectionHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Logger(r.Context()).Info("hello")
		}),
	)

	handlerToTest.ServeHTTP(httptest.NewRecorder(), req.WithContext(ctx))

	assert.Equal(t, 1, observedLogs.Len())

	loggedEntry := observedLogs.AllUntimed()[0]
	assert.Equal(t, "hello", loggedEntry.Message)

	assert.Equal(t, 3, len(loggedEntry.Context))
	assert.Equal(t, "method", loggedEntry.Context[0].Key)
	assert.Equal(t, "path", loggedEntry.Context[1].Key)
	assert.Equal(t, "request_id", loggedEntry.Context[2].Key)
}
