package service

import (
	"discordcommandbot/pkg/logger"
	"time"
)

// ClearTwitchChannels removes all channels from monitoring and resets notification cache.
func (s *Service) ClearTwitchChannels() int {
	s.twitchMu.Lock()
	defer s.twitchMu.Unlock()

	removedCount := len(s.twitchChannels)
	s.twitchChannels = []string{}
	s.twitchLastNotified = make(map[string]time.Time)

	logger.Info("🧹 Cleared Twitch channel list: %d removed", removedCount)

	return removedCount
}
