package service

import (
	"context"
	"discordcommandbot/integration"
	"discordcommandbot/internal/config"
	"sync"
	"time"
)

type Service struct {
	config       *config.Config
	integrations *integration.Integrations

	// Twitch monitoring fields
	twitchChannels        []string
	twitchMu              sync.RWMutex
	twitchMonitoringCtx   context.Context
	twitchCancelFunc      context.CancelFunc
	twitchIsMonitoring    bool
	twitchNotifyChannelID string
	twitchCheckInterval   time.Duration
	twitchLastNotified    map[string]time.Time
}

func NewService(cfg *config.Config, integrations *integration.Integrations) *Service {
	return &Service{
		config:             cfg,
		integrations:       integrations,
		twitchChannels:     []string{},
		twitchLastNotified: make(map[string]time.Time),
	}
}
