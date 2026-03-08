package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
	"enque-learning/service"
)

type UnknownCommandHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewUnknownCommandHandler(discord *discord.Discord, service *service.Service) *UnknownCommandHandler {
	return &UnknownCommandHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *UnknownCommandHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload for unknown command", nil)
	}

	logger.Debug("❓ Handling unknown command: %s from user: %s", payload.Command, payload.Username)

	response := h.Service.ProcessUnknownCommand(payload.Command)

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send unknown command response", err)
	}

	return nil
}
