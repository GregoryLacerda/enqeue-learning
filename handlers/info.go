package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
	"enque-learning/service"
)

type InfoCommandHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewInfoCommandHandler(discord *discord.Discord, service *service.Service) *InfoCommandHandler {
	return &InfoCommandHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *InfoCommandHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload for info command", nil)
	}

	logger.Debug("ℹ️ Handling info command from user: %s", payload.Username)

	infoData := service.InfoData{
		Username:  payload.Username,
		UserID:    payload.UserID,
		Command:   payload.Command,
		ChannelID: payload.ChannelID,
		GuildID:   payload.GuildID,
		Timestamp: payload.Timestamp,
	}

	response := h.Service.ProcessInfo(infoData)

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send info response", err)
	}

	return nil
}
