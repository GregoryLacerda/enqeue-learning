package handlers

import (
	"discordcommandbot/constants"
	"discordcommandbot/events"
	"discordcommandbot/integration/discord"
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"
	"discordcommandbot/service"
	"fmt"
)

type TwitchClearStreamsHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewTwitchClearStreamsHandler(discord *discord.Discord, service *service.Service) *TwitchClearStreamsHandler {
	return &TwitchClearStreamsHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *TwitchClearStreamsHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload", nil)
	}

	logger.Debug("🧹 Handling TwitchClearStreams command from user: %s", payload.Username)

	removedCount := h.Service.ClearTwitchChannels()
	response := fmt.Sprintf(constants.TwitchClearSuccess, removedCount)

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send response", err)
	}

	return nil
}
