package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (sdk *DiscordSdk) GetGuild(ctx context.Context, guildId types.Snowflake) (*types.Guild, error) {
	resp, err := sdk.ApiClient.getGuildById(ctx, guildId)
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
	var guild types.Guild
	json.Unmarshal(body, &guild)
	sdk.Log(fmt.Sprintf("Successfully got guild %s (#%s)", guild.Name, guildId), utils.SUCCESS)
	sdk.Repo.InsertGuild(guild)
	return &guild, nil
}

func (sdk *DiscordSdk) GetGuildChannels(ctx context.Context, guildId types.Snowflake) ([]types.Channel, error) {
	resp, err := sdk.ApiClient.getGuildChannels(ctx, guildId)
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
	}
	var channels []types.Channel
	json.Unmarshal(body, &channels)
	sdk.Log(fmt.Sprintf("Successfully got %d channels in guild %s", len(channels), guildId), utils.SUCCESS)
	sdk.Repo.InsertMultipleChannels(channels)
	return channels, nil
}

func (sdk *DiscordSdk) SearchGuildMessages(ctx context.Context, guildId types.Snowflake, options *SearchGuildMessagesOptions) ([]types.Message, error) {
	resp, err := sdk.ApiClient.searchGuildMessages(ctx, guildId, options)
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
	}
	var result types.GuildSearchResult
	json.Unmarshal(body, &result)
	messages := utils.Map(result.MessageArrays, func(arr []types.Message) types.Message { return arr[0] })
	if options.Offset > 0 {
		sdk.Log(fmt.Sprintf("Got %d messages (offset %d) in guild %s", len(result.MessageArrays), options.Offset, guildId), utils.SUCCESS)
	} else if options.MinId != nil {
		sdk.Log(fmt.Sprintf("Got %d messages (after %s) in guild %s", len(result.MessageArrays), *options.MinId, guildId), utils.SUCCESS)
	} else {
		sdk.Log(fmt.Sprintf("Got %d messages in guild %s", len(result.MessageArrays), guildId), utils.SUCCESS)
	}
	sdk.Repo.InsertMultipleMessages(messages, "PENDING")
	return messages, nil
}
