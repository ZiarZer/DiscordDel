package discord

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (sdk *DiscordSdk) GetMessageReactions(authorizationToken string, channelId types.Snowflake, messageId types.Snowflake, emoji string, burst bool, options *GetMessageReactionsOptions) []types.User {
	resp, err := sdk.ApiClient.getMessageReactions(authorizationToken, channelId, messageId, emoji, burst, options)
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
	var usersReacted []types.User
	json.Unmarshal(body, &usersReacted)
	sdk.TempLog(fmt.Sprintf("Got %d users who reacted %s on message %s", len(usersReacted), emoji, messageId), utils.SUCCESS)
	return usersReacted
}

func (sdk *DiscordSdk) DeleteMessage(authorizationToken string, channelId types.Snowflake, messageId types.Snowflake) bool {
	resp, err := sdk.ApiClient.deleteMessage(authorizationToken, channelId, messageId)
	if err != nil {
		utils.InternalLog(err.Error(), utils.ERROR)
		return false
	}
	if resp.StatusCode != 204 {
		sdk.Log(fmt.Sprintf("Failed to delete message %s", messageId), utils.ERROR)
		return false
	}
	sdk.Log(fmt.Sprintf("Deleted message %s", messageId), utils.SUCCESS)
	return true
}
