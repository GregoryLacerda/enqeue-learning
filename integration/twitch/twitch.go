package twitch

import "enque-learning/internal/config"

type Twitch struct {
	Config config.Config
}

func NewTwitchIntegration(config config.Config) *Twitch {
	return &Twitch{
		Config: config,
	}
}
