package internal

import (
	"context"
	"fmt"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDependenciesAreSatisfied(t *testing.T) {
	ctx := context.Background()
	err := fx.ValidateApp(Module(ctx))
	assert.NoError(t, err)
}

func TestSvc_Healthcheck(t *testing.T) {
	ctx := context.Background()

	var port Port
	app := fxtest.New(t, Module(ctx), fx.Populate(&port))
	defer app.RequireStart().RequireStop()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/healthcheck", port))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, string(body), `{"status":"ok"}`)
}

func TestMiddlewareOrder(t *testing.T) {
	// We want to ensure that the order of the middlewares is the correct one.
	// This ensures that if a panic happens then we get the access logs afterwards
	ctx := context.Background()

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	ctx = shell.WithLogger(ctx, observedLogger.With(zap.String("test", "flag")))

	req := httptest.NewRequest("GET", "http://testing/foo", nil)

	h := wireUpMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			shell.Logger(r.Context()).Info("About to panic!")
			panic("test")
		}),
	)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, req.WithContext(ctx))

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	assert.Equal(t, 4, observedLogs.Len())

	assert.Equal(t, "request.start", observedLogs.AllUntimed()[0].Message)
	// This ensures that our baseContext is injected properly
	assert.Equal(t, "test", observedLogs.AllUntimed()[0].Context[0].Key)
	assert.Equal(t, "flag", observedLogs.AllUntimed()[0].Context[0].String)

	assert.Equal(t, "About to panic!", observedLogs.AllUntimed()[1].Message)
	assert.Equal(t, "request_panic", observedLogs.AllUntimed()[2].Message)
	assert.Equal(t, "request.finish", observedLogs.AllUntimed()[3].Message)
}
