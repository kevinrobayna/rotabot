package internal

import (
	"context"

	"github.com/kevinrobayna/rotabot/internal/config"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type commandSvc struct {
	cfg    *config.AppConfig
	client *slack.Client
	commands
}

type commands interface {
	Handle(ctx context.Context, cmd *slack.SlashCommand) error
}

func (svc *commandSvc) Handle(ctx context.Context, cmd *slack.SlashCommand) error {
	l := shell.Logger(ctx).
		With(zap.String("cmd", cmd.Command)).
		With(zap.String("txt", cmd.Text)).
		With(zap.String("user_id", cmd.UserID)).
		With(zap.String("channel_id", cmd.ChannelID)).
		With(zap.String("team_id", cmd.TeamID)).
		With(zap.String("trigger_id", cmd.TriggerID))
	switch cmd.Text {
	default:
		return svc.handleUnknown(shell.WithLogger(ctx, l), cmd)
	}
}

func (svc *commandSvc) handleUnknown(ctx context.Context, cmd *slack.SlashCommand) error {
	l := shell.Logger(ctx)
	l.Debug("Posting help message")

	_, err := svc.client.OpenViewContext(ctx, cmd.TriggerID, generateModal())
	if err != nil {
		l.Warn("Failed to post help message", zap.Error(err))
		return err
	}

	l.Debug("Help message posted successfully")
	return nil
}

func generateModal() slack.ModalViewRequest {
	return slack.ModalViewRequest{
		Type:  slack.VTModal,
		Title: slack.NewTextBlockObject(slack.MarkdownType, "Rotabot :robot_face:", false, false),
		Close: slack.NewTextBlockObject(slack.MarkdownType, "Close", false, false),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewDividerBlock(),
			},
		},
	}
}
