package service

import (
	"context"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
)

// StartTwitchMonitoringForever inicia o monitoramento indefinidamente
func (s *Service) StartTwitchMonitoringForever(channelID string, intervalMinutes int) error {
	s.twitchMu.Lock()
	defer s.twitchMu.Unlock()

	if s.twitchIsMonitoring {
		return errors.NewService("monitoring already in progress", nil)
	}

	if len(s.twitchChannels) == 0 {
		return errors.NewService("no channels added for monitoring", nil)
	}

	s.twitchNotifyChannelID = channelID
	s.twitchIsMonitoring = true

	// Cria contexto cancelável (sem timeout)
	ctx, cancel := context.WithCancel(context.Background())
	s.twitchMonitoringCtx = ctx
	s.twitchCancelFunc = cancel

	logger.Info("🚀 Starting Twitch INFINITE monitoring: %d channels (interval: %d min)",
		len(s.twitchChannels), intervalMinutes)

	// Inicia goroutine de monitoramento
	go s.monitorTwitchStreams(intervalMinutes, true)

	return nil
}
