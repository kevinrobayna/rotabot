package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kevinrobayna/rotabot/internal/rotabot"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Resource interface {
	HealthCheck() http.HandlerFunc
	SlackEvents() http.HandlerFunc
}

type slackConfig struct {
	signingSecret string `yaml:"signing_secret"`
	clientSecret  string `yaml:"client_secret"`
}

type resource struct {
	slackConfig slackConfig
	handler     rotabot.Handler
}

func (s *resource) HealthCheck() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if _, err := slackClient(s.slackConfig, req.Context()); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "OK")
	}
}

func (s *resource) SlackEvents() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		sv, err := slack.NewSecretsVerifier(req.Header, s.slackConfig.signingSecret)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal(body, &r)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "text")
			fmt.Fprint(res, r.Challenge)
		}

		client, err := slackClient(s.slackConfig, req.Context())
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.handler.Handle(req.Context(), client, &eventsAPIEvent)
		res.WriteHeader(http.StatusOK)
	}
}

type zapWrap struct {
	l *zap.Logger
}

func (z zapWrap) Output(_ int, message string) error {
	z.l.Info(message)
	return nil
}

func slackClient(cfg slackConfig, ctx context.Context) (*slack.Client, error) {
	c := slack.New(
		cfg.clientSecret,
		slack.OptionLog(zapWrap{l: shell.Logger(ctx)}),
		slack.OptionDebug(true),
	)

	if _, err := c.AuthTestContext(ctx); err != nil {
		shell.Logger(ctx).Error("SlackClientError", zap.String("error", err.Error()))
		return nil, err
	}

	return c, nil
}
