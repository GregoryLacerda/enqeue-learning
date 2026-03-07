package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/service"
	"fmt"
	"log"
)

type CalcCommandHandler struct {
	Discord *discord.Discord
	Service *service.Service
}

func NewCalcCommandHandler(discord *discord.Discord, service *service.Service) *CalcCommandHandler {
	return &CalcCommandHandler{
		Discord: discord,
		Service: service,
	}
}

func (h *CalcCommandHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		return fmt.Errorf("invalid payload for calc command")
	}

	log.Printf("handling calc command from user: %s", payload.Username)

	response := h.Service.ProcessCalc(payload.Arguments)

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return fmt.Errorf("failed to send calc response: %w", err)
	}

	return nil
}
