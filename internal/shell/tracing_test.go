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

	handlerToTest := HttpRequestIdMiddleware(
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

	handlerToTest := HttpRequestIdMiddleware(
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
