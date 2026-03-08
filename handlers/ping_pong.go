package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
	"enque-learning/service"
)

type PingCommandHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewPingCommandHandler(discord *discord.Discord, service *service.Service) *PingCommandHandler {
	return &PingCommandHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *PingCommandHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload", nil)
	}

	logger.Debug("🏓 Processing ping command from user: %s", payload.Username)

	response := h.Service.ProcessPing()

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send pong message", err)
	}

	return nil
}
