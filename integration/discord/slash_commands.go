package discord

import (
	"discordcommandbot/events"
	"discordcommandbot/pkg/logger"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const metadataEphemeralAck = "ephemeral_ack"

var twitchChannelSuggestions = []string{
	"gaules",
	"alanzoka",
	"cellbit",
	"casimito",
	"brtt",
	"yoda",
	"fps_shaka",
	"xqc",
	"kai_cenat",
	"ibai",
	"tarik",
	"shroud",
	"pokimane",
	"sodapoppin",
	"asmongold",
}

var slashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Checks if the bot is online and responding.",
	},
	{
		Name:        "hello",
		Description: "Sends a greeting message.",
	},
	{
		Name:        "help",
		Description: "Shows all available commands.",
	},
	{
		Name:        "info",
		Description: "Shows user and system information.",
	},
	{
		Name:        "calc",
		Description: "Calculates a simple expression like '2 + 2'.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "expression",
				Description: "Expression in format number operator number (example: 2 + 2).",
				Required:    true,
			},
		},
	},
	{
		Name:        "twitch",
		Description: "Twitch monitoring management.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Adds Twitch channels for monitoring.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "channels",
						Description:  "Space-separated channel list (example: gaules brtt).",
						Required:     true,
						Autocomplete: true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "ephemeral",
						Description: "Sends command acknowledgement only to you.",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "Lists monitored Twitch channels.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "ephemeral",
						Description: "Sends command acknowledgement only to you.",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "clear",
				Description: "Clears all monitored Twitch channels.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "ephemeral",
						Description: "Sends command acknowledgement only to you.",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "start",
				Description: "Starts Twitch monitoring for a limited duration.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "duration_minutes",
						Description: "Total monitoring duration in minutes.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "check_interval_minutes",
						Description: "Interval between checks in minutes.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "ephemeral",
						Description: "Sends command acknowledgement only to you.",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "startforever",
				Description: "Starts Twitch monitoring without time limit.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "check_interval_minutes",
						Description: "Interval between checks in minutes.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "ephemeral",
						Description: "Sends command acknowledgement only to you.",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "stop",
				Description: "Stops active Twitch monitoring.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "ephemeral",
						Description: "Sends command acknowledgement only to you.",
						Required:    false,
					},
				},
			},
		},
	},
}

func (d *Discord) registerSlashCommands() {
	appID := d.Session.State.User.ID
	guildID := d.Config.GuildID
	existingCommands, err := d.Session.ApplicationCommands(appID, guildID)
	if err != nil {
		logger.Warn("❌ Failed to list existing slash commands: %v", err)
		return
	}

	if guildID == "" {
		logger.Info("🔍 Registering slash commands globally (can take time to propagate)")
	} else {
		logger.Info("🔍 Registering slash commands for guild: %s", guildID)
	}

	existingByName := make(map[string]*discordgo.ApplicationCommand, len(existingCommands))
	for _, cmd := range existingCommands {
		existingByName[cmd.Name] = cmd
	}

	for _, cmd := range slashCommands {
		existing, ok := existingByName[cmd.Name]
		if ok {
			_, err = d.Session.ApplicationCommandEdit(appID, guildID, existing.ID, cmd)
			if err != nil {
				logger.Warn("❌ Failed to update slash command %s: %v", cmd.Name, err)
				continue
			}
			logger.Info("✅ Slash command updated: /%s", cmd.Name)
			continue
		}

		_, err = d.Session.ApplicationCommandCreate(appID, guildID, cmd)
		if err != nil {
			logger.Warn("❌ Failed to register slash command %s: %v", cmd.Name, err)
			continue
		}
		logger.Info("✅ Slash command registered: /%s", cmd.Name)
	}
}

func (d *Discord) interactionAutocompleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	data := i.ApplicationCommandData()
	if data.Name != "twitch" || len(data.Options) == 0 {
		return
	}

	subcommand := data.Options[0]
	if subcommand.Name != "add" {
		return
	}

	focusedValue, ok := getFocusedOptionValue(subcommand.Options, "channels")
	if !ok {
		return
	}

	choices := buildChannelAutocompleteChoices(focusedValue)
	err := d.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		logger.Warn("❌ Failed to respond to autocomplete: %v", err)
	}
}

func (d *Discord) interactionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	payload := d.buildSlashPayload(i)
	payload.Command = resolveSlashCommandName(data)
	payload.Arguments = extractSlashArguments(data)
	payload.Metadata = extractSlashMetadata(data)
	if payload.Command == "" {
		d.respondToInteraction(i, "❌ Unknown slash command payload.", true)
		return
	}

	event := events.NewEvent("discord.command.received")
	event.Payload = payload

	err := d.Dispatcher.Dispatch(event)
	if err != nil {
		logger.Warn("❌ Error processing slash command: %v", err)
		d.respondToInteraction(i, "❌ Failed to process command.", shouldUseEphemeralAck(payload.Metadata))
		return
	}

	displayName := slashDisplayName(data)
	d.respondToInteraction(i, fmt.Sprintf("✅ Command /%s received. Check this channel for the response.", displayName), shouldUseEphemeralAck(payload.Metadata))
}

func (d *Discord) buildSlashPayload(i *discordgo.InteractionCreate) DiscordCommandPayload {
	userID := ""
	username := "unknown"

	if i.Member != nil && i.Member.User != nil {
		userID = i.Member.User.ID
		username = i.Member.User.Username
	} else if i.User != nil {
		userID = i.User.ID
		username = i.User.Username
	}

	guildID := ""
	if i.GuildID != "" {
		guildID = i.GuildID
	}

	return DiscordCommandPayload{
		UserID:    userID,
		Username:  username,
		ChannelID: i.ChannelID,
		GuildID:   guildID,
		Command:   "",
		Arguments: []string{},
		MessageID: "",
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func extractSlashArguments(data discordgo.ApplicationCommandInteractionData) []string {
	if len(data.Options) == 0 {
		return []string{}
	}

	switch data.Name {
	case "calc":
		expression := data.Options[0].StringValue()
		return strings.Fields(expression)
	case "twitch":
		return extractTwitchSubcommandArguments(data.Options[0])
	default:
		return []string{}
	}
}

func resolveSlashCommandName(data discordgo.ApplicationCommandInteractionData) string {
	if data.Name != "twitch" {
		return data.Name
	}

	if len(data.Options) == 0 {
		return ""
	}

	subcommand := data.Options[0]
	switch subcommand.Name {
	case "add":
		return "twitch.add"
	case "list":
		return "twitch.list"
	case "clear":
		return "twitch.clear"
	case "start":
		return "twitch.start"
	case "startforever":
		return "twitch.startforever"
	case "stop":
		return "twitch.stop"
	default:
		return ""
	}
}

func slashDisplayName(data discordgo.ApplicationCommandInteractionData) string {
	if data.Name != "twitch" || len(data.Options) == 0 {
		return data.Name
	}

	return fmt.Sprintf("%s %s", data.Name, data.Options[0].Name)
}

func extractTwitchSubcommandArguments(subcommand *discordgo.ApplicationCommandInteractionDataOption) []string {
	switch subcommand.Name {
	case "add":
		channels := getStringOption(subcommand.Options, "channels")
		return strings.Fields(channels)
	case "start":
		duration := fmt.Sprintf("%d", getIntOption(subcommand.Options, "duration_minutes"))
		interval := fmt.Sprintf("%d", getIntOption(subcommand.Options, "check_interval_minutes"))
		return []string{duration, interval}
	case "startforever":
		interval := fmt.Sprintf("%d", getIntOption(subcommand.Options, "check_interval_minutes"))
		return []string{interval}
	default:
		return []string{}
	}
}

func extractSlashMetadata(data discordgo.ApplicationCommandInteractionData) map[string]string {
	metadata := map[string]string{}
	if data.Name != "twitch" || len(data.Options) == 0 {
		return metadata
	}

	subcommand := data.Options[0]
	ephemeral, found := getBoolOption(subcommand.Options, "ephemeral")
	if found {
		metadata[metadataEphemeralAck] = fmt.Sprintf("%t", ephemeral)
	}

	return metadata
}

func shouldUseEphemeralAck(metadata map[string]string) bool {
	if metadata == nil {
		return true
	}

	value, ok := metadata[metadataEphemeralAck]
	if !ok {
		return true
	}

	return value == "true"
}

func getStringOption(options []*discordgo.ApplicationCommandInteractionDataOption, optionName string) string {
	for _, option := range options {
		if option.Name == optionName {
			return option.StringValue()
		}
	}

	return ""
}

func getIntOption(options []*discordgo.ApplicationCommandInteractionDataOption, optionName string) int64 {
	for _, option := range options {
		if option.Name == optionName {
			return option.IntValue()
		}
	}

	return 0
}

func getBoolOption(options []*discordgo.ApplicationCommandInteractionDataOption, optionName string) (bool, bool) {
	for _, option := range options {
		if option.Name == optionName {
			return option.BoolValue(), true
		}
	}

	return false, false
}

func getFocusedOptionValue(options []*discordgo.ApplicationCommandInteractionDataOption, optionName string) (string, bool) {
	for _, option := range options {
		if option.Name == optionName && option.Focused {
			return option.StringValue(), true
		}
	}

	return "", false
}

func buildChannelAutocompleteChoices(currentValue string) []*discordgo.ApplicationCommandOptionChoice {
	baseValue, query := splitAutocompleteInput(currentValue)
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, 25)
	query = strings.ToLower(query)

	for _, candidate := range twitchChannelSuggestions {
		if query != "" && !strings.HasPrefix(candidate, query) {
			continue
		}

		choiceValue := baseValue + candidate
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  candidate,
			Value: choiceValue,
		})

		if len(choices) >= 25 {
			break
		}
	}

	return choices
}

func splitAutocompleteInput(input string) (string, string) {
	if strings.HasSuffix(input, " ") {
		if strings.TrimSpace(input) == "" {
			return "", ""
		}

		return strings.TrimSpace(input) + " ", ""
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", ""
	}

	if len(parts) == 1 {
		return "", parts[0]
	}

	base := strings.Join(parts[:len(parts)-1], " ") + " "
	query := parts[len(parts)-1]
	return base, query
}

func (d *Discord) respondToInteraction(i *discordgo.InteractionCreate, message string, ephemeral bool) {
	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	err := d.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   flags,
		},
	})
	if err != nil {
		logger.Warn("❌ Failed to respond to slash interaction: %v", err)
	}
}
