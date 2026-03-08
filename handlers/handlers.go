package handlers

import (
	"encoding/json"
	"discordcommandbot/events"
	"discordcommandbot/integration/discord"
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"
	"discordcommandbot/service"
	"fmt"
	"strings"
)

type ResponseHandler struct {
	Discord    *discord.Discord
	dispatcher *events.EventDispatcher
	service    *service.Service
}

func NewResponseHandler(discord *discord.Discord, dispatcher *events.EventDispatcher, service *service.Service) *ResponseHandler {

	// Register specific handlers for each command
	//greetings handlers
	dispatcher.RegisterHandler("discord.command.hello", NewHelloCommandHandler(discord, service))
	dispatcher.RegisterHandler("discord.command.oi", NewHelloCommandHandler(discord, service))
	dispatcher.RegisterHandler("discord.command.olá", NewHelloCommandHandler(discord, service))

	//help handlers
	dispatcher.RegisterHandler("discord.command.help", NewHelpCommandHandler(discord, service))
	dispatcher.RegisterHandler("discord.command.ajuda", NewHelpCommandHandler(discord, service))
	dispatcher.RegisterHandler("discord.command.info", NewInfoCommandHandler(discord, service))

	//other handlers
	dispatcher.RegisterHandler("discord.command.ping", NewPingCommandHandler(discord, service))
	dispatcher.RegisterHandler("discord.command.calc", NewCalcCommandHandler(discord, service))

	// Twitch handlers
	dispatcher.RegisterHandler("discord.command.twitchaddstream", NewTwitchAddStreamHandler(discord, service))
	dispatcher.RegisterHandler("discord.command.twitchstreammonitoring", NewTwitchStreamMonitoringHandler(discord, service))
	dispatcher.RegisterHandler("discord.command.twitchstreammonitoringforever", NewTwitchStreamMonitoringForeverHandler(discord, service))
	dispatcher.RegisterHandler("discord.command.twitchstopmonitoring", NewTwitchStopMonitoringHandler(discord, service))

	return &ResponseHandler{
		Discord:    discord,
		dispatcher: dispatcher,
		service:    service,
	}
}

func (h *ResponseHandler) ProcessMessage(message []byte) error {
	var payload discord.DiscordCommandPayload
	err := json.Unmarshal(message, &payload)
	if err != nil {
		return errors.NewHandler("failed to unmarshal message", err)
	}

	logger.Debug("🔍 Processing response to Discord channel %s: %s", payload.ChannelID, payload.Command)

	eventName := fmt.Sprintf("discord.command.%s", strings.ToLower(payload.Command))
	event := events.NewEvent(eventName)
	event.Payload = payload

	// Always try dispatch first (goes through event system)
	err = h.dispatcher.Dispatch(event)
	if err != nil {
		return errors.NewHandler("failed to dispatch event", err)
	}

	// If there were no registered handlers, treat as unknown command
	if !h.dispatcher.HasAnyHandler(eventName) {
		unknownHandler := NewUnknownCommandHandler(h.Discord, h.service)
		err = unknownHandler.HandleEvent(event)
		if err != nil {
			return errors.NewHandler("failed to handle unknown command", err)
		}
	}

	logger.Debug("✅ Response sent to Discord with success")
	return nil
}

/*
func (h *ResponseHandler) processComand(payload discord.DiscordCommandPayload) string {
	command := strings.ToLower(payload.Command)

	switch command {
	case "ping":
		return "Pong! 🏓"
	case "hello", "oi", "olá":
		return fmt.Sprintf("👋 Olá, %s! Como posso ajudar?", payload.Username)

	case "help", "ajuda":
		return `📚 **Comandos Disponíveis:**

		!ping - Testa se o bot está respondendo
		!hello - Recebe uma saudação
		!calc <expressão> - Calcula uma expressão matemática
		!info - Mostra informações sobre o sistema
		!help - Mostra esta mensagem de ajuda`

	case "calc":
		if len(payload.Arguments) == 0 {
			return "❌ Uso: !calc <expressão>\nExemplo: !calc 2 + 2"
		}

		expression := strings.Join(payload.Arguments, " ")
		return fmt.Sprintf("🧮 Calculando: %s\n(Implementar lógica de cálculo)", expression)

	case "info":
		return fmt.Sprintf(`ℹ️ **Informações do Sistema:**

		👤 Usuário: %s
		🆔 User ID: %s
		📝 Comando: %s
		⏰ Timestamp: %s`,
			payload.Username,
			payload.UserID,
			payload.Command,
			payload.Timestamp,
		)

	default:
		return fmt.Sprintf("❓ Comando desconhecido: `%s`\nUse !help para ver os comandos disponíveis.", command)
	}
}
*/
