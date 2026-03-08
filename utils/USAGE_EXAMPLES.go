package utils

/*
USAGE EXAMPLES FOR THE ERROR AND LOGGING SYSTEM

This file demonstrates how to use the custom error and logging system
in the enque-learning project.

================================================================================
SETUP (already done in cmd/main.go)
================================================================================

	import "enque-learning/pkg/logger"
	import "enque-learning/pkg/errors"

	// Initialize logger at startup
	logger.Init(config.DebugMode)

================================================================================
LOGGING EXAMPLES
================================================================================

1. DEBUG LOGS (only shown when DEBUG_MODE=true)

	logger.Debug("Checking user permissions: %s", userID)
	logger.Debug("Cache hit for key: %s", cacheKey)

2. INFO LOGS (general information)

	logger.Info("✅ Discord bot connected successfully")
	logger.Info("📨 Message received from channel: %s", channelID)

3. WARN LOGS (warnings, non-critical issues)

	logger.Warn("⚠️ API rate limit approaching: %d/%d", current, max)
	logger.Warn("⏰ Monitoring timeout in 5 minutes")

4. CRITICAL LOGS (critical errors, system failures)

	logger.Critical("❌ Failed to connect to RabbitMQ: %v", err)
	logger.Critical("🛑 System shutdown initiated")

================================================================================
ERROR CREATION EXAMPLES
================================================================================

1. VALIDATION ERRORS (user input errors)

	// Simple validation error
	err := errors.NewValidation("missing required argument", nil)

	// With formatted message
	err := errors.NewValidationf("invalid channel name: %s", channelName)

	// With context
	err := errors.NewValidation("invalid command format", nil).
		WithContext("command", command).
		WithContext("user", userID)

2. CONFIG ERRORS (configuration issues - CRITICAL)

	err := errors.NewConfig("Discord token not found", nil)
	err := errors.NewConfigf("invalid port number: %s", port)

3. INTEGRATION ERRORS (external service errors - Discord, RabbitMQ, Twitch)

	err := errors.NewIntegration("failed to send Discord message", originalErr)
	err := errors.NewIntegrationf("Twitch API returned status %d", statusCode)

4. SERVICE ERRORS (business logic errors)

	err := errors.NewService("monitoring already in progress", nil)
	err := errors.NewServicef("no channels found for user: %s", userID)

5. HANDLER ERRORS (request handling errors)

	err := errors.NewHandler("failed to parse payload", originalErr)
	err := errors.NewHandlerf("unknown command: %s", command)

6. API ERRORS (API call failures)

	err := errors.NewAPI("Twitch API request failed", originalErr).
		WithContext("endpoint", "/streams").
		WithContext("status_code", 429)

7. NETWORK ERRORS (connection issues)

	err := errors.NewNetwork("connection timeout", originalErr)
	err := errors.NewNetworkf("failed to resolve host: %s", hostname)

8. AUTH ERRORS (authentication/authorization failures - CRITICAL)

	err := errors.NewAuth("invalid Discord token", nil)
	err := errors.NewAuthf("Twitch OAuth failed: %v", originalErr)

9. DATABASE ERRORS (database operations - CRITICAL)

	err := errors.NewDatabase("failed to connect to database", originalErr)
	err := errors.NewDatabasef("query timeout: %s", query)

================================================================================
LOGGING ERRORS
================================================================================

// Automatically log at the error's level
appErr := errors.NewAPI("Twitch API request failed", err)
logger.LogError(appErr)

// Or log at specific level
logger.Warn("API error: %v", appErr)
logger.Critical("Critical failure: %v", appErr)

================================================================================
PRACTICAL EXAMPLES IN CODE
================================================================================

// Example 1: Handler validation
func (h *CommandHandler) HandleEvent(event events.EventInterface) error {
	payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
	if !ok {
		err := errors.NewHandler("invalid payload type", nil)
		logger.LogError(err)
		return err
	}

	logger.Debug("Processing command: %s from user: %s", payload.Command, payload.Username)

	if len(payload.Arguments) < 1 {
		err := errors.NewValidationf("missing required argument for command: %s", payload.Command)
		logger.LogError(err)
		return err
	}

	logger.Info("✅ Command validated: %s", payload.Command)
	return nil
}

// Example 2: Service layer error handling
func (s *Service) StartTwitchMonitoring(duration, interval time.Duration) error {
	logger.Debug("Attempting to start Twitch monitoring: duration=%v, interval=%v", duration, interval)

	s.twitchMu.Lock()
	defer s.twitchMu.Unlock()

	if s.twitchIsMonitoring {
		err := errors.NewService("monitoring already in progress", nil).
			WithContext("duration_requested", duration).
			WithContext("interval_requested", interval)
		logger.LogError(err)
		return err
	}

	if len(s.twitchChannels) == 0 {
		err := errors.NewValidation("no channels added for monitoring", nil)
		logger.LogError(err)
		return err
	}

	logger.Info("🚀 Starting Twitch monitoring for %d channels", len(s.twitchChannels))
	// ... monitoring logic

	return nil
}

// Example 3: Integration error handling
func (d *Discord) SendMessage(channelID, message string) error {
	logger.Debug("Sending Discord message to channel: %s", channelID)

	_, err := d.Session.ChannelMessageSend(channelID, message)
	if err != nil {
		appErr := errors.NewIntegration("failed to send Discord message", err).
			WithContext("channel_id", channelID).
			WithContext("message_length", len(message))
		logger.LogError(appErr)
		return appErr
	}

	logger.Debug("✅ Message sent successfully to channel: %s", channelID)
	return nil
}

// Example 4: Configuration validation
func LoadConfig() (*Config, error) {
	logger.Debug("Loading configuration from .env file")

	err := godotenv.Load()
	if err != nil {
		appErr := errors.NewConfig("failed to load .env file", err)
		logger.LogError(appErr)
		return nil, appErr
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		appErr := errors.NewConfig("DISCORD_TOKEN not found in environment", nil)
		logger.LogError(appErr)
		return nil, appErr
	}

	logger.Info("✅ Configuration loaded successfully")
	return &Config{Token: token}, nil
}

================================================================================
ERROR LEVELS GUIDE
================================================================================

DEBUG (only shown when DEBUG_MODE=true):
- Detailed debugging information
- Variable values during execution
- Step-by-step flow tracking
- Cache hits/misses
- Internal state changes

INFO:
- Successful operations
- User actions
- System events
- Connection status
- Command execution

WARN:
- Recoverable errors
- Deprecated features
- Rate limit warnings
- Timeouts (non-critical)
- API errors (with retry)
- Service degradation

CRITICAL:
- System failures
- Configuration errors
- Authentication failures
- Database connection errors
- Unrecoverable errors
- Security issues

================================================================================
MIGRATION GUIDE (from fmt.Errorf to utils errors)
================================================================================

BEFORE:
	return fmt.Errorf("monitoring already in progress")

AFTER:
	return errors.NewService("monitoring already in progress", nil)

---

BEFORE:
	return fmt.Errorf("invalid channel name: %s", channel)

AFTER:
	return errors.NewValidationf("invalid channel name: %s", channel)

---

BEFORE:
	log.Printf("Starting monitoring for user: %s", userID)

AFTER:
	logger.Info("🚀 Starting monitoring for user: %s", userID)

---

BEFORE:
	log.Printf("DEBUG: Cache value: %v", cacheVal)

AFTER:
	logger.Debug("Cache value: %v", cacheVal)

---

BEFORE:
	return fmt.Errorf("failed to connect: %w", err)

AFTER:
	return errors.NewNetwork("failed to connect", err)

================================================================================
*/
