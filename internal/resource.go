package internal

import (
	"encoding/json"
	"net/http"

	"github.com/kevinrobayna/rotabot/internal/config"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type resource struct {
	cfg      *config.AppConfig
	commands commandSvc
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

		s, err := slack.SlashCommandParse(r)
		if err != nil {
			l.Error("failed to parse slash command", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		switch s.Command {
		case "/rotabot":
			if err = resource.commands.Handle(r.Context(), &s); err != nil {
				l.Error("something went wrong while processing command", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			l.Error("unknown_command", zap.String("command", s.Command))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
