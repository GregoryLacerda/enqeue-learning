package config

import (
	"log"
	"os"
	"strconv"

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
	Token         string
	CommandPrefix string
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

	return &Config{
		DiscordConfig: DiscordConfig{
			Token:         os.Getenv("DISCORD_TOKEN"),
			CommandPrefix: os.Getenv("DISCORD_COMMAND_PREFIX"),
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
		},
		WebServerPort: os.Getenv("WEB_SERVER_PORT"),
		DebugMode:     debugMode,
		LogLevel:      logLevel,
	}
}
