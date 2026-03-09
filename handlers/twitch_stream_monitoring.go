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

type TwitchStreamMonitoringHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewTwitchStreamMonitoringHandler(discord *discord.Discord, service *service.Service) *TwitchStreamMonitoringHandler {
	return &TwitchStreamMonitoringHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *TwitchStreamMonitoringHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload", nil)
	}

	logger.Debug("📺 Handling TwitchStreamMonitoring command from user: %s", payload.Username)

	// Verifica se foram passados os 2 argumentos necessários
	if len(payload.Arguments) != 2 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchStartUsage)
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

	// Parse duração
	duration, err := strconv.Atoi(payload.Arguments[0])
	if err != nil || duration <= 0 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchInvalidDuration)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Parse intervalo
	interval, err := strconv.Atoi(payload.Arguments[1])
	if err != nil || interval <= 0 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchInvalidInterval)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Inicia monitoramento
	err = h.Service.StartTwitchMonitoring(payload.ChannelID, duration, interval)
	if err != nil {
		response := fmt.Sprintf(constants.TwitchStartError, err.Error())
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Resposta de sucesso
	response := fmt.Sprintf(constants.TwitchStartSuccess, duration, interval, len(channels))
	err = h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send response", err)
	}

	return nil
}
