package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kevinrobayna/rotabot/internal/config"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type AppSuite struct {
	suite.Suite
}

func NewTestConfig() *config.AppConfig {
	return &config.AppConfig{
		Debug: false,
		Slack: config.SlackConfig{
			SigningSecret: os.Getenv("SLACK_SIGNING_SECRET"),
			ClientSecret:  os.Getenv("SLACK_CLIENT_SECRET"),
		},
	}
}

func (suite *AppSuite) SetupTest() {
}

func (suite *AppSuite) TestDependenciesAreSatisfied() {
	ctx := context.Background()
	err := fx.ValidateApp(Module(ctx), fx.Provide(NewTestConfig))

	suite.NoError(err)
}

func (suite *AppSuite) TestMiddlewareOrder() {
	// We want to ensure that the order of the middlewares is the correct one.
	// This ensures that if a panic happens then we get the access logs afterwards
	ctx := context.Background()

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	ctx = shell.WithLogger(ctx, observedLogger.With(zap.String("test", "flag")))

	req := httptest.NewRequest(http.MethodGet, "http://testing/foo", nil)

	h := wireUpMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			shell.Logger(r.Context()).Info("About to panic!")
			panic("test")
		}),
	)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, req.WithContext(ctx))

	suite.Equal(http.StatusInternalServerError, res.Code)

	suite.Equal(4, observedLogs.Len())

	suite.Equal("request.start", observedLogs.AllUntimed()[0].Message)
	// This ensures that our baseContext is injected properly
	suite.Equal("test", observedLogs.AllUntimed()[0].Context[0].Key)
	suite.Equal("flag", observedLogs.AllUntimed()[0].Context[0].String)

	suite.Equal("About to panic!", observedLogs.AllUntimed()[1].Message)
	suite.Equal("request_panic", observedLogs.AllUntimed()[2].Message)
	suite.Equal("request.finish", observedLogs.AllUntimed()[3].Message)
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}
