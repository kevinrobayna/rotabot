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

	blocks := []slack.Block{
		slack.NewContextBlock(
			"help",
			[]slack.MixedElement{
				slack.NewTextBlockObject(slack.MarkdownType, "This is the list of actions available for Rotabot: :robot_face:", false, false),
				slack.NewTextBlockObject(slack.MarkdownType, "- create name_of_rota schedule: Creates a new rota with the given schedule(DAILY, WEEKLY)", false, false),
				slack.NewTextBlockObject(slack.MarkdownType, "- add_member name_of_rota \"@User\": Adds a member to the given rota", false, false),
				slack.NewTextBlockObject(slack.MarkdownType, "- list: Lists all the rotas available on this channel", false, false),
				slack.NewTextBlockObject(slack.MarkdownType, "- current name_of_rota: Shows the current member on the given rota", false, false),
			}...,
		),
	}

	_, _, err := svc.client.PostMessageContext(
		ctx,
		cmd.ChannelID,

		slack.MsgOptionBlocks(blocks...),
		slack.MsgOptionAsUser(true),
		slack.MsgOptionPostEphemeral(cmd.UserID),
	)

	if err != nil {
		l.Warn("Failed to post help message", zap.Error(err))
		return err
	}

	l.Debug("Help message posted successfully")
	return nil
}
