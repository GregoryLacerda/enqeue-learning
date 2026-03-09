package handlers

import (
	"discordcommandbot/constants"
	"discordcommandbot/events"
	"discordcommandbot/integration/discord"
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"
	"discordcommandbot/service"
	"fmt"
	"strings"
)

type TwitchListStreamsHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewTwitchListStreamsHandler(discord *discord.Discord, service *service.Service) *TwitchListStreamsHandler {
	return &TwitchListStreamsHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *TwitchListStreamsHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload", nil)
	}

	logger.Debug("📋 Handling TwitchListStreams command from user: %s", payload.Username)

	channels := h.Service.ListTwitchChannels()
	if len(channels) == 0 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchListStreamsEmpty)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	formattedChannels := "- " + strings.Join(channels, "\n- ")
	response := fmt.Sprintf(constants.TwitchListStreamsSuccess, len(channels), formattedChannels)

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send response", err)
	}

	return nil
}
