package discord

import "enque-learning/internal/config"

type Discord struct {
	Config config.Config
}

func NewDiscordIntegration(config config.Config) *Discord {
	return &Discord{
		Config: config,
	}
}
