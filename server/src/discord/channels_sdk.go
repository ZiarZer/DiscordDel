package discord

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (sdk *DiscordSdk) GetChannel(channelId types.Snowflake, authorizationToken string) *types.Channel {
	resp, err := sdk.ApiClient.getChannelById(channelId, authorizationToken)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil
	}
	if resp.StatusCode != 200 {
		sdk.Log(string(body), utils.ERROR)
		return nil
	}
	var channel types.Channel
	json.Unmarshal(body, &channel)
	sdk.Log(fmt.Sprintf("Successfully got channel %s (#%s)", *channel.Name, channelId), utils.SUCCESS)
	sdk.Repo.InsertChannel(channel)
	return &channel
}

func (sdk *DiscordSdk) GetChannelMessages(authorizationToken string, channelId types.Snowflake, options *GetChannelMessagesOptions) []types.Message {
	resp, err := sdk.ApiClient.getChannelMessages(channelId, options, authorizationToken)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil
	}
	if resp.StatusCode != 200 {
		sdk.Log(string(body), utils.ERROR)
		return nil
	}
	var messages []types.Message
	json.Unmarshal(body, &messages)
	sdk.Log(fmt.Sprintf("Got %d messages in channel %s", len(messages), channelId), utils.SUCCESS)
	return messages
}

func (sdk *DiscordSdk) SearchChannelThreads(authorizationToken string, mainChannelId types.Snowflake, options *SearchChannelThreadsOptions) []types.Channel {
	resp, err := sdk.ApiClient.searchChannelThreads(authorizationToken, mainChannelId, options)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil
	}

	if resp.StatusCode != 200 {
		sdk.Log(string(body), utils.ERROR)
		return nil
	}
	var result types.ThreadsResult
	json.Unmarshal(body, &result)
	sdk.Log(fmt.Sprintf("Got %d threads in channel %s (offset %d)", len(result.Threads), mainChannelId, options.Offset), utils.SUCCESS)
	return result.Threads
}

func (sdk *DiscordSdk) GetThreadsData(authorizationToken string, mainChannelId types.Snowflake, threadIds []types.Snowflake) *types.ThreadsDataResult {
	resp, err := sdk.ApiClient.getThreadsData(authorizationToken, mainChannelId, threadIds)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return nil
	}

	if resp.StatusCode != 200 {
		sdk.Log(string(body), utils.ERROR)
		return nil
	}
	var result types.ThreadsDataResult
	json.Unmarshal(body, &result)
	sdk.Log(fmt.Sprintf("Got data for %d threads of channel %s", len(result.Threads), mainChannelId), utils.SUCCESS)
	return &result
}

func (sdk *DiscordSdk) UnarchiveThread(authorizationToken string, threadId types.Snowflake) bool {
	resp, err := sdk.ApiClient.modifyChannel(authorizationToken, threadId, ModifyChannelBody{Archived: utils.MakePointer(false)})
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return false
	}
	if resp.StatusCode != 200 {
		sdk.Log(fmt.Sprintf("Failed to unarchive thread %s", threadId), utils.ERROR)
		return false
	}
	sdk.Log(fmt.Sprintf("Successfully unarchived thread %s", threadId), utils.SUCCESS)
	return true
}
