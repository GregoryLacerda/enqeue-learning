package service

import (
	"discordcommandbot/pkg/logger"
	"strings"
)

// AddTwitchChannels adiciona um ou mais canais à lista de monitoramento
func (s *Service) AddTwitchChannels(channels ...string) []string {
	s.twitchMu.Lock()
	defer s.twitchMu.Unlock()

	addedChannels := make([]string, 0, len(channels))

	for _, channel := range channels {
		normalizedChannel := normalizeTwitchChannel(channel)
		if normalizedChannel == "" {
			continue
		}

		// Avoid duplicates
		if !s.containsTwitchChannel(normalizedChannel) {
			s.twitchChannels = append(s.twitchChannels, normalizedChannel)
			addedChannels = append(addedChannels, normalizedChannel)
			logger.Debug("✅ Twitch channel added: %s", normalizedChannel)
		}
	}

	return addedChannels
}

func normalizeTwitchChannel(channel string) string {
	normalized := strings.TrimSpace(strings.ToLower(channel))
	normalized = strings.TrimPrefix(normalized, "https://")
	normalized = strings.TrimPrefix(normalized, "http://")
	normalized = strings.TrimPrefix(normalized, "www.")
	normalized = strings.TrimPrefix(normalized, "twitch.tv/")
	normalized = strings.TrimPrefix(normalized, "@")
	normalized = strings.Trim(normalized, "/")

	return normalized
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
