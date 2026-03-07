package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/service"
	"fmt"
	"log"
)

type UnknownCommandHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewUnknownCommandHandler(discord *discord.Discord, service *service.Service) *UnknownCommandHandler {
	return &UnknownCommandHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *UnknownCommandHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return fmt.Errorf("invalid payload for unknown command")
	}

	log.Printf("handling unknown command: %s from user: %s", payload.Command, payload.Username)

	response := h.Service.ProcessUnknownCommand(payload.Command)

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return fmt.Errorf("failed to send unknown command response: %w", err)
	}

	return nil
}
