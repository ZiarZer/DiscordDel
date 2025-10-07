package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	Delay               int // Delay in ms
	lastRequestUnixTime int64
}

func (apiClient *ApiClient) request(ctx context.Context, method string, endpoint string, body *[]byte, retriesLeft int) (*http.Response, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		break
	}

	url := fmt.Sprintf("%s/%s", DiscordApibaseURL, endpoint)
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(*body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Authorization", ctx.Value(types.CtxKey{Key: "authorizationToken"}).(string))
	httpClient := &http.Client{}

	timeSinceLastRequest := time.Now().UnixMilli() - int64(apiClient.lastRequestUnixTime)
	if timeSinceLastRequest < int64(apiClient.Delay) {
		time.Sleep(time.Duration(apiClient.Delay*1000000) - time.Duration(timeSinceLastRequest*1000000))
	}
	resp, err := httpClient.Do(req)
	apiClient.lastRequestUnixTime = time.Now().UnixMilli()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 429 {
		utils.InternalLog("Rate limited", utils.WARNING)
		apiClient.Delay = int(math.Round(1.5 * float64(apiClient.Delay)))
		utils.InternalLog(fmt.Sprintf("Delay multiplied by 1.5, current value: %dms", apiClient.Delay), utils.INFO)

		retryAfter, err := strconv.ParseFloat(resp.Header.Get("retry_after"), 32)
		if err != nil {
			retryAfter = 2.
		}
		time.Sleep(time.Duration(retryAfter * 1000000000))
		if retriesLeft > 0 {
			return apiClient.request(ctx, method, endpoint, body, retriesLeft-1)
		}
	}
	return resp, nil
}

func (apiClient *ApiClient) login(ctx context.Context) (*http.Response, error) {
	return apiClient.request(ctx, "GET", "users/@me", nil, 3)
}

func (apiClient *ApiClient) getUserGuilds(ctx context.Context) (*http.Response, error) {
	return apiClient.request(ctx, "GET", "users/@me/guilds", nil, 3)
}

func (apiClient *ApiClient) getGuildById(ctx context.Context, guildId types.Snowflake) (*http.Response, error) {
	return apiClient.request(ctx, "GET", fmt.Sprintf("guilds/%s", guildId), nil, 3)
}

func (apiClient *ApiClient) getGuildChannels(ctx context.Context, guildId types.Snowflake) (*http.Response, error) {
	return apiClient.request(ctx, "GET", fmt.Sprintf("guilds/%s/channels", guildId), nil, 3)
}

type SearchGuildMessagesOptions struct {
	AuthorIds  []types.Snowflake
	ChannelIds []types.Snowflake
	Offset     int
	SortBy     *string
	SortOrder  *string
	MinId      *types.Snowflake
	MaxId      *types.Snowflake
}

func (apiClient *ApiClient) searchGuildMessages(ctx context.Context, guildId types.Snowflake, options *SearchGuildMessagesOptions) (*http.Response, error) {
	searchParams := url.Values{}
	if options != nil {
		if options.Offset > 9975 {
			return nil, errors.New("offset must be less than or equal to 9975")
		}
		searchParams.Add("offset", strconv.Itoa(options.Offset))

		for _, authorId := range options.AuthorIds {
			searchParams.Add("author_id", string(authorId))
		}
		for _, channelId := range options.ChannelIds {
			searchParams.Add("channel_id", string(channelId))
		}
		if options.SortOrder != nil {
			searchParams.Add("sort_order", *options.SortOrder)
		}
		if options.SortBy != nil {
			searchParams.Add("sort_by", *options.SortBy)
		}
		if options.MinId != nil {
			searchParams.Add("min_id", string(*options.MinId))
		}
		if options.MaxId != nil {
			searchParams.Add("max_Id", string(*options.MaxId))
		}
	}
	return apiClient.request(ctx, "GET", fmt.Sprintf("guilds/%s/messages/search?%s", guildId, searchParams.Encode()), nil, 3)
}

func (apiClient *ApiClient) getChannelById(ctx context.Context, channelId types.Snowflake) (*http.Response, error) {
	return apiClient.request(ctx, "GET", fmt.Sprintf("channels/%s", channelId), nil, 3)
}

type GetChannelMessagesOptions struct {
	Limit  *int
	Before *types.Snowflake
	After  *types.Snowflake
	Around *types.Snowflake
}

func (apiClient *ApiClient) getChannelMessages(ctx context.Context, channelId types.Snowflake, options *GetChannelMessagesOptions) (*http.Response, error) {
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
	return apiClient.request(ctx, "GET", fmt.Sprintf("channels/%s/messages?%s", channelId, searchParams.Encode()), nil, 3)
}

type SearchChannelThreadsOptions struct {
	Offset    int
	Limit     *int
	Archived  *bool
	SortBy    *string
	SortOrder *string
}

func (apiClient *ApiClient) searchChannelThreads(ctx context.Context, mainChannelId types.Snowflake, options *SearchChannelThreadsOptions) (*http.Response, error) {
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
	return apiClient.request(ctx, "GET", fmt.Sprintf("channels/%s/threads/search?%s", mainChannelId, searchParams.Encode()), nil, 3)
}

type GetThreadDataBody struct {
	ThreadIds []types.Snowflake `json:"thread_ids"`
}

func (apiClient *ApiClient) getThreadsData(ctx context.Context, mainChannelId types.Snowflake, threadIds []types.Snowflake) (*http.Response, error) {
	body := GetThreadDataBody{ThreadIds: threadIds}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		utils.InternalLog("Failed to serialize thread data", utils.ERROR)
		return nil, err
	}
	return apiClient.request(ctx, "POST", fmt.Sprintf("channels/%s/post-data", mainChannelId), &jsonBody, 3)
}

type ModifyChannelBody struct {
	Archived *bool `json:"archived"`
}

func (apiClient *ApiClient) modifyChannel(ctx context.Context, channelId types.Snowflake, body ModifyChannelBody) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		utils.InternalLog("Failed to serialize update channel data", utils.ERROR)
		return nil, err
	}
	return apiClient.request(ctx, "PATCH", fmt.Sprintf("channels/%s", channelId), &jsonBody, 3)
}

type GetMessageReactionsOptions struct {
	Limit *int
	After *types.Snowflake
}

func (apiClient *ApiClient) deleteMessage(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake) (*http.Response, error) {
	return apiClient.request(ctx, "DELETE", fmt.Sprintf("channels/%s/messages/%s", channelId, messageId), nil, 3)
}

func (apiClient *ApiClient) getMessageReactions(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake, emoji string, burst bool, options *GetMessageReactionsOptions) (*http.Response, error) {
	searchParams := url.Values{}
	if burst {
		searchParams.Add("type", "1")
	} else {
		searchParams.Add("type", "0")
	}
	if options != nil {
		if options.Limit != nil {
			searchParams.Add("limit", strconv.Itoa(*options.Limit))
		}
		if options.After != nil {
			searchParams.Add("after", string(*options.After))
		}
	}
	return apiClient.request(ctx, "GET", fmt.Sprintf("channels/%s/messages/%s/reactions/%s?%s", channelId, messageId, emoji, searchParams.Encode()), nil, 3)
}

func (apiClient *ApiClient) deleteOwnRection(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake, emoji string) (*http.Response, error) {
	return apiClient.request(ctx, "DELETE", fmt.Sprintf("channels/%s/messages/%s/reactions/%s/@me", channelId, messageId, emoji), nil, 3)
}
