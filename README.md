# 🤖 enque-learning

Event-driven Discord bot built with Go and RabbitMQ for asynchronous command processing.

![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)
![Discord](https://img.shields.io/badge/Discord-Bot-5865F2?logo=discord)
![RabbitMQ](https://img.shields.io/badge/RabbitMQ-FF6600?logo=rabbitmq)
![License](https://img.shields.io/badge/license-MIT-green)

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
  - [Discord Bot Setup](#discord-bot-setup)
  - [Twitch Application Setup](#twitch-application-setup)
  - [RabbitMQ Setup](#rabbitmq-setup)
  - [Environment Variables](#environment-variables)
- [Running the Bot](#running-the-bot)
- [Available Commands](#available-commands)
  - [General Commands](#general-commands)
  - [Twitch Commands](#twitch-commands)
- [Architecture](#architecture)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

---

## 🎯 Overview

**enque-learning** is a Discord bot that demonstrates event-driven architecture using Go. Commands are received from Discord, enqueued in RabbitMQ for asynchronous processing, and handled by dedicated event handlers.

### Key Highlights

- **Event-Driven Architecture**: Uses custom EventDispatcher pattern
- **Asynchronous Processing**: RabbitMQ for reliable message queuing
- **Modular Design**: One file per command for easy maintenance
- **Twitch Integration**: Monitor Twitch streams and get Discord notifications
- **Type-Safe**: Leverages Go's strong typing and interfaces

---

## ✨ Features

- ✅ **General Commands**: Ping, Help, Info, Calculator
- ✅ **Twitch Stream Monitoring**: Get notified when streamers go live
- ✅ **Event System**: Extensible command architecture
- ✅ **Graceful Shutdown**: Proper cleanup of resources
- ✅ **Anti-Spam System**: Intelligent notification cooldowns
- ✅ **Concurrent Processing**: Safe goroutine handling with context

---

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.26 or higher** - [Download Go](https://golang.org/dl/)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **Git** - [Install Git](https://git-scm.com/downloads)
- **Discord Account** - To create a bot application
- **Twitch Account** - To create a developer application (optional, for Twitch features)

---

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/enque-learning.git
cd enque-learning
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Start RabbitMQ

```bash
docker-compose up -d
```

Verify RabbitMQ is running:
- **AMQP**: `amqp://localhost:5672`
- **Management UI**: http://localhost:15672 (guest/guest)

---

## ⚙️ Configuration

### Discord Bot Setup

#### Step 1: Create Discord Application

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Click **"New Application"**
3. Give it a name (e.g., "enque-learning")
4. Click **"Create"**

#### Step 2: Create Bot User

1. In your application, go to **"Bot"** section
2. Click **"Add Bot"**
3. Confirm by clicking **"Yes, do it!"**

#### Step 3: Get Bot Token

1. In the **"Bot"** section, click **"Reset Token"**
2. Copy the token (⚠️ **NEVER share this token publicly!**)
3. Save it for the `.env` file

#### Step 4: Configure Bot Permissions

1. In the **"Bot"** section, enable these **Privileged Gateway Intents**:
   - ✅ **Presence Intent**
   - ✅ **Server Members Intent**
   - ✅ **Message Content Intent**

#### Step 5: Invite Bot to Server

1. Go to **"OAuth2"** → **"URL Generator"**
2. Select scopes:
   - ✅ `bot`
   - ✅ `applications.commands`
3. Select bot permissions:
   - ✅ `Send Messages`
   - ✅ `Read Messages/View Channels`
   - ✅ `Read Message History`
   - ✅ `Embed Links`
4. Copy the generated URL and paste it in your browser
5. Select your server and authorize the bot

---

### Twitch Application Setup

#### Step 1: Create Twitch Application

1. Go to [Twitch Developer Console](https://dev.twitch.tv/console/apps)
2. Log in with your Twitch account
3. Click **"Register Your Application"**

#### Step 2: Configure Application

Fill in the required information:

- **Name**: `enque-learning-bot` (or your preferred name)
- **OAuth Redirect URLs**: `http://localhost` (for local development)
- **Category**: Select **"Application Integration"** or **"Bot"**

Click **"Create"**

#### Step 3: Get Credentials

1. After creation, click **"Manage"** on your application
2. Copy the **Client ID**
3. Click **"New Secret"** to generate a **Client Secret**
4. Copy the secret immediately (⚠️ **You can only see it once!**)
5. Save both for the `.env` file

---

### RabbitMQ Setup

RabbitMQ runs via Docker Compose (already configured in `docker-compose.yml`):

```yaml
services:
  rabbitmq:
    image: rabbitmq:3.12-management
    container_name: rabbitmq
    ports:
      - "5672:5672"   # AMQP protocol
      - "15672:15672" # Management UI
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
```

**No additional setup required!** The bot automatically creates the exchange, queue, and bindings.

---

### Environment Variables

#### Step 1: Create `.env` File

Copy the example file:

```bash
cp .env.example .env
```

Or create manually:

```bash
touch .env
```

#### Step 2: Configure Variables

Edit `.env` with your credentials:

```env
# Discord Configuration
DISCORD_TOKEN=your_discord_bot_token_here
DISCORD_COMMAND_PREFIX=!

# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
QUEUE_NAME=discord-commands
EXCHANGE_NAME=discord-exchange
ROUTING_KEY=discord.command

# Twitch Configuration (Optional - only if using Twitch features)
TWITCH_CLIENT_ID=your_twitch_client_id
TWITCH_CLIENT_SECRET=your_twitch_client_secret

# Server Configuration (Optional)
WEB_SERVER_PORT=8080
```

#### Configuration Details

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `DISCORD_TOKEN` | ✅ Yes | Discord bot token from Developer Portal | `MTIzNDU2Nzg5MDEyMzQ1Njc4OQ.ABCDEF.xyz123...` |
| `DISCORD_COMMAND_PREFIX` | No | Command prefix (default: `!`) | `!` |
| `RABBITMQ_URL` | ✅ Yes | RabbitMQ connection string | `amqp://guest:guest@localhost:5672/` |
| `QUEUE_NAME` | No | RabbitMQ queue name | `discord-commands` |
| `EXCHANGE_NAME` | No | RabbitMQ exchange name | `discord-exchange` |
| `ROUTING_KEY` | No | RabbitMQ routing key | `discord.command` |
| `TWITCH_CLIENT_ID` | For Twitch | Twitch application client ID | `abc123xyz456...` |
| `TWITCH_CLIENT_SECRET` | For Twitch | Twitch application client secret | `def789uvw012...` |
| `WEB_SERVER_PORT` | No | Web server port (future use) | `8080` |
| `DEBUG_MODE` | No | Enable debug logging (default: `false`) | `true` or `false` |
| `LOG_LEVEL` | No | Minimum log level (default: `info`) | `debug`, `info`, `warn`, `critical` |

---

## 🏃 Running the Bot

### Local Development

#### Start RabbitMQ (if not already running)

```bash
docker-compose up -d
```

#### Run the Bot

```bash
go run cmd/main.go
```

Or build and run:

```bash
go build -o bot.exe cmd/main.go
./bot.exe
```

#### Verify Bot is Running

You should see logs like:

```
🚀 Starting complete system (Producer + Consumer)...
✅ Discord bot connected and online!
✅ Complete system started successfully!
```

#### Test in Discord

Send a command in any channel where the bot has access:

```
!ping
```

Expected response:

```
🏓 Pong!
```

### Production Deployment

For production, consider:

1. **Environment variables**: Use secrets management (e.g., GitHub Secrets, AWS Secrets Manager)
2. **Process manager**: Use systemd, Docker, or Kubernetes
3. **Monitoring**: Add logging and metrics (Prometheus, Grafana)
4. **High availability**: Multiple bot instances with load balancing

---

## 📚 Available Commands

### General Commands

| Command | Description | Example |
|---------|-------------|---------|
| `!ping` | Check if bot is responding | `!ping` |
| `!hello` | Get a personalized greeting | `!hello` |
| `!help` | List all available commands | `!help` |
| `!info` | Show bot system information | `!info` |
| `!calc <expression>` | Calculate mathematical expression | `!calc 2 + 2 * 3` |

---

### Twitch Commands

Monitor Twitch streams and receive notifications in Discord when streamers go live!

#### Prerequisites

Make sure you have configured `TWITCH_CLIENT_ID` and `TWITCH_CLIENT_SECRET` in your `.env` file.

---

#### `!TwitchAddStream`

Add Twitch channels to your monitoring list.

**Syntax:**
```
!TwitchAddStream <channel1> [channel2] [channel3] ...
```

**Examples:**
```
!TwitchAddStream gaules
!TwitchAddStream gaules brtt loud_coringa
```

**Response:**
```
✅ Channels added successfully!

📺 Channels: gaules, brtt, loud_coringa
📋 Total monitored channels: 3
```

---

#### `!TwitchStreamMonitoring`

Start monitoring streams for a limited time.

**Syntax:**
```
!TwitchStreamMonitoring <duration_minutes> <check_interval_minutes>
```

**Parameters:**
- `duration_minutes`: Total monitoring duration
- `check_interval_minutes`: How often to check for live streams

**Examples:**
```
!TwitchStreamMonitoring 60 5
```
*Monitor for 60 minutes, checking every 5 minutes*

**Response:**
```
🚀 Twitch monitoring started!

⏱️ Duration: 60 minutes
🔄 Check interval: 5 minutes
📺 Monitored channels: 3
```

**Live Stream Notification:**
```
🎮 Gaules is LIVE!

📺 Title: VALORANT RANKED - RADIANT CLIMB
🎯 Playing: VALORANT
👥 Viewers: 12,450
⏱️ Live for: 1h 23min
🔗 Watch: https://twitch.tv/gaules
```

---

#### `!TwitchStreamMonitoringForever`

Start indefinite stream monitoring until manually stopped.

**Syntax:**
```
!TwitchStreamMonitoringForever <check_interval_minutes>
```

**Example:**
```
!TwitchStreamMonitoringForever 10
```
*Monitor indefinitely, checking every 10 minutes*

**Response:**
```
🚀 Twitch INFINITE monitoring started!

🔄 Check interval: 10 minutes
📺 Monitored channels: 3

⚠️ Use !TwitchStopMonitoring to stop
```

---

#### `!TwitchStopMonitoring`

Stop active stream monitoring.

**Syntax:**
```
!TwitchStopMonitoring
```

**Response:**
```
✅ Twitch monitoring stopped successfully!
```

---

### Twitch Features

#### Anti-Spam System
- Notifications sent once per stream every **30 minutes**
- Prevents duplicate notifications for the same stream

#### Notification Information
Each notification includes:
- Streamer name
- Stream title
- Game/category
- Viewer count
- Stream duration
- Direct link to stream

#### Monitoring Limits
- Only **one monitoring process** can be active at a time
- Attempting to start a new monitor while one is active will return an error

---

## 🏗️ Architecture

### System Overview

```
Discord User → Discord Bot → RabbitMQ → Event Dispatcher → Handlers → Services → Response
```

### Event Flow

1. **User sends command** in Discord (e.g., `!ping`)
2. **Discord integration** captures message and emits event `discord.command.received`
3. **CommandHandler** receives event and publishes to RabbitMQ
4. **RabbitMQ** queues the message
5. **Consumer** reads from queue and creates specific event (e.g., `discord.command.ping`)
6. **EventDispatcher** routes to appropriate handler
7. **Handler** processes command via **Service** layer
8. **Response** sent back to Discord

### Key Components

- **EventDispatcher**: Central event routing system
- **Handlers**: One file per command for modularity
- **Services**: Business logic layer
- **Integrations**: Discord, RabbitMQ, Twitch clients

For complete architecture documentation, see [AGENTS.md](AGENTS.md).

---

## 🛠️ Development

### Project Structure

```
enque-learning/
├── .github/
│   └── copilot-instructions.md    # AI development guidelines
├── cmd/
│   └── main.go                    # Application entry point
├── constants/
│   └── commands.go                # Command messages & constants
├── events/
│   ├── event.go                   # Event implementation
│   ├── dispatcher.go              # Event dispatcher
│   └── interfaces.go              # Event interfaces
├── handlers/
│   ├── command.go                 # Main command handler
│   ├── handlers.go                # Handler registration
│   ├── ping_pong.go               # Ping command
│   ├── hello.go                   # Hello command
│   ├── help.go                    # Help command
│   ├── info.go                    # Info command
│   ├── calc.go                    # Calculator command
│   ├── twitch_add_stream.go       # Twitch add stream
│   ├── twitch_stream_monitoring.go # Twitch monitoring
│   ├── twitch_stream_monitoring_forever.go
│   └── twitch_stop_monitoring.go  # Stop monitoring
├── integration/
│   ├── discord/                   # Discord client
│   ├── rabbitmq/                  # RabbitMQ client
│   └── twitch/                    # Twitch API client
├── service/
│   ├── service.go                 # Service struct
│   ├── calc.go                    # Calculator service
│   ├── hello.go                   # Hello service
│   ├── help.go                    # Help service
│   ├── info.go                    # Info service
│   ├── ping.go                    # Ping service
│   ├── twitch_add_stream.go       # Twitch add/remove
│   ├── twitch_stream_monitoring.go
│   ├── twitch_stream_monitoring_forever.go
│   ├── twitch_stop_monitoring.go
│   └── twitch_utils.go            # Shared Twitch utilities
├── server/
│   └── server.go                  # Server orchestration
├── internal/
│   ├── config/                    # Configuration loader
│   └── errors/                    # Custom errors
├── docker-compose.yml             # RabbitMQ setup
├── go.mod                         # Go dependencies
├── AGENTS.md                      # Complete architecture docs
└── README.md                      # This file
```

### Adding New Commands

See [.github/copilot-instructions.md](.github/copilot-instructions.md) for detailed guidelines on adding new commands.

**Quick Steps:**

1. Add constants in `constants/commands.go`
2. Create service file `service/<command>.go`
3. Create handler file `handlers/<command>.go`
4. Register handler in `handlers/handlers.go`

### Code Standards

- ✅ **ONE FILE PER COMMAND** - Mandatory pattern
- ✅ **English only** - All logs, errors, comments
- ✅ **Emoji prefixes** - For logs (✅ ❌ 🚀 🔍 🛑)
- ✅ **Constants** - No hardcoded messages
- ✅ **Custom errors** - Use `errors.New*()` constructors

### Error and Logging System

The project uses a custom error and logging system located in `utils/`.

#### Error Types

Create errors using specific constructors:

```go
// Validation errors (user input)
err := errors.NewValidation("missing required argument", nil)
err := errors.NewValidationf("invalid channel: %s", channel)

// Service errors (business logic)
err := errors.NewService("monitoring already in progress", nil)

// Integration errors (Discord, RabbitMQ, Twitch)
err := errors.NewIntegration("failed to send message", originalErr)

// API errors (external API calls)
err := errors.NewAPI("Twitch API request failed", err)

// Config errors (configuration issues)
err := errors.NewConfig("Discord token not found", nil)

// Network errors (connection issues)
err := errors.NewNetwork("connection timeout", err)
```

#### Error Levels

Errors are automatically categorized by severity:
- **DEBUG**: Development/debugging information (only shown with `DEBUG_MODE=true`)
- **INFO**: Validation errors, user input issues
- **WARN**: Service errors, API failures, recoverable issues
- **CRITICAL**: Config errors, auth failures, system errors

#### Logging

Use the logging functions instead of `log.Printf`:

```go
// Debug logs (only shown when DEBUG_MODE=true)
logger.Debug("Checking cache for key: %s", key)

// Info logs
logger.Info("✅ Command executed successfully")

// Warn logs
logger.Warn("⚠️ API rate limit approaching")

// Critical logs
logger.Critical("❌ Failed to connect to database")

// Auto-detect level from AppError
appErr := errors.NewAPI("request failed", err)
logger.LogError(appErr)  // Logs at WARN level
```

#### Debug Mode

Enable debug logging by setting `DEBUG_MODE=true` in `.env`:

```env
DEBUG_MODE=true
LOG_LEVEL=debug
```

Debug logs provide detailed information about:
- Variable values during execution
- Step-by-step flow tracking
- Cache operations
- Internal state changes

**See [utils/USAGE_EXAMPLES.go](utils/USAGE_EXAMPLES.go) for complete usage examples.**

### Build & Test

```bash
# Update dependencies
go mod tidy

# Build application
go build -o bot.exe cmd/main.go

# Run tests (when available)
go test ./...

# Check for errors
go vet ./...

# Format code
go fmt ./...
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

### Development Workflow

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Follow code standards**: See [.github/copilot-instructions.md](.github/copilot-instructions.md)
4. **Commit your changes**: `git commit -m "feat: add amazing feature"`
5. **Push to branch**: `git push origin feature/amazing-feature`
6. **Open a Pull Request**

### Commit Convention

Use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `refactor:` Code refactoring
- `test:` Adding tests
- `chore:` Maintenance tasks

### Code Review

All pull requests require:
- ✅ Passing build (`go build`)
- ✅ Code follows project patterns
- ✅ English-only code and comments
- ✅ No hardcoded credentials

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 📞 Support

### Documentation

- **Architecture**: [AGENTS.md](AGENTS.md) - Complete system documentation
- **Development Guidelines**: [.github/copilot-instructions.md](.github/copilot-instructions.md)
- **Implementation Plan**: [PLANO_DE_IMPLEMENTACAO.md](PLANO_DE_IMPLEMENTACAO.md) (Portuguese)

### Troubleshooting

#### Bot not responding
1. Verify `DISCORD_TOKEN` is correct
2. Check bot permissions in Discord server
3. Ensure bot has "Message Content Intent" enabled
4. Check RabbitMQ is running: `docker ps | grep rabbitmq`

#### Twitch notifications not working
1. Verify `TWITCH_CLIENT_ID` and `TWITCH_CLIENT_SECRET` are set
2. Confirm streamer names are correct (case-insensitive)
3. Remember the 30-minute anti-spam cooldown
4. Check bot logs for API errors

#### RabbitMQ connection failed
1. Ensure Docker is running: `docker-compose up -d`
2. Verify RabbitMQ URL in `.env`: `amqp://guest:guest@localhost:5672/`
3. Check RabbitMQ logs: `docker-compose logs rabbitmq`

### Getting Help

- **Issues**: [GitHub Issues](https://github.com/your-username/enque-learning/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/enque-learning/discussions)

---

## 🙏 Acknowledgments

- [discordgo](https://github.com/bwmarrin/discordgo) - Discord API wrapper
- [amqp091-go](https://github.com/rabbitmq/amqp091-go) - RabbitMQ client
- [godotenv](https://github.com/joho/godotenv) - Environment variable loader

---

**Made with ❤️ and Go**

*Last updated: March 2026*

