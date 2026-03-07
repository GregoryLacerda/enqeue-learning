package server

import (
	"enque-learning/events"
	"enque-learning/handlers"
	"enque-learning/integration"
	"enque-learning/internal/config"
	"enque-learning/service"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	config          *config.Config
	EventDispatcher *events.EventDispatcher
	integrations    *integration.Integrations
	service         *service.Service
	CommandHandler  *handlers.CommandHandler
	ResponseHandler *handlers.ResponseHandler
}

func NewServer(cfg *config.Config, eventDispatcher *events.EventDispatcher, integrations *integration.Integrations, service *service.Service) *Server {
	return &Server{
		config:          cfg,
		EventDispatcher: eventDispatcher,
		integrations:    integrations,
		service:         service,
		CommandHandler:  handlers.NewCommandHandler(integrations.RabbitMQ, integrations.Discord),
		ResponseHandler: handlers.NewResponseHandler(integrations.Discord, eventDispatcher, service),
	}
}

func (s *Server) StartAll() error {
	log.Println("🚀 Starting complete system (Producer + Consumer)...")

	s.EventDispatcher.RegisterHandler("discord.command.received", s.CommandHandler)

	// Start Discord
	err := s.integrations.Discord.Start()
	if err != nil {
		return fmt.Errorf("failed to start Discord: %w", err)
	}

	// Start Consumer
	msgs, err := s.integrations.RabbitMQ.Consumer()
	if err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	log.Println("✅ Complete system started successfully!")

	// Process messages in background

	for msg := range msgs {
		log.Printf("📨 Message received from queue")

		go func(msg amqp.Delivery) {
			err := s.ResponseHandler.ProcessMessage(msg.Body)
			if err != nil {
				log.Printf("❌ Error processing message: %v", err)
				msg.Nack(false, true)
			} else {
				log.Printf("✅ Message processed successfully")
				msg.Ack(false)
			}
		}(msg)
	}

	// Wait for shutdown
	s.waitForShutdown()

	return s.Shutdown()
}

func (s *Server) waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("⚠️ Interrupt signal received...")
}

func (s *Server) Shutdown() error {
	log.Println("🛑 Shutting down system...")

	if s.integrations.Discord != nil {
		s.integrations.Discord.Stop()
	}

	if s.integrations.RabbitMQ != nil {
		s.integrations.RabbitMQ.Close()
	}

	log.Println("✅ System shut down successfully!")
	return nil
}
