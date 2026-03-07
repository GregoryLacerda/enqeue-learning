package utils

import (
	"fmt"
	"strings"
)

func GetStreamURL(streamChannels []string) string {

	streams := strings.Join(streamChannels, "&user_login=")

	url := fmt.Sprintf("https://api.twitch.tv/helix/streams?user_login=%s", streams)

	return url
}
