package config

import "github.com/urfave/cli/v2"

func NewAppConfig(c *cli.Context) *AppConfig {
	return &AppConfig{
		Debug: c.Bool("verbose"),
		Slack: SlackConfig{
			SigningSecret: c.String("slack.signing_secret"),
			ClientSecret:  c.String("slack.client_secret"),
		},
	}
}

type AppConfig struct {
	Debug bool
	Slack SlackConfig
}

type SlackConfig struct {
	// Basic Information > App Credentials > Signing Secret
	SigningSecret string
	// OAuth & Permissions > OAuth Tokens for Your Workspace > Bot User OAuth Access Token
	ClientSecret string
}
