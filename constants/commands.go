package constants

const (
	// Help Command
	CommandHelpMessage = `📚 **Available Commands:**

				**General Commands:**
				**!ping** - Tests if the bot is responding
				**!hello** / **!oi** - Receives a greeting
				**!calc <expression>** - Calculates a mathematical expression (e.g., !calc 2 + 2)
				**!info** - Shows information about you and the system
				**!help** - Shows this help message

				**Twitch Commands:**
				**!TwitchAddStream <channel1> [channel2]...** - Adds channels for monitoring
				**!TwitchStreamMonitoring <duration_min> <interval_min>** - Monitors streams for a limited time
				**!TwitchStreamMonitoringForever <interval_min>** - Monitors streams indefinitely
				**!TwitchStopMonitoring** - Stops active monitoring`

	// Ping Command
	PingMessage = "🏓 Pong!"

	// Hello Command
	HelloMessageTemplate = "👋 Hello, %s! How can I help?"

	// Info Command
	InfoMessageTemplate = `ℹ️ **System Information:**

👤 **User:** %s
🆔 **User ID:** %s
📝 **Command:** %s
📅 **Channel ID:** %s
🏢 **Guild ID:** %s
⏰ **Timestamp:** %s`

	// Unknown Command
	UnknownCommandTemplate = "❓ Unknown command: `%s`\n\nUse `!help` to see available commands."

	// Calc Command
	CalcUsageMessage        = "❌ Usage: `!calc <expression>`\n**Example:** !calc 2 + 2"
	CalcErrorTemplate       = "❌ Error calculating: %s"
	CalcResultTemplate      = "🧮 **Result:**\n`%s = %.2f`"
	CalcInvalidFormat       = "expected format: number operator number (e.g., 2 + 2)"
	CalcInvalidFirstNumber  = "invalid first number: %s"
	CalcInvalidSecondNumber = "invalid second number: %s"
	CalcDivisionByZero      = "division by zero"
	CalcInvalidOperator     = "invalid operator: %s (use: +, -, *, /, ^)"

	// Twitch Commands
	TwitchAddStreamUsage                 = "❌ Usage: `!TwitchAddStream <channel1> [channel2] [channel3]...`\n**Example:** !TwitchAddStream gaules brtt"
	TwitchAddStreamSuccess               = "✅ **Channels added successfully!**\n\n📺 Channels: %s\n📋 Total monitored channels: %d"
	TwitchStreamMonitoringUsage          = "❌ Usage: `!TwitchStreamMonitoring <duration_minutes> <interval_minutes>`\n**Example:** !TwitchStreamMonitoring 60 5"
	TwitchStreamMonitoringStarted        = "🚀 **Twitch monitoring started!**\n\n⏱️ Duration: %d minutes\n🔄 Check interval: %d minutes\n📺 Monitored channels: %d"
	TwitchStreamMonitoringError          = "❌ Error starting monitoring: %s"
	TwitchStreamMonitoringForeverUsage   = "❌ Usage: `!TwitchStreamMonitoringForever <interval_minutes>`\n**Example:** !TwitchStreamMonitoringForever 10"
	TwitchStreamMonitoringForeverStarted = "🚀 **Twitch INFINITE monitoring started!**\n\n🔄 Check interval: %d minutes\n📺 Monitored channels: %d\n\n⚠️ Use `!TwitchStopMonitoring` to stop"
	TwitchStopMonitoringSuccess          = "✅ **Twitch monitoring stopped successfully!**"
	TwitchStopMonitoringError            = "❌ Error stopping monitoring: %s"
	TwitchNoChannelsAdded                = "❌ No channels added yet!\n\nUse `!TwitchAddStream <channel>` to add channels."
	TwitchInvalidDuration                = "❌ Invalid duration. Use an integer number of minutes."
	TwitchInvalidInterval                = "❌ Invalid interval. Use an integer number of minutes."
)
