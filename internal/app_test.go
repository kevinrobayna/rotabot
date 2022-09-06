package internal

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"io"
	"net/http"
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

	_, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)

}
