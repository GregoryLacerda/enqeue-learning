package service

import (
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
)

// StopTwitchMonitoring para o monitoramento em andamento
func (s *Service) StopTwitchMonitoring() error {
	s.twitchMu.Lock()
	defer s.twitchMu.Unlock()

	if !s.twitchIsMonitoring {
		return errors.NewService("no monitoring in progress", nil)
	}

	if s.twitchCancelFunc != nil {
		s.twitchCancelFunc()
	}

	s.twitchIsMonitoring = false
	logger.Info("🛑 Twitch monitoring stopped")

	return nil
}
