package internal

import (
	"context"
	"github.com/kevinrobayna/rotabot/gen/rotabot"
	"github.com/kevinrobayna/rotabot/internal/shell"
)

func NewRotabotService() rotabot.Service {
	return &svc{}
}

type svc struct {
}

func (s svc) Healthcheck(ctx context.Context) (err error) {
	shell.Logger(ctx).Info("Healthcheck")
	return nil
}

var _ rotabot.Service = &svc{}
