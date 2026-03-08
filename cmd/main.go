package main

import (
	"context"
	"discordcommandbot/events"
	"discordcommandbot/integration"
	"discordcommandbot/internal/config"
	"discordcommandbot/pkg/logger"
	"discordcommandbot/server"
	"discordcommandbot/service"
)

func main() {

	ctx := context.Background()

	cfg := config.LoadConfig()

	// Initialize logger with debug mode from config
	logger.Init(cfg.DebugMode)
	logger.Info("🚀 Starting DiscordCommandBot bot...")
	logger.Debug("Configuration loaded: DebugMode=%v, LogLevel=%s", cfg.DebugMode, cfg.LogLevel)

	dispatcher := events.NewEventDispatcher()

	integrations, err := integration.NewIntegrations(ctx, cfg, dispatcher)
	if err != nil {
		endAsError(err)
	}

	srv := service.NewService(cfg, integrations)

	server := server.NewServer(cfg, dispatcher, integrations, srv)

	if err := server.StartAll(); err != nil {
		endAsError(err)
	}
}

func endAsError(err error) {
	logger.Critical("Fatal error: %v", err)
	panic(err)
}
