package discord

// DiscordCommandPayload represents the payload of a Discord command
type DiscordCommandPayload struct {
	UserID    string            `json:"user_id"`
	Username  string            `json:"username"`
	ChannelID string            `json:"channel_id"`
	GuildID   string            `json:"guild_id"`
	Command   string            `json:"command"`
	Arguments []string          `json:"arguments"`
	MessageID string            `json:"message_id"`
	Timestamp string            `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// DiscordResponsePayload represents the response to be sent to Discord
type DiscordResponsePayload struct {
	ChannelID string `json:"channel_id"`
	Message   string `json:"message"`
	MessageID string `json:"message_id,omitempty"` // To reply to a specific message
}
