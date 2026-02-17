package integration

import (
	"enque-learning/integration/discord"
	"enque-learning/integration/rabbitmq"
	"enque-learning/integration/twitch"
	"enque-learning/internal/config"
)

type Integrations struct {
	Discord  *discord.Discord
	Twitch   *twitch.Twitch
	RabbitMQ *rabbitmq.RabbitMQ
}

func NewIntegrations(config config.Config) *Integrations {
	return &Integrations{
		Discord:  discord.NewDiscordIntegration(config),
		Twitch:   twitch.NewTwitchIntegration(config),
		RabbitMQ: rabbitmq.NewRabbitMQIntegration(config),
	}
}
