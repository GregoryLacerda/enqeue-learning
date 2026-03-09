# DiscordCommandBot

Event-driven Discord bot in Go with RabbitMQ for asynchronous command processing.

## Overview

Discord commands are received, queued in RabbitMQ, and processed by command handlers through an event dispatcher.

Flow:

`Discord -> Command Event -> RabbitMQ -> Consumer -> Specific Handler -> Discord Response`

## Main Features

- Event-driven architecture with `EventDispatcher`
- Async command processing with RabbitMQ
- Modular command structure (one file per command)
- Twitch stream monitoring with Discord notifications
- Custom error and logging system (`pkg/errors`, `pkg/logger`)

## Prerequisites

- Go `1.26+`
- Docker + Docker Compose
- Discord bot token
- Twitch credentials (optional, only for Twitch commands)

## Quick Start

1. Clone and enter project:

```bash
git clone https://github.com/GregoryLacerda/DiscordCommandBot.git
cd DiscordCommandBot
```

2. Install dependencies:

```bash
go mod download
```

3. Start RabbitMQ:

```bash
docker-compose up -d
```

4. Create `.env` in project root:

```env
DISCORD_TOKEN=your_discord_bot_token
# Optional: Register slash commands instantly in one server
DISCORD_GUILD_ID=your_discord_server_id

RABBITMQ_URL=amqp://guest:guest@localhost:5672/
QUEUE_NAME=discord-commands
EXCHANGE_NAME=discord-exchange
ROUTING_KEY=discord.command

# Optional (Twitch commands)
TWITCH_CLIENT_ID=your_twitch_client_id
TWITCH_CLIENT_SECRET=your_twitch_client_secret

# Optional
DEBUG_MODE=false
LOG_LEVEL=info
WEB_SERVER_PORT=8080
```

5. Run the bot:

```bash
go run cmd/main.go
```

## RabbitMQ Local Access

- AMQP: `amqp://localhost:5672`
- UI: `http://localhost:15672`
- User/Pass: `guest` / `guest`

## Commands

The bot uses slash commands (`/`) only.

### General

- `/ping` - Checks if the bot is online
- `/hello` - Sends a greeting
- `/help` - Lists available commands
- `/info` - Shows user/system information
- `/calc expression:<text>` - Calculates an expression

Examples:

```text
/ping
/calc expression:2 + 2 * 3
```

### Twitch

- `/twitch add channels:<channel1 channel2 ...> [ephemeral:true|false]`
- `/twitch list [ephemeral:true|false]`
- `/twitch clear [ephemeral:true|false]`
- `/twitch start duration_minutes:<int> check_interval_minutes:<int> [ephemeral:true|false]`
- `/twitch startforever check_interval_minutes:<int> [ephemeral:true|false]`
- `/twitch stop [ephemeral:true|false]`

Notes:

- Requires `TWITCH_CLIENT_ID` and `TWITCH_CLIENT_SECRET`
- Anti-spam cooldown for repeated notifications is enabled
- Only one monitoring process can run at a time
- Slash command descriptions and option hints are shown directly in Discord while typing
- `/twitch add` supports channel autocomplete while typing
- `ephemeral` is optional for Twitch admin actions and controls if the command ack is private or public
- For instant slash updates in your server, set `DISCORD_GUILD_ID` in `.env`

## Project Structure

```text
cmd/              # app entrypoint
constants/        # command texts and constants
events/           # event contracts + dispatcher
handlers/         # command handlers
integration/      # Discord, RabbitMQ, Twitch clients
service/          # business logic layer
server/           # startup and graceful shutdown
pkg/errors/       # typed app errors
pkg/logger/       # logging utilities
```

## Development Rules (Important)

- One file per command in `handlers/` and `service/`
- No business logic in handlers (use service layer)
- Use `pkg/errors` constructors instead of `fmt.Errorf`
- Use `pkg/logger` instead of `log.Printf`
- Keep logs, errors, comments in English

## Build and Validation

```bash
go mod tidy
go build -o bot.exe cmd/main.go
go test ./...
```

## Troubleshooting

### Bot does not respond

- Validate `DISCORD_TOKEN`
- Check bot permissions on the server/channel
- Confirm RabbitMQ is running

### RabbitMQ connection errors

- Start services: `docker-compose up -d`
- Check logs: `docker-compose logs -f rabbitmq`
- Confirm `RABBITMQ_URL` in `.env`

### Twitch monitoring issues

- Confirm Twitch credentials in `.env`
- Check channel names passed to commands
- Check bot logs for API/integration errors

## Documentation

- Architecture and deep context: `AGENTS.md`
- Copilot coding rules: `.github/copilot-instructions.md`
- Error/log usage examples: `utils/USAGE_EXAMPLES.go`

## License

MIT - see `LICENSE`.