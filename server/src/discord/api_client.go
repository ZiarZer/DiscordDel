package discord

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

const DiscordApibaseURL = "https://discord.com/api/v9"

type ApiClient struct {
	Delay               int
	lastRequestUnixTime int64
}

func (apiClient *ApiClient) request(method string, endpoint string, authorizationToken string, retriesLeft int) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", DiscordApibaseURL, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authorizationToken)
	httpClient := &http.Client{}

	timeSinceLastRequest := time.Now().UnixMilli() - int64(apiClient.lastRequestUnixTime)
	if timeSinceLastRequest < int64(apiClient.Delay) {
		time.Sleep(time.Duration(apiClient.Delay*1000000) - time.Duration(timeSinceLastRequest*1000000))
	} else {
		utils.InternalLog(fmt.Sprintf("%d", timeSinceLastRequest), nil)
	}
	resp, err := httpClient.Do(req)
	apiClient.lastRequestUnixTime = time.Now().UnixMilli()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 429 && retriesLeft > 0 {
		apiClient.Delay = int(math.Round(1.5 * float64(apiClient.Delay)))
		utils.InternalLog(fmt.Sprintf("Rate limited - Delay multiplied by 1.5, current value: %dms", apiClient.Delay), utils.INFO)
		return apiClient.request(method, endpoint, authorizationToken, retriesLeft-1)
	}
	return resp, nil
}

func (apiClient *ApiClient) login(authorizationToken string) (*http.Response, error) {
	return apiClient.request("GET", "users/@me", authorizationToken, 3)
}

func (apiClient *ApiClient) getUserGuilds(authorizationToken string) (*http.Response, error) {
	return apiClient.request("GET", "users/@me/guilds", authorizationToken, 3)
}

func (apiClient *ApiClient) getGuildById(guildId types.Snowflake, authorizationToken string) (*http.Response, error) {
	return apiClient.request("GET", fmt.Sprintf("guilds/%s", guildId), authorizationToken, 3)
}

func (apiClient *ApiClient) getGuildChannels(guildId types.Snowflake, authorizationToken string) (*http.Response, error) {
	return apiClient.request("GET", fmt.Sprintf("guilds/%s/channels", guildId), authorizationToken, 3)
}

func (apiClient *ApiClient) getChannelById(channelId types.Snowflake, authorizationToken string) (*http.Response, error) {
	return apiClient.request("GET", fmt.Sprintf("channels/%s", channelId), authorizationToken, 3)
}

type GetChannelMessagesOptions struct {
	Limit  *int
	Before *types.Snowflake
	After  *types.Snowflake
	Around *types.Snowflake
}

func (apiClient *ApiClient) getChannelMessages(channelId types.Snowflake, options *GetChannelMessagesOptions, authorizationToken string) (*http.Response, error) {
	searchParams := url.Values{}
	if options != nil {
		if options.Limit != nil {
			searchParams.Add("limit", strconv.Itoa(*options.Limit))
		}
		if options.Before != nil {
			searchParams.Add("before", string(*options.Before))
		}
		if options.After != nil {
			searchParams.Add("after", string(*options.After))
		}
		if options.Around != nil {
			searchParams.Add("around", string(*options.Around))
		}
	}
	return apiClient.request("GET", fmt.Sprintf("channels/%s/messages?%s", channelId, searchParams.Encode()), authorizationToken, 3)
}

type SearchChannelThreadsOptions struct {
	Offset    int
	Limit     *int
	Archived  *bool
	SortBy    *string
	SortOrder *string
}

func (apiClient *ApiClient) searchChannelThreads(authorizationToken string, mainChannelId types.Snowflake, options *SearchChannelThreadsOptions) (*http.Response, error) {
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
	return apiClient.request("GET", fmt.Sprintf("/channels/%s/threads/search?%s", mainChannelId, searchParams.Encode()), authorizationToken, 3)
}
