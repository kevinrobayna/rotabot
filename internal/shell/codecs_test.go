package shell

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestEncoderWithLogs_WithoutErrors(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	ctx := WithLogger(
		context.Background(),
		zap.New(observedZapCore),
	)

	res := httptest.NewRecorder()
	rd := ResponseEncoderWithLogs(ctx, res)

	data := map[string]interface{}{"foo": "bar"}
	err := rd.Encode(&data)
	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\"}\n", res.Body.String())
	assert.Equal(t, 0, len(observedLogs.AllUntimed()))
}

type badData struct{}

func (b badData) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot be marshalled to JSON")
}

func TestRequestEncoderWithLogs_WithErrors(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	ctx := WithLogger(
		context.Background(),
		zap.New(observedZapCore),
	)

	res := httptest.NewRecorder()
	rd := ResponseEncoderWithLogs(ctx, res)

	err := rd.Encode(badData{})
	assert.Error(t, err)
	assert.Equal(t, 1, len(observedLogs.AllUntimed()))
	assert.Equal(t, "failed to decode response", observedLogs.AllUntimed()[0].Message)
}

func TestRequestDecoderWithLogs_WithoutErrors(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	ctx := WithLogger(
		context.Background(),
		zap.New(observedZapCore),
	)

	req := httptest.NewRequest(http.MethodGet, "/hello", strings.NewReader(`{"key": "value"}`))
	req = req.WithContext(ctx)
	rd := RequestDecoderWithLogs(req)

	var target map[string]interface{}
	err := rd.Decode(&target)
	assert.NoError(t, err)
	assert.Equal(t, target, map[string]interface{}{"key": "value"})

	assert.Equal(t, 0, len(observedLogs.AllUntimed()))
}

func TestRequestDecoderWithLogs_WithErrors(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	ctx := WithLogger(
		context.Background(),
		zap.New(observedZapCore),
	)

	strings.NewReader("{\"bad")
	req := httptest.NewRequest(http.MethodGet, "/hello", strings.NewReader("{bad"))
	req = req.WithContext(ctx)
	rd := RequestDecoderWithLogs(req)

	var target map[string]interface{}
	err := rd.Decode(&target)
	assert.Error(t, err)

	assert.Equal(t, 1, len(observedLogs.AllUntimed()))
	assert.Equal(t, "failed to decode request", observedLogs.AllUntimed()[0].Message)
}
