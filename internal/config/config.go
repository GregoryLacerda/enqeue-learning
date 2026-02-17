package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL   string
	QueueName     string
	ExchangeName  string
	WebServerPort string
}

func LoadConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Fatal Error loading .env")
	}

	return &Config{
		RabbitMQURL:   os.Getenv("RABBITMQ_URL"),
		QueueName:     os.Getenv("QUEUE_NAME"),
		ExchangeName:  os.Getenv("EXCHANGE_NAME"),
		WebServerPort: os.Getenv("WEB_SERVER_PORT"),
	}
}
