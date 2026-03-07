package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/service"
	"fmt"
	"log"
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
		return fmt.Errorf("invalid payload for help command")
	}

	log.Printf("handling help command from user: %s", payload.Username)

	response := h.service.ProcessHelp()

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return fmt.Errorf("failed to send help response: %w", err)
	}

	return nil
}
