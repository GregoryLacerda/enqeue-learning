package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/service"
	"fmt"
	"log"
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
		return fmt.Errorf("invalid payload")
	}

	log.Printf("processing ping command from user: %s", payload.Username)

	response := h.Service.ProcessPing()

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return fmt.Errorf("failed to send pong message: %w", err)
	}

	return nil
}
