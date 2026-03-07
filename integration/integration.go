package integration

import (
	"context"
	"enque-learning/events"
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

func NewIntegrations(ctx context.Context, config *config.Config, dispatcher events.EventDispatcherInterface) (*Integrations, error) {
	discordIntegration, err := discord.NewDiscordIntegration(&config.DiscordConfig, dispatcher)
	if err != nil {
		return nil, err
	}

	twitchIntegration, err := twitch.NewTwitchIntegration(ctx, &config.TwitchConfig)
	if err != nil {
		return nil, err
	}

	rabbitMQIntegration, err := rabbitmq.NewRabbitMQIntegration(ctx, &config.RabbitMQConfig)
	if err != nil {
		return nil, err
	}

	return &Integrations{
		Discord:  discordIntegration,
		Twitch:   twitchIntegration,
		RabbitMQ: rabbitMQIntegration,
	}, nil
}
