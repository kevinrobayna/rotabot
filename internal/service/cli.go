package service

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"net"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "rotabot",
		Usage:  "Starts the rotabot server and its dependencies",
		Action: serviceFunc(),
	}
}

func serviceFunc() cli.ActionFunc {
	return func(c *cli.Context) error {
		ctx := c.Context

		addr := ":8080"
		l, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen on %s: %w", addr, err)
		}
		defer l.Close()

		srv := NewApp(&AppDeps{
			Component: "rotabot.server",
			Service:   NewRotabotService(),
			Context:   ctx,
			Listener:  l,
		})
		return srv.Run()
	}
}
