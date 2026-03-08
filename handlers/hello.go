package handlers

import (
	"discordcommandbot/events"
	"discordcommandbot/integration/discord"
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"
	"discordcommandbot/service"
)

type HelloCommandHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewHelloCommandHandler(discord *discord.Discord, service *service.Service) *HelloCommandHandler {
	return &HelloCommandHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *HelloCommandHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload for hello command", nil)
	}

	logger.Debug("👋 Handling hello command from user: %s", payload.Username)

	response := h.Service.ProcessHello(payload.Username)

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send hello response", err)
	}

	return nil
}
