package discord

import (
	"discordcommandbot/events"
	"discordcommandbot/internal/config"
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Config     *config.DiscordConfig
	Session    *discordgo.Session
	Dispatcher events.EventDispatcherInterface
}

func NewDiscordIntegration(config *config.DiscordConfig, dispatcher events.EventDispatcherInterface) (*Discord, error) {

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, errors.NewIntegration("failed to create Discord session", err)
	}

	discord := &Discord{
		Config:     config,
		Dispatcher: dispatcher,
		Session:    session,
	}

	session.AddHandler(discord.interactionCreateHandler)
	session.AddHandler(discord.interactionAutocompleteHandler)

	return discord, nil
}

func (d *Discord) Start() error {
	err := d.Session.Open()
	if err != nil {
		return errors.NewIntegration("failed to open Discord session", err)
	}

	logger.Info("✅ Discord bot connected and online!")

	// Register slash commands after bot identity is available.
	d.registerSlashCommands()

	return nil
}

func (d *Discord) Stop() error {
	return d.Session.Close()
}

func (d *Discord) SendMessage(channelID, message string) error {
	_, err := d.Session.ChannelMessageSend(channelID, message)
	if err != nil {
		return errors.NewIntegration("failed to send message to Discord", err)
	}
	return nil
}

func (d *Discord) ReplyToMessage(channelID, messageID, message string) error {
	if messageID == "" {
		return d.SendMessage(channelID, message)
	}

	_, err := d.Session.ChannelMessageSendReply(channelID, message, &discordgo.MessageReference{
		MessageID: messageID,
		ChannelID: channelID,
	})
	if err != nil {
		return errors.NewIntegration("failed to reply to message", err)
	}
	return nil
}
