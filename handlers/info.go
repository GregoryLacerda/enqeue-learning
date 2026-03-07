package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/service"
	"fmt"
	"log"
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
		return fmt.Errorf("invalid payload for info command")
	}

	log.Printf("handling info command from user: %s", payload.Username)

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
		return fmt.Errorf("failed to send info response: %w", err)
	}

	return nil
}
