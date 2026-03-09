package service

import (
	"discordcommandbot/integration/twitch/models"
	"discordcommandbot/pkg/logger"
	"fmt"
	"slices"
	"time"
)

// IsTwitchMonitoring retorna se o monitoramento está ativo
func (s *Service) IsTwitchMonitoring() bool {
	s.twitchMu.RLock()
	defer s.twitchMu.RUnlock()
	return s.twitchIsMonitoring
}

// containsTwitchChannel verifica se um canal já está na lista (sem lock, uso interno)
func (s *Service) containsTwitchChannel(channel string) bool {
	return slices.Contains(s.twitchChannels, channel)
}

// monitorTwitchStreams é a goroutine que executa o monitoramento
func (s *Service) monitorTwitchStreams(intervalMinutes int, forever bool) {
	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	defer ticker.Stop()

	// Primeira verificação imediata
	s.checkTwitchStreams()

	for {
		select {
		case <-s.twitchMonitoringCtx.Done():
			s.twitchMu.Lock()
			s.twitchIsMonitoring = false
			s.twitchMu.Unlock()

			var reason string
			if forever {
				reason = "manually cancelled"
			} else {
				reason = "time limit reached"
			}

			logger.Info("⏰ Twitch monitoring finished: %s", reason)

			// Notify on Discord
			if s.twitchNotifyChannelID != "" {
				message := fmt.Sprintf("⏰ **Twitch monitoring finished**: %s", reason)
				s.integrations.Discord.SendMessage(s.twitchNotifyChannelID, message)
			}
			return

		case <-ticker.C:
			s.checkTwitchStreams()
		}
	}
}

// checkTwitchStreams verifica o status das streams e envia notificações
func (s *Service) checkTwitchStreams() {
	s.twitchMu.RLock()
	channels := make([]string, len(s.twitchChannels))
	copy(channels, s.twitchChannels)
	notifyChannelID := s.twitchNotifyChannelID
	s.twitchMu.RUnlock()

	if len(channels) == 0 {
		return
	}

	logger.Debug("🔍 Checking streams: %v", channels)

	streamsResponse, err := s.integrations.Twitch.GetStreams(channels)
	if err != nil {
		logger.Warn("❌ Error fetching streams: %v", err)
		return
	}

	if len(streamsResponse.Data) == 0 {
		logger.Debug("📭 No streams online at the moment")
		return
	}

	// Notifica sobre streams online
	for _, stream := range streamsResponse.Data {
		// Verifica se já notificou recentemente (evita spam)
		s.twitchMu.RLock()
		lastTime, exists := s.twitchLastNotified[stream.UserLogin]
		s.twitchMu.RUnlock()

		// Se já notificou há menos de 30 minutos, pula
		if exists && time.Since(lastTime) < 30*time.Minute {
			continue
		}

		message := s.formatTwitchStreamNotification(stream)

		if notifyChannelID != "" {
			err := s.integrations.Discord.SendMessage(notifyChannelID, message)
			if err != nil {
				logger.Warn("❌ Error sending notification: %v", err)
			} else {
				logger.Debug("✅ Notification sent for channel: %s", stream.UserLogin)

				// Atualiza timestamp da última notificação
				s.twitchMu.Lock()
				s.twitchLastNotified[stream.UserLogin] = time.Now()
				s.twitchMu.Unlock()
			}
		}
	}
}

// formatTwitchStreamNotification formata a mensagem de notificação
func (s *Service) formatTwitchStreamNotification(stream models.StreamData) string {
	duration := time.Since(stream.StartedAt)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf(`🎮 **%s is LIVE!**

📺 **Title:** %s
🎯 **Playing:** %s
👥 **Viewers:** %d
⏱️ **Live for:** %dh %dmin
🔗 **Watch:** <https://twitch.tv/%s>`,
		stream.UserName,
		stream.Title,
		stream.GameName,
		stream.ViewerCount,
		hours, minutes,
		stream.UserLogin,
	)
}
