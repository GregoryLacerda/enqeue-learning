package handlers

import (
	"enque-learning/events"
	"enque-learning/integration/discord"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
	"enque-learning/service"
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
		return errors.NewHandler("invalid payload for calc command", nil)
	}

	logger.Debug("🧮 Handling calc command from user: %s", payload.Username)

	response := h.Service.ProcessCalc(payload.Arguments)

	err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
	if err != nil {
		return errors.NewIntegration("failed to send calc response", err)
	}

	return nil
}
