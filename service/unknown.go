package service

import (
	"discordcommandbot/constants"
	"fmt"
)

func (s *Service) ProcessUnknownCommand(command string) string {
	return fmt.Sprintf(constants.UnknownCommandTemplate, command)
}
