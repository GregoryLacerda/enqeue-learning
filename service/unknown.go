package service

import (
	"enque-learning/constants"
	"fmt"
)

func (s *Service) ProcessUnknownCommand(command string) string {
	return fmt.Sprintf(constants.UnknownCommandTemplate, command)
}
