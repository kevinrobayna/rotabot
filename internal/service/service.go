package service

import (
	"context"
	"github.com/kevinrobayna/rotabot/gen/rotabot"
	"github.com/kevinrobayna/rotabot/internal/zapctx"
)

func NewRotabotService() rotabot.Service {
	return &svc{}
}

type svc struct {
}

func (w svc) Healthcheck(c context.Context) (res string, err error) {
	zapctx.Logger(c).Info("Healthcheck")
	return "OK", nil
}

func (w svc) Home(c context.Context) (res string, err error) {
	zapctx.Logger(c).Info("Home")
	return "Home", nil
}

var _ rotabot.Service = &svc{}
