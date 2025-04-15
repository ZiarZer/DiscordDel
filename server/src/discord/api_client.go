package discord

import (
	"fmt"
	"net/http"
)

const DiscordApibaseURL = "https://discord.com/api/v9"

func request(method string, endpoint string, authorizationToken string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", DiscordApibaseURL, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authorizationToken)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func login(authorizationToken string) (*http.Response, error) {
	return request("GET", "users/@me", authorizationToken)
}

func getGuildById(guildId string, authorizationToken string) (*http.Response, error) {
	return request("GET", fmt.Sprintf("guilds/%s", guildId), authorizationToken)
}

func getChannelById(channelId string, authorizationToken string) (*http.Response, error) {
	return request("GET", fmt.Sprintf("channels/%s", channelId), authorizationToken)
}
