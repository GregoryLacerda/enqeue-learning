package service

import (
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"
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
	s.twitchCheckInterval = 0
	logger.Info("🛑 Twitch monitoring stopped")

	return nil
}
