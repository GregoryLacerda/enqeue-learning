package handlers

import (
	"encoding/json"
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/integration/rabbitmq"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
)

type CommandHandler struct {
	RabbitMQ *rabbitmq.RabbitMQ
	Discord  *discord.Discord
}

func NewCommandHandler(rabbitMQ *rabbitmq.RabbitMQ, discord *discord.Discord) *CommandHandler {
	return &CommandHandler{
		RabbitMQ: rabbitMQ,
		Discord:  discord,
	}
}

func (h *CommandHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload", nil)
	}

	logger.Debug("🔍 Processing command: %s", payload.Command)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return errors.NewHandler("failed to marshal payload", err)
	}

	err = h.RabbitMQ.Publisher(jsonData)
	if err != nil {
		return errors.NewIntegration("failed to publish message", err)
	}

	return nil
}
