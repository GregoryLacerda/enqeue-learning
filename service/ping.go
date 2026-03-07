package service

import "enque-learning/constants"

func (s *Service) ProcessPing() string {
	return constants.PingMessage
}
