package twitch

import (
	"context"
	"encoding/json"
	"discordcommandbot/integration/twitch/models"
	"discordcommandbot/integration/twitch/utils"
	"discordcommandbot/internal/config"
	"io"
	"net/http"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

type Twitch struct {
	Config *config.TwitchConfig
	ctx    context.Context
}

func NewTwitchIntegration(ctx context.Context, config *config.TwitchConfig) (*Twitch, error) {
	return &Twitch{
		Config: config,
		ctx:    ctx,
	}, nil
}

func (t *Twitch) getToken() (string, error) {

	oauth2Config := &clientcredentials.Config{
		ClientID:     t.Config.ClientID,
		ClientSecret: t.Config.ClientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	token, err := oauth2Config.Token(t.ctx)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func (t *Twitch) GetStreams(streamChannels []string) (streamsResponse models.StreamResponse, err error) {

	token, err := t.getToken()
	if err != nil {
		return streamsResponse, err
	}

	req, err := http.NewRequest(http.MethodGet, utils.GetStreamURL(streamChannels), nil)
	if err != nil {
		return streamsResponse, err
	}

	req.Header.Set("Client-ID", t.Config.ClientID)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return streamsResponse, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return streamsResponse, err
	}

	err = json.Unmarshal(body, &streamsResponse)
	if err != nil {
		return streamsResponse, err
	}

	return streamsResponse, nil
}
