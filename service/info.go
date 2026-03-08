package service

import (
	"discordcommandbot/constants"
	"fmt"
)

type InfoData struct {
	Username  string
	UserID    string
	Command   string
	ChannelID string
	GuildID   string
	Timestamp string
}

func (s *Service) ProcessInfo(info InfoData) string {
	return fmt.Sprintf(
		constants.InfoMessageTemplate,
		info.Username,
		info.UserID,
		info.Command,
		info.ChannelID,
		info.GuildID,
		info.Timestamp,
	)
}
