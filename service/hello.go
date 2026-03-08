package service

import (
	"discordcommandbot/constants"
	"fmt"
)

func (s *Service) ProcessHello(username string) string {
	return fmt.Sprintf(constants.HelloMessageTemplate, username)
}
