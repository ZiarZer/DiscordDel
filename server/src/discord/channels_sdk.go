package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (sdk *DiscordSdk) GetChannel(ctx context.Context, channelId types.Snowflake) (*types.Channel, error) {
	resp, err := sdk.ApiClient.getChannelById(ctx, channelId)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil, err
	}
	if resp.StatusCode != 200 {
		sdk.Log(string(body), utils.ERROR)
		return nil, nil
	}
	var channel types.Channel
	json.Unmarshal(body, &channel)
	sdk.Log(fmt.Sprintf("Successfully got channel %s (#%s)", *channel.Name, channelId), utils.SUCCESS)
	sdk.Repo.InsertChannel(channel)
	return &channel, nil
}

func (sdk *DiscordSdk) GetChannelMessages(ctx context.Context, channelId types.Snowflake, options *GetChannelMessagesOptions) ([]types.Message, error) {
	resp, err := sdk.ApiClient.getChannelMessages(ctx, channelId, options)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil, err
	}
	if resp.StatusCode != 200 {
		sdk.Log(string(body), utils.ERROR)
		return nil, nil
	}
	var messages []types.Message
	json.Unmarshal(body, &messages)
	sdk.Log(fmt.Sprintf("Got %d messages in channel %s", len(messages), channelId), utils.SUCCESS)
	return messages, nil
}

func (sdk *DiscordSdk) SearchChannelThreads(ctx context.Context, mainChannelId types.Snowflake, options *SearchChannelThreadsOptions) ([]types.Channel, error) {
	resp, err := sdk.ApiClient.searchChannelThreads(ctx, mainChannelId, options)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil, err
	}

	if resp.StatusCode != 200 {
		sdk.Log(string(body), utils.ERROR)
		return nil, nil
	}
	var result types.ThreadsResult
	json.Unmarshal(body, &result)
	sdk.Log(fmt.Sprintf("Got %d threads in channel %s (offset %d)", len(result.Threads), mainChannelId, options.Offset), utils.SUCCESS)
	return result.Threads, nil
}

func (sdk *DiscordSdk) GetThreadsData(ctx context.Context, mainChannelId types.Snowflake, threadIds []types.Snowflake) (*types.ThreadsDataResult, error) {
	resp, err := sdk.ApiClient.getThreadsData(ctx, mainChannelId, threadIds)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil, err
	}

	if resp.StatusCode != 200 {
		sdk.Log(string(body), utils.ERROR)
		return nil, nil
	}
	var result types.ThreadsDataResult
	json.Unmarshal(body, &result)
	sdk.Log(fmt.Sprintf("Got data for %d threads of channel %s", len(result.Threads), mainChannelId), utils.SUCCESS)
	return &result, nil
}

func (sdk *DiscordSdk) UnarchiveThread(ctx context.Context, threadId types.Snowflake) (bool, error) {
	resp, err := sdk.ApiClient.modifyChannel(ctx, threadId, ModifyChannelBody{Archived: utils.MakePointer(false)})
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return false, err
	}
	if resp.StatusCode != 200 {
		sdk.Log(fmt.Sprintf("Failed to unarchive thread %s", threadId), utils.ERROR)
		return false, nil
	}
	sdk.Log(fmt.Sprintf("Successfully unarchived thread %s", threadId), utils.SUCCESS)
	return true, nil
}
