package service

import (
	"context"
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"
	"time"
)

// StartTwitchMonitoring inicia o monitoramento por um tempo limitado
func (s *Service) StartTwitchMonitoring(channelID string, durationMinutes int, intervalMinutes int) error {
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
	s.twitchCheckInterval = time.Duration(intervalMinutes) * time.Minute

	// Cria contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(durationMinutes)*time.Minute)
	s.twitchMonitoringCtx = ctx
	s.twitchCancelFunc = cancel

	logger.Info("🚀 Starting Twitch monitoring: %d channels for %d minutes (interval: %d min)",
		len(s.twitchChannels), durationMinutes, intervalMinutes)
	logger.Info("🔧 Twitch notify settings: mode=%s, effective_cooldown=%d min",
		s.config.TwitchConfig.NotifyMode, intervalMinutes)

	// Inicia goroutine de monitoramento
	go s.monitorTwitchStreams(intervalMinutes, false)

	return nil
}
