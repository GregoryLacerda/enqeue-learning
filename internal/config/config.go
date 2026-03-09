package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordConfig  DiscordConfig
	RabbitMQConfig RabbitMQConfig
	TwitchConfig   TwitchConfig

	// Server
	WebServerPort string
	DebugMode     bool
	LogLevel      string
}

type DiscordConfig struct {
	Token   string
	GuildID string
}

type RabbitMQConfig struct {
	URL          string
	QueueName    string
	ExchangeName string
	RoutingKey   string
}

type TwitchConfig struct {
	ClientID     string
	ClientSecret string
	NotifyMode   string
}

func LoadConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Fatal Error loading .env")
	}

	// Parse debug mode (defaults to false)
	debugMode := false
	if debugEnv := os.Getenv("DEBUG_MODE"); debugEnv != "" {
		debugMode, _ = strconv.ParseBool(debugEnv)
	}

	// Get log level (defaults to "info")
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	// Twitch notification mode (defaults to always)
	notifyMode := strings.ToLower(strings.TrimSpace(os.Getenv("TWITCH_NOTIFY_MODE")))
	if notifyMode != "cooldown" {
		notifyMode = "always"
	}

	return &Config{
		DiscordConfig: DiscordConfig{
			Token:   os.Getenv("DISCORD_TOKEN"),
			GuildID: os.Getenv("DISCORD_GUILD_ID"),
		},
		RabbitMQConfig: RabbitMQConfig{
			URL:          os.Getenv("RABBITMQ_URL"),
			QueueName:    os.Getenv("QUEUE_NAME"),
			ExchangeName: os.Getenv("EXCHANGE_NAME"),
			RoutingKey:   os.Getenv("ROUTING_KEY"),
		},
		TwitchConfig: TwitchConfig{
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
			NotifyMode:   notifyMode,
		},
		WebServerPort: os.Getenv("WEB_SERVER_PORT"),
		DebugMode:     debugMode,
		LogLevel:      logLevel,
	}
}
