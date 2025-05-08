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

type SearchChannelThreadsOptions struct {
	Offset    int
	Limit     *int
	Archived  *bool
	SortBy    *string
	SortOrder *string
}

func searchChannelThreads(authorizationToken string, mainChannelId string, options *SearchChannelThreadsOptions) (*http.Response, error) {
	searchParams := url.Values{}
	if options != nil {
		searchParams.Add("offset", strconv.Itoa(options.Offset))
		if options.Limit != nil {
			searchParams.Add("limit", strconv.Itoa(*options.Limit))
		}
		if options.Archived != nil {
			searchParams.Add("archived", strconv.FormatBool(*options.Archived))
		}
		if options.SortBy != nil {
			searchParams.Add("sort_by", *options.SortBy)
		}
		if options.SortOrder != nil {
			searchParams.Add("sort_order", *options.SortOrder)
		}
	}
	return request("GET", fmt.Sprintf("/channels/%s/threads/search?%s", mainChannelId, searchParams.Encode()), authorizationToken)
}
