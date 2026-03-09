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

type TwitchAddStreamHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewTwitchAddStreamHandler(discord *discord.Discord, service *service.Service) *TwitchAddStreamHandler {
	return &TwitchAddStreamHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *TwitchAddStreamHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload", nil)
	}

	logger.Debug("🟣 Handling TwitchAddStream command from user: %s", payload.Username)

	// Verifica se foram passados argumentos (canais)
	if len(payload.Arguments) == 0 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchAddUsage)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Adiciona os canais
	addedChannels := h.Service.AddTwitchChannels(payload.Arguments...)
	if len(addedChannels) == 0 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchAddUsage)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Monta resposta
	channelsList := strings.Join(addedChannels, ", ")
	totalChannels := len(h.Service.GetTwitchChannels())
	response := fmt.Sprintf(constants.TwitchAddSuccess, channelsList, totalChannels)

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send response", err)
	}

	return nil
}
