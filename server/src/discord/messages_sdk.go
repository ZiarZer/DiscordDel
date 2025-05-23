package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (sdk *DiscordSdk) GetMessageReactions(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake, emoji string, burst bool, options *GetMessageReactionsOptions) ([]types.User, error) {
	resp, err := sdk.ApiClient.getMessageReactions(ctx, channelId, messageId, emoji, burst, options)
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
	var usersReacted []types.User
	json.Unmarshal(body, &usersReacted)
	sdk.TempLog(fmt.Sprintf("Got %d users who reacted %s on message %s", len(usersReacted), emoji, messageId), utils.SUCCESS)
	return usersReacted, nil
}

func (sdk *DiscordSdk) DeleteMessage(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake) (bool, error) {
	resp, err := sdk.ApiClient.deleteMessage(ctx, channelId, messageId)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return false, err
	}
	if resp.StatusCode != 204 {
		sdk.Log(fmt.Sprintf("Failed to delete message %s", messageId), utils.ERROR)
		return false, nil
	}
	sdk.Log(fmt.Sprintf("Deleted message %s", messageId), utils.SUCCESS)
	return true, nil
}

func (sdk *DiscordSdk) DeleteOwnReaction(ctx context.Context, channelId types.Snowflake, messageId types.Snowflake, emoji string) (bool, error) {
	resp, err := sdk.ApiClient.deleteOwnRection(ctx, channelId, messageId, emoji)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return false, err
	}
	if resp.StatusCode != 204 {
		sdk.Log(fmt.Sprintf("Failed to delete reaction %s on message %s", emoji, messageId), utils.ERROR)
		return false, nil
	}
	sdk.Log(fmt.Sprintf("Deleted reaction %s on message %s", emoji, messageId), utils.SUCCESS)
	return true, nil
}
