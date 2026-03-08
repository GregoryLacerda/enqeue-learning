package service

import "discordcommandbot/constants"

func (s *Service) ProcessHelp() string {
	return constants.CommandHelpMessage
}
