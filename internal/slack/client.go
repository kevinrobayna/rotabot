package slack

import (
	"context"
	"github.com/kevinrobayna/rotabot/internal/config"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type zapWrap struct {
	l *zap.Logger
}

func (z zapWrap) Output(_ int, message string) error {
	z.l.Info(message)
	return nil
}

func slackClient(cfg config.AppConfig, ctx context.Context) (*slack.Client, error) {
	c := slack.New(
		cfg.Slack.ClientSecret,
		slack.OptionLog(zapWrap{l: shell.Logger(ctx)}),
		slack.OptionDebug(cfg.Debug),
	)

	if _, err := c.AuthTestContext(ctx); err != nil {
		shell.Logger(ctx).Error("SlackClientError", zap.String("error", err.Error()))
		return nil, err
	}

	return c, nil
}
