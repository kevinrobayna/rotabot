package rotabot

import (
	"context"
	"github.com/kevinrobayna/rotabot/internal/shell"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

type Handler struct {
}

func (h Handler) Handle(ctx context.Context, client *slack.Client, event *slackevents.EventsAPIEvent) {
	if event.Type == slackevents.CallbackEvent {
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			h.handleAppMentionEvent(ctx, client, ev)
		}
	}
}

func (h Handler) handleAppMentionEvent(ctx context.Context, client *slack.Client, event *slackevents.AppMentionEvent) {
	_, _, err := client.PostMessageContext(
		ctx,
		event.Channel,
		slack.MsgOptionText("Hello world!", false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		shell.Logger(ctx).Error("failed to post message", zap.Error(err))
		return
	}
}
