package server

import (
	"enque-learning/events"
	"enque-learning/handlers"
	"enque-learning/integration"
	"enque-learning/internal/config"
	"enque-learning/pkg/errors"
	"enque-learning/pkg/logger"
	"enque-learning/service"
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
	logger.Info("🚀 Starting complete system (Producer + Consumer)...")

	s.EventDispatcher.RegisterHandler("discord.command.received", s.CommandHandler)

	// Start Discord
	err := s.integrations.Discord.Start()
	if err != nil {
		return errors.NewIntegration("failed to start Discord", err)
	}

	// Start Consumer
	msgs, err := s.integrations.RabbitMQ.Consumer()
	if err != nil {
		return errors.NewIntegration("failed to start consumer", err)
	}

	logger.Info("✅ Complete system started successfully!")

	// Process messages in background

	for msg := range msgs {
		logger.Debug("📨 Message received from queue")

		go func(msg amqp.Delivery) {
			err := s.ResponseHandler.ProcessMessage(msg.Body)
			if err != nil {
				logger.Warn("❌ Error processing message: %v", err)
				msg.Nack(false, true)
			} else {
				logger.Debug("✅ Message processed successfully")
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
	logger.Warn("⚠️ Interrupt signal received...")
}

func (s *Server) Shutdown() error {
	logger.Info("🛑 Shutting down system...")

	if s.integrations.Discord != nil {
		s.integrations.Discord.Stop()
	}

	if s.integrations.RabbitMQ != nil {
		s.integrations.RabbitMQ.Close()
	}

	logger.Info("✅ System shut down successfully!")
	return nil
}
