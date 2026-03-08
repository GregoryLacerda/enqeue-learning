package handlers

import (
	"discordcommandbot/events"
	"discordcommandbot/integration/discord"
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"
	"discordcommandbot/service"
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
