package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
	"enque-learning/service"
)

type HelpCommandHandler struct {
	Discord *discord.Discord
	service *service.Service
}

func NewHelpCommandHandler(discord *discord.Discord, service *service.Service) *HelpCommandHandler {
	return &HelpCommandHandler{
		Discord: discord,
		service: service,
	}
}

func (h *HelpCommandHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return errors.NewHandler("invalid payload for help command", nil)
	}

	logger.Debug("📚 Handling help command from user: %s", payload.Username)

	response := h.service.ProcessHelp()

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send help response", err)
	}

	return nil
}
