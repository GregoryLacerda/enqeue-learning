
# Agent Context - enque-learning Project

## 📋 Overview

**enque-learning** is an event-driven Discord bot written in Go that uses RabbitMQ for asynchronous command processing. The project implements a robust event system based on the Event Dispatcher pattern, allowing extensibility and separation of concerns.

### Purpose
- Receive commands from Discord (e.g., `!ping`, `!help`, `!calc`)
- Enqueue commands in RabbitMQ for asynchronous processing
- Process commands through specific handlers
- Respond to Discord with the processing result

### Core Technologies
- **Language**: Go 1.26
- **Discord Bot**: discordgo (v0.29.0)
- **Message Broker**: RabbitMQ with amqp091-go (v1.10.0)
- **Identifiers**: UUID (google/uuid v1.6.0)
- **Configuration**: godotenv (v1.5.1)

---

## 🏗️ Architecture

### Complete Execution Flow

```
┌─────────────────┐
│  Discord User   │
│  sends: !ping   │
└────────┬────────┘
         │
         ▼
┌────────────────────────────────────────────────────────┐
│                    Discord Integration                  │
│  - Captures message with prefix "!"                    │
│  - Creates DiscordCommandPayload                       │
│  - Emits event: "discord.command.received"             │
└────────┬───────────────────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────────────────────┐
│              Event Dispatcher                           │
│  - Receives event "discord.command.received"           │
│  - Finds registered handlers                           │
│  - Calls CommandHandler.HandleEvent()                  │
└────────┬───────────────────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────────────────────┐
│               Command Handler                           │
│  - Serializes payload to JSON                          │
│  - Publishes to RabbitMQ (Exchange + Routing Key)      │
└────────┬───────────────────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────────────────────┐
│                  RabbitMQ Queue                         │
│  - Queue: discord-commands                             │
│  - Exchange: discord-exchange (type: topic)            │
│  - Routing Key: discord.command                        │
└────────┬───────────────────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────────────────────┐
│             RabbitMQ Consumer                           │
│  - Consumes messages from queue                        │
│  - Deserializes JSON to DiscordCommandPayload          │
│  - Calls ResponseHandler.ProcessMessage()              │
└────────┬───────────────────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────────────────────┐
│              Response Handler                           │
│  - Creates specific event: "discord.command.ping"      │
│  - Dispatch to EventDispatcher                         │
│  - If no handler → UnknownCommandHandler              │
└────────┬───────────────────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────────────────────┐
│           Command-Specific Handler                      │
│  Examples:                                             │
│  - PingCommandHandler → responds "🏓 Pong!"           │
│  - HelpCommandHandler → lists commands                 │
│  - CalcCommandHandler → calculates expression          │
└────────┬───────────────────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────────────────────┐
│            Discord Integration                          │
│  - ReplyToMessage(channelID, messageID, response)     │
│  - Sends response to Discord channel                   │
└────────────────────────────────────────────────────────┘
```

#### 1. **Event System** (`events/`)
- **EventInterface**: Contract for events (ID, Name, Date, Payload)
- **Event**: Concrete implementation of events with UUID
- **EventHandlerInterface**: Contract for handlers (`HandleEvent()`)
- **EventDispatcher**: Central registry that maps events to handlers
  - Supports multiple handlers per event
  - Parallel execution with WaitGroups
  - Methods: `RegisterHandler()`, `Dispatch()`, `RemoveHandler()`

#### 2. **Integrations** (`integration/`)
- **Discord**: 
  - Manages Discord connection via discordgo
  - Captures messages with prefix (default: `!`)
  - Emits `discord.command.received` events
  - Methods: `SendMessage()`, `ReplyToMessage()`
  
- **RabbitMQ**:
  - Manages RabbitMQ connection
  - Setup of Exchange (topic), Queue and Binding
  - Methods: `Publisher(body)`, `Consumer()`, `Close()`

- **Twitch**: Prepared structure (not fully implemented)

#### 3. **Handlers** (`handlers/`)
Two-level system:

**Level 1: Main Handlers**
- `CommandHandler`: Receives commands from Discord → publishes to RabbitMQ
- `ResponseHandler`: Consumes from RabbitMQ → dispatches specific events

**Level 2: Command-Specific Handlers**
- `PingCommandHandler`: Responds "🏓 Pong!"
- `HelloCommandHandler`: Personalized greeting
- `HelpCommandHandler`: Lists available commands
- `InfoCommandHandler`: Shows system information
- `CalcCommandHandler`: Calculates mathematical expressions
- `UnknownCommandHandler`: Fallback for unknown commands

#### 4. **Server** (`server/`)
- Orchestrates initialization of all components
- Registers handlers in EventDispatcher
- Starts Discord and RabbitMQ Consumer
- Manages graceful shutdown (SIGTERM, SIGINT)

#### 5. **Service** (`service/`)
- Business logic layer
- Currently contains `HelpService` to generate command list

#### 6. **Configuration** (`internal/config/`)
- Loads environment variables via `.env`
- `Config` struct with typed fields:
  - Discord: Token, CommandPrefix
  - RabbitMQ: URL, QueueName, ExchangeName, RoutingKey
  - Server: WebServerPort

---

## 📁 Detailed Directory Structure

```
enque-learning/
│
├── cmd/
│   └── main.go                    # Entry point: initializes complete system
│
├── events/                        # Event system (core)
│   ├── event.go                   # Event implementation with UUID
│   ├── interfaces.go              # EventInterface, EventHandlerInterface
│   ├── dispatcher.go              # EventDispatcher (registry + dispatch)
│   └── discord_payload.go         # DiscordCommandPayload structure
│
├── integration/                   # External integrations
│   ├── integration.go             # Factory: NewIntegrations()
│   ├── discord/
│   │   ├── discord.go             # Discord client + messageHandler
│   │   └── models.go              # DiscordCommandPayload
│   ├── rabbitmq/
│   │   └── rabbitmq.go            # Publisher + Consumer + Setup
│   └── twitch/
│       └── twitch.go              # Prepared for Twitch (future)
│
├── handlers/                      # Command handlers
│   ├── command.go                 # CommandHandler (Discord → RabbitMQ)
│   ├── handlers.go                # ResponseHandler (RabbitMQ → Dispatch)
│   ├── ping_pong.go               # PingCommandHandler
│   ├── hello.go                   # HelloCommandHandler
│   ├── help.go                    # HelpCommandHandler
│   ├── info.go                    # InfoCommandHandler
│   ├── calc.go                    # CalcCommandHandler
│   └── unknown.go                 # UnknownCommandHandler
│
├── service/                       # Business logic
│   ├── service.go                 # Service struct
│   └── help.go                    # HelpService (available commands)
│
├── server/                        # Orchestration
│   └── server.go                  # Server: StartAll(), Shutdown()
│
├── internal/                      # Internal packages
│   ├── config/
│   │   └── config.go              # LoadConfig() + Config struct
│   └── errors/
│       └── errors.go              # Custom errors (AlreadyRegisteredError)
│
├── utils/                         # Utilities
│   └── utils.go
│
├── scripts/                       # Auxiliary scripts
│   └── run-local.sh               # Script for local development
│
├── docker-compose.yml             # RabbitMQ + Management UI
├── go.mod                         # Go dependencies
├── PLANO_DE_IMPLEMENTACAO.md      # Documentation: implementation plan
└── MELHORIA_USANDO_EVENTOS.md     # Documentation: event system improvements
```

---

## 🔑 Key Components and Interfaces

### EventInterface
```go
type EventInterface interface {
    GetName() string
    GetDate() string
    GetID() string
    GetPayload() interface{}
}
```

### EventHandlerInterface
```go
type EventHandlerInterface interface {
    HandleEvent(event EventInterface) error
}
```

### DiscordCommandPayload
```go
type DiscordCommandPayload struct {
    UserID    string   // Discord User ID
    Username  string   // Discord Username
    ChannelID string   // Channel where command was sent
    GuildID   string   // Discord Server (Guild)
    Command   string   // Command (e.g., "ping", "help")
    Arguments []string // Command arguments
    MessageID string   // Original message ID
    Timestamp string   // Command timestamp
}
```

---

## 🎯 Patterns and Conventions

### Event Naming

**System Events:**
- `discord.command.received` - Command received from Discord (before enqueuing)

**Command-Specific Events:**
- `discord.command.<command_name>` - Command being processed
  - Example: `discord.command.ping`
  - Example: `discord.command.help`
  - Example: `discord.command.calc`

### Handler Structure

All command handlers follow this pattern:

```go
type XCommandHandler struct {
    Discord *discord.Discord
    // Other dependencies (service, etc)
}

func NewXCommandHandler(discord *discord.Discord) *XCommandHandler {
    return &XCommandHandler{Discord: discord}
}

func (h *XCommandHandler) HandleEvent(event events.EventInterface) error {
    // 1. Cast the payload
    payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
    if !ok {
        return fmt.Errorf("invalid payload")
    }
    
    // 2. Log for tracing
    log.Printf("handling X command from user: %s", payload.Username)
    
    // 3. Command logic
    response := processLogic(payload)
    
    // 4. Respond to Discord
    err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
    if err != nil {
        return fmt.Errorf("failed to send response: %w", err)
    }
    
    return nil
}
```

### Handler Registration

Handlers are registered in `ResponseHandler.NewResponseHandler()`:

```go
dispatcher.RegisterHandler("discord.command.ping", NewPingCommandHandler(discord))
dispatcher.RegisterHandler("discord.command.help", NewHelpCommandHandler(discord, service))
// ... other commands
```

---

## 🚀 How to Add a New Command

### Step by Step

#### 1. Create new file in `handlers/`
```go
// handlers/mytask.go
package handlers

import (
    "enque-learning/events"
    "enque-learning/integration/discord"
    "fmt"
    "log"
)

type MyTaskCommandHandler struct {
    Discord *discord.Discord
}

func NewMyTaskCommandHandler(discord *discord.Discord) *MyTaskCommandHandler {
    return &MyTaskCommandHandler{Discord: discord}
}

func (h *MyTaskCommandHandler) HandleEvent(event events.EventInterface) error {
    payload, ok := event.GetPayload().(discord.DiscordCommandPayload)
    if !ok {
        return fmt.Errorf("invalid payload")
    }
    
    log.Printf("handling mytask command from user: %s", payload.Username)
    
    // Your logic here
    response := "✅ Task executed successfully!"
    
    err := h.Discord.ReplyToMessage(payload.ChannelID, payload.MessageID, response)
    if err != nil {
        return fmt.Errorf("failed to send response: %w", err)
    }
    
    return nil
}
```

#### 2. Register handler in `handlers/handlers.go`

In the `NewResponseHandler()` method, add:
```go
dispatcher.RegisterHandler("discord.command.mytask", NewMyTaskCommandHandler(discord))
```

#### 3. (Optional) Add to Help

If you have `HelpService` in `service/help.go`, add command to the list.

#### 4. Test

In Discord: `!mytask`

---

## ⚙️ Configuration

### Environment Variables (.env)

```env
# Discord
DISCORD_TOKEN=your_bot_token_here
DISCORD_COMMAND_PREFIX=!

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
QUEUE_NAME=discord-commands
EXCHANGE_NAME=discord-exchange
ROUTING_KEY=discord.command

# Server (optional)
WEB_SERVER_PORT=8080
```

### Docker Compose (RabbitMQ)

```yaml
services:
  rabbitmq:
    image: rabbitmq:3.12-management
    container_name: rabbitmq
    ports:
      - "5672:5672"   # AMQP
      - "15672:15672" # Management UI
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
```

**Start RabbitMQ:**
```bash
docker-compose up -d
```

**Access Management UI:**
http://localhost:15672 (guest/guest)

---

## 🔧 Useful Commands

### Development

```bash
# Install dependencies
go mod download

# Run application
go run cmd/main.go

# Build
go build -o bot cmd/main.go

# Execute
./bot

# Tests (if any)
go test ./...
```

### Docker

```bash
# Start RabbitMQ
docker-compose up -d

# View RabbitMQ logs
docker-compose logs -f rabbitmq

# Stop RabbitMQ
docker-compose down

# Clean volumes (complete reset)
docker-compose down -v
```

---

## 📊 Detailed Data Flow

### 1. Command Reception

```
User types: !ping in Discord
  ↓
discord.messageHandler() captures
  ↓
Checks prefix "!"
  ↓
Parse: command="ping", args=[]
  ↓
Creates DiscordCommandPayload
  ↓
Creates Event("discord.command.received")
  ↓
dispatcher.Dispatch(event)
```

### 2. Initial Processing (CommandHandler)

```
CommandHandler.HandleEvent() receives
  ↓
Serializes payload to JSON
  ↓
RabbitMQ.Publisher(jsonData)
  ↓
Message goes to Exchange → Queue
```

### 3. Consumption and Processing (ResponseHandler)

```
RabbitMQ.Consumer() receives message
  ↓
server.go: go func() processes
  ↓
ResponseHandler.ProcessMessage(body)
  ↓
Deserializes JSON → DiscordCommandPayload
  ↓
Creates Event("discord.command.ping")
  ↓
dispatcher.Dispatch(event)
```

### 4. Specific Handler Execution

```
PingCommandHandler.HandleEvent() called
  ↓
response = "🏓 Pong!"
  ↓
Discord.ReplyToMessage(channelID, messageID, response)
  ↓
User sees response in Discord
  ↓
msg.Ack(false) - confirms processing
```

---

## 🧩 Project Dependencies

### Direct Dependencies

```go
require (
    github.com/bwmarrin/discordgo v0.29.0      // Discord client
    github.com/google/uuid v1.6.0               // UUID generation
    github.com/joho/godotenv v1.5.1             // .env loading
    github.com/rabbitmq/amqp091-go v1.10.0      // RabbitMQ client
)
```

### Indirect Dependencies

```go
require (
    github.com/gorilla/websocket v1.4.2                // WebSocket (discordgo)
    golang.org/x/crypto v0.0.0-20210421170649-...      // Crypto (discordgo)
    golang.org/x/sys v0.0.0-20201119102817-...         // System calls
)
```

---

## 🎭 States and Lifecycle

### Initialization (main.go → server.go)

1. **LoadConfig()** - Loads .env
2. **NewEventDispatcher()** - Creates event registry
3. **NewIntegrations()** - Initializes Discord, RabbitMQ, Twitch
4. **NewService()** - Initializes business services
5. **NewServer()** - Creates orchestrator
   - Creates CommandHandler
   - Creates ResponseHandler (registers specific handlers)
6. **server.StartAll()**
   - Registers `discord.command.received` → CommandHandler
   - Starts Discord session (`.Open()`)
   - Starts RabbitMQ Consumer
   - Message processing loop
   - Waits for shutdown signal

### Graceful Shutdown

```
SIGTERM/SIGINT received
  ↓
server.waitForShutdown() unblocks
  ↓
server.Shutdown()
  ↓
Discord.Stop() - closes session
  ↓
RabbitMQ.Close() - closes channel and connection
  ↓
Log: "✅ System shut down successfully!"
```

---

## 🐛 Debugging and Logs

### Important Logs

The system uses `log.Printf()` and `log.Println()` extensively:

**Discord:**
- `"bot Discord conectado e online!"`
- `"erro ao processar comando: %v"`

**RabbitMQ:**
- `"RabbitMQ sucess configured"`
- `"📨 Mensagem recebida da fila"`
- `"✅ Mensagem processada com sucesso"`
- `"❌ Erro ao processar mensagem: %v"`

**Server:**
- `"🚀 Iniciando sistema completo (Producer + Consumer)..."`
- `"✅ Sistema completo iniciado com sucesso!"`
- `"⚠️ Sinal de interrupção recebido..."`
- `"🛑 Encerrando sistema..."`

**Handlers:**
- `"processing command: %s"`
- `"handling X command from user: %s"`
- `"processing response to Discord channel %s: %s"`

### RabbitMQ Management

Access http://localhost:15672 to:
- View number of messages in queue
- Monitor active consumers
- Inspect message payloads
- View publish/consume rates

---

## 📚 Additional Documentation

The project includes documentation in Portuguese:

- **PLANO_DE_IMPLEMENTACAO.md**: Step-by-step guide for original implementation
- **MELHORIA_USANDO_EVENTOS.md**: Explanation of migration to event system

---

## 🔮 Future Extensions

### Areas Prepared for Expansion

1. **Twitch Integration** - Structure already exists in `integration/twitch/`
2. **Web Server** - Config has `WebServerPort` (not implemented)
3. **More Commands** - System is plug-and-play via handlers
4. **Persistence** - Add database for history
5. **Metrics** - Prometheus for monitoring
6. **Tests** - Structure facilitates unit tests with mocks

---

## 💡 Tips for AIs

### When Modifying Code

1. **Always implement EventHandlerInterface**: `HandleEvent(event EventInterface) error`
2. **Register handlers in ResponseHandler**: Use `dispatcher.RegisterHandler()`
3. **Use log.Printf for tracing**: Facilitates debugging
4. **Error wrapping**: Use `fmt.Errorf("context: %w", err)` for stack traces
5. **Payload type assertion**: Always check `ok` when casting

### When Adding Features

1. **New commands**: Create handler in `handlers/`, register in `handlers.go`
2. **New integrations**: Add in `integration/`, update `integration.go`
3. **New business logic**: Create in `service/`, inject via constructor
4. **New events**: Follow pattern `<source>.<action>.<detail>`

### When Debugging

1. **Check .env**: Make sure all variables are defined
2. **RabbitMQ running**: `docker ps | grep rabbitmq`
3. **Valid Discord token**: Common error on initialization
4. **Consumer logs**: Messages being processed but no response?
5. **Handler registered**: Check `NewResponseHandler()` has event registration

---

## 🎯 Currently Implemented Commands

| Command | Handler | Description | Aliases |
|---------|---------|-------------|--------|
| `!ping` | PingCommandHandler | Tests if bot is responding | - |
| `!hello` | HelloCommandHandler | Personalized greeting | `!oi`, `!olá` |
| `!help` | HelpCommandHandler | Lists available commands | `!ajuda` |
| `!info` | InfoCommandHandler | System information | - |
| `!calc` | CalcCommandHandler | Calculates mathematical expression | - |
| _unknown_ | UnknownCommandHandler | Fallback for invalid commands | - |

---

## 📞 Important Discord Integration Methods

```go
// Send simple message
Discord.SendMessage(channelID string, message string) error

// Reply to a specific message (with reference)
Discord.ReplyToMessage(channelID string, messageID string, message string) error
```

---

## 🔄 Concurrency and WaitGroups

The `EventDispatcher` processes handlers in parallel using `sync.WaitGroup`:

```go
func (e *EventDispatcher) Dispatch(event EventInterface) error {
    if handlers, ok := e.handlers[event.GetName()]; ok {
        wg := sync.WaitGroup{}
        wg.Add(len(handlers))
        errorChannel := make(chan error, len(handlers))
        
        for _, handler := range handlers {
            wg.Go(func() {
                defer wg.Done()
                err := handler.HandleEvent(event)
                if err != nil {
                    errorChannel <- err
                }
            })
        }
        
        wg.Wait()
        // ... error handling
    }
    return nil
}
```

**Implication**: Handlers must be thread-safe if they share state.

---

## 🏁 Quick Start

```bash
# 1. Clone repo (example)
git clone <repo-url>
cd enque-learning

# 2. Create .env
cp .env.example .env
# Edit with your DISCORD_TOKEN

# 3. Start RabbitMQ
docker-compose up -d

# 4. Install dependencies
go mod download

# 5. Run bot
go run cmd/main.go

# 6. Test in Discord
# !ping
# !help
# !hello
```

---

## 📝 Final Notes

This project is an excellent example of:
- ✅ **Event-Driven Architecture** in Go
- ✅ **Pub/Sub Pattern** with RabbitMQ
- ✅ **Dependency Injection** manual
- ✅ **Interface-based Design** for testability
- ✅ **Graceful Shutdown** with signal handling
- ✅ **Separation of Concerns** clear between layers

**Last Update**: March 2026 (Context date)

---

_This document was created to facilitate code understanding by AI systems and developers. Keep it updated as the project evolves._
