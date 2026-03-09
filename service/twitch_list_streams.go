package service

// ListTwitchChannels returns a copy of all monitored Twitch channels.
func (s *Service) ListTwitchChannels() []string {
	return s.GetTwitchChannels()
}
