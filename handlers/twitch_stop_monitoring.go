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

type TwitchStopMonitoringHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewTwitchStopMonitoringHandler(discord *discord.Discord, service *service.Service) *TwitchStopMonitoringHandler {
	return &TwitchStopMonitoringHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *TwitchStopMonitoringHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload", nil)
	}

	logger.Debug("🛑 Handling TwitchStopMonitoring command from user: %s", payload.Username)

	// Para o monitoramento
	err := h.Service.StopTwitchMonitoring()
	if err != nil {
		response := fmt.Sprintf(constants.TwitchStopMonitoringError, err.Error())
		err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
		if err != nil {
			return errors.NewIntegration("failed to send response", err)
		}
		return nil
	}

	// Resposta de sucesso
	err = h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, constants.TwitchStopMonitoringSuccess)
	if err != nil {
		return errors.NewIntegration("failed to send response", err)
	}

	return nil
}
