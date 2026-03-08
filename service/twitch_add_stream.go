package service

import "enque-learning/pkg/logger"

// AddTwitchChannels adiciona um ou mais canais à lista de monitoramento
func (s *Service) AddTwitchChannels(channels ...string) {
	s.twitchMu.Lock()
	defer s.twitchMu.Unlock()

	for _, channel := range channels {
		// Avoid duplicates
		if !s.containsTwitchChannel(channel) {
			s.twitchChannels = append(s.twitchChannels, channel)
			logger.Debug("✅ Twitch channel added: %s", channel)
		}
	}
}

// RemoveTwitchChannel remove um canal da lista
func (s *Service) RemoveTwitchChannel(channel string) {
	s.twitchMu.Lock()
	defer s.twitchMu.Unlock()

	for i, ch := range s.twitchChannels {
		if ch == channel {
			s.twitchChannels = append(s.twitchChannels[:i], s.twitchChannels[i+1:]...)
			logger.Debug("🗑️ Twitch channel removed: %s", channel)
			break
		}
	}
}

// GetTwitchChannels retorna a lista de canais
func (s *Service) GetTwitchChannels() []string {
	s.twitchMu.RLock()
	defer s.twitchMu.RUnlock()

	channelsCopy := make([]string, len(s.twitchChannels))
	copy(channelsCopy, s.twitchChannels)
	return channelsCopy
}
