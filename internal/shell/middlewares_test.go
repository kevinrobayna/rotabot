package shell

import (
	uuidGen "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareAddsUUIDToRequestContext(t *testing.T) {
	req := httptest.NewRequest("GET", "http://testing", nil)

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
	req := httptest.NewRequest("GET", "http://testing", nil)
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

//func TestPanicRecoveryMiddleware(t *testing.T) {
//	ctx := context.Background()
//
//	mockEndpoint := EndpointRecoverMiddleware()(
//		func(ctx context.Context, req interface{}) (interface{}, error) {
//			panic("test")
//		},
//	)
//
//	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
//	observedLogger := zap.New(observedZapCore)
//	ctx = WithLogger(ctx, observedLogger)
//
//	_, err := mockEndpoint(ctx, nil)
//	assert.Error(t, err)
//
//	var expectedError *goa.ServiceError
//	assert.True(t, errors.As(err, &expectedError))
//	assert.NotEmpty(t, expectedError.ID)
//
//	assert.Equal(t, 1, observedLogs.Len())
//
//	loggedEntry := observedLogs.AllUntimed()[0]
//	assert.Equal(t, "request_panic", loggedEntry.Message)
//	assert.Equal(t, 2, len(loggedEntry.Context))
//	assert.Equal(t, "stacktrace", loggedEntry.Context[0].Key)
//}
//
//func TestPanicRecoveryMiddlewareWithError(t *testing.T) {
//	ctx := context.Background()
//
//	mockEndpoint := EndpointRecoverMiddleware()(
//		func(ctx context.Context, req interface{}) (interface{}, error) {
//			panic(errors.New("hello darkness my old friend"))
//		},
//	)
//
//	observedZapCore, _ := observer.New(zap.InfoLevel)
//	observedLogger := zap.New(observedZapCore)
//	ctx = WithLogger(ctx, observedLogger)
//
//	_, err := mockEndpoint(ctx, nil)
//	assert.Error(t, err)
//
//	var expectedError *goa.ServiceError
//	assert.True(t, errors.As(err, &expectedError))
//	assert.NotEmpty(t, expectedError.ID)
//	assert.Equal(t, "unexpected error occurred", expectedError.Message)
//}
