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
