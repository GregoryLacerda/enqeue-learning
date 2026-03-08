package service

import "discordcommandbot/constants"

func (s *Service) ProcessPing() string {
	return constants.PingMessage
}
