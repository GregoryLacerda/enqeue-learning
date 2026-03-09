package constants

const (
	// Help Command
	CommandHelpMessage = `📚 **Available Commands:**

				**General Commands:**
				**/ping** - Tests if the bot is responding
				**/hello** - Receives a greeting
				**/calc expression:<text>** - Calculates a mathematical expression (e.g., /calc expression:2 + 2)
				**/info** - Shows information about you and the system
				**/help** - Shows this help message

				**Twitch Commands:**
				**/twitch add channels:<channel1 channel2 ...>** - Adds channels for monitoring
				**/twitch list** - Lists monitored channels
				**/twitch clear** - Clears all monitored channels
				**/twitch start duration_minutes:<int> check_interval_minutes:<int>** - Monitors streams for a limited time
				**/twitch startforever check_interval_minutes:<int>** - Monitors streams indefinitely
				**/twitch stop** - Stops active monitoring`

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
	UnknownCommandTemplate = "❓ Unknown command: `%s`\n\nUse `/help` to see available commands."

	// Calc Command
	CalcUsageMessage        = "❌ Usage: `/calc expression:<text>`\n**Example:** /calc expression:2 + 2"
	CalcErrorTemplate       = "❌ Error calculating: %s"
	CalcResultTemplate      = "🧮 **Result:**\n`%s = %.2f`"
	CalcInvalidFormat       = "expected format: number operator number (e.g., 2 + 2)"
	CalcInvalidFirstNumber  = "invalid first number: %s"
	CalcInvalidSecondNumber = "invalid second number: %s"
	CalcDivisionByZero      = "division by zero"
	CalcInvalidOperator     = "invalid operator: %s (use: +, -, *, /, ^)"

	// Twitch Commands
	TwitchAddStreamUsage                 = "❌ Usage: `/twitch add channels:<channel1 channel2 ...>`\n**Example:** /twitch add channels:gaules brtt"
	TwitchAddStreamSuccess               = "✅ **Channels added successfully!**\n\n📺 Channels: %s\n📋 Total monitored channels: %d"
	TwitchListStreamsEmpty               = "📭 **No channels configured for monitoring.**\n\nUse `/twitch add` to add channels."
	TwitchListStreamsSuccess             = "📋 **Monitored Twitch channels (%d):**\n\n%s"
	TwitchClearStreamsSuccess            = "🧹 **Channel list cleared successfully!**\n\n📋 Removed channels: %d"
	TwitchStreamMonitoringUsage          = "❌ Usage: `/twitch start duration_minutes:<int> check_interval_minutes:<int>`\n**Example:** /twitch start duration_minutes:60 check_interval_minutes:5"
	TwitchStreamMonitoringStarted        = "🚀 **Twitch monitoring started!**\n\n⏱️ Duration: %d minutes\n🔄 Check interval: %d minutes\n📺 Monitored channels: %d"
	TwitchStreamMonitoringError          = "❌ Error starting monitoring: %s"
	TwitchStreamMonitoringForeverUsage   = "❌ Usage: `/twitch startforever check_interval_minutes:<int>`\n**Example:** /twitch startforever check_interval_minutes:10"
	TwitchStreamMonitoringForeverStarted = "🚀 **Twitch INFINITE monitoring started!**\n\n🔄 Check interval: %d minutes\n📺 Monitored channels: %d\n\n⚠️ Use `/twitch stop` to stop"
	TwitchStopMonitoringSuccess          = "✅ **Twitch monitoring stopped successfully!**"
	TwitchStopMonitoringError            = "❌ Error stopping monitoring: %s"
	TwitchNoChannelsAdded                = "❌ No channels added yet!\n\nUse `/twitch add` to add channels."
	TwitchInvalidDuration                = "❌ Invalid duration. Use an integer number of minutes."
	TwitchInvalidInterval                = "❌ Invalid interval. Use an integer number of minutes."
)
