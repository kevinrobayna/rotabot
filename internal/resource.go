package internal

import (
	"encoding/json"
	"github.com/kevinrobayna/rotabot/internal/config"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type resource struct {
	cfg *config.AppConfig
	endpoints
}

type endpoints interface {
	HealthCheck() http.HandlerFunc
	HandleSlashCommand() http.HandlerFunc
}

func (resource *resource) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func (resource *resource) HandleSlashCommand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shell.Logger(r.Context())

		verifier, err := slack.NewSecretsVerifier(r.Header, resource.cfg.Slack.SigningSecret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = verifier.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch s.Command {
		case "/rotabot":
			// TODO: handle command, for now we are just acknowledging the request
			w.WriteHeader(http.StatusOK)
		default:
			l.Error("unknown_command", zap.String("command", s.Command))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
