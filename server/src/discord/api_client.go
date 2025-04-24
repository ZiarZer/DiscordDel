package discord

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

func getUserGuilds(authorizationToken string) (*http.Response, error) {
	return request("GET", "users/@me/guilds", authorizationToken)
}

func getGuildById(guildId string, authorizationToken string) (*http.Response, error) {
	return request("GET", fmt.Sprintf("guilds/%s", guildId), authorizationToken)
}

func getGuildChannels(guildId string, authorizationToken string) (*http.Response, error) {
	return request("GET", fmt.Sprintf("guilds/%s/channels", guildId), authorizationToken)
}

func getChannelById(channelId string, authorizationToken string) (*http.Response, error) {
	return request("GET", fmt.Sprintf("channels/%s", channelId), authorizationToken)
}

type GetChannelMessagesOptions struct {
	Limit  *int
	Before *string
	After  *string
	Around *string
}

func getChannelMessages(channelId string, options *GetChannelMessagesOptions, authorizationToken string) (*http.Response, error) {
	searchParams := url.Values{}
	if options != nil {
		if options.Limit != nil {
			searchParams.Add("limit", strconv.Itoa(*options.Limit))
		}
		if options.Before != nil {
			searchParams.Add("before", *options.Before)
		}
		if options.After != nil {
			searchParams.Add("after", *options.After)
		}
		if options.Around != nil {
			searchParams.Add("around", *options.Around)
		}
	}
	return request("GET", fmt.Sprintf("channels/%s/messages?%s", channelId, searchParams.Encode()), authorizationToken)
}
