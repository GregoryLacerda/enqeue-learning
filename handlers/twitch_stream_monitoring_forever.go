package handlers

import (
	"enque-learning/constants"
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
	"enque-learning/service"
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
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchStreamMonitoringForeverUsage)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Verifica se há canais adicionados
	channels := h.Service.GetTwitchChannels()
	if len(channels) == 0 {
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchNoChannelsAdded)
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
		response := fmt.Sprintf(constants.TwitchStreamMonitoringError, err.Error())
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Resposta de sucesso
	response := fmt.Sprintf(constants.TwitchStreamMonitoringForeverStarted, interval, len(channels))
	err = h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send response", err)
	}

	return nil
}
