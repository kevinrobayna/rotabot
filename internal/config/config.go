package config

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
