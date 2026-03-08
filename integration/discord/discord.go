package discord

import (
	"enque-learning/events"
	"enque-learning/internal/config"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
	"strings"

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

	session.AddHandler(discord.messageHandler)

	return discord, nil
}

func (d *Discord) Start() error {
	err := d.Session.Open()
	if err != nil {
		return errors.NewIntegration("failed to open Discord session", err)
	}

	logger.Info("✅ Discord bot connected and online!")

	return nil
}

func (d *Discord) Stop() error {
	return d.Session.Close()
}

func (d *Discord) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, d.Config.CommandPrefix) {
		return
	}

	content := strings.TrimPrefix(m.Content, d.Config.CommandPrefix)
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}

	payload := DiscordCommandPayload{
		UserID:    m.Author.ID,
		Username:  m.Author.Username,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
		Command:   command,
		Arguments: args,
		MessageID: m.ID,
		Timestamp: m.Timestamp.String(),
	}

	event := events.NewEvent("discord.command.received")
	event.Payload = payload

	err := d.Dispatcher.Dispatch(event)
	if err != nil {
		logger.Warn("❌ Error processing command: %v", err)
		d.SendMessage(m.ChannelID, "❌ Erro ao processar comando!")
	}
}

func (d *Discord) SendMessage(channelID, message string) error {
	_, err := d.Session.ChannelMessageSend(channelID, message)
	if err != nil {
		return errors.NewIntegration("failed to send message to Discord", err)
	}
	return nil
}

func (d *Discord) ReplyToMessage(channelID, messageID, message string) error {
	_, err := d.Session.ChannelMessageSendReply(channelID, message, &discordgo.MessageReference{
		MessageID: messageID,
		ChannelID: channelID,
	})
	if err != nil {
		return errors.NewIntegration("failed to reply to message", err)
	}
	return nil
}
