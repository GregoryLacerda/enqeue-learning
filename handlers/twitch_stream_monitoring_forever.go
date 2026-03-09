package handlers

import (
	"discordcommandbot/constants"
	"discordcommandbot/events"
	"discordcommandbot/integration/discord"
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"
	"discordcommandbot/service"
	"fmt"
	"strconv"
)

type TwitchStreamMonitoringForeverHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewTwitchStreamMonitoringForeverHandler(discord *discord.Discord, service *service.Service) *TwitchStreamMonitoringForeverHandler {
	return &TwitchStreamMonitoringForeverHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *TwitchStreamMonitoringForeverHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload", nil)
	}

	logger.Debug("♾️ Handling TwitchStreamMonitoringForever command from user: %s", payload.Username)

	// Verifica se foi passado o argumento necessário (intervalo)
	if len(payload.Arguments) != 1 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchStartForeverUsage)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Verifica se há canais adicionados
	channels := h.Service.GetTwitchChannels()
	if len(channels) == 0 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchChannelsEmpty)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Parse intervalo
	interval, err := strconv.Atoi(payload.Arguments[0])
	if err != nil || interval <= 0 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchInvalidInterval)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Inicia monitoramento infinito
	err = h.Service.StartTwitchMonitoringForever(payload.ChannelID, interval)
	if err != nil {
		response := fmt.Sprintf(constants.TwitchStartError, err.Error())
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Resposta de sucesso
	response := fmt.Sprintf(constants.TwitchStartForeverSuccess, interval, len(channels))
	err = h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send response", err)
	}

	return nil
}
