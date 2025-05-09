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
	sdk.Log(fmt.Sprintf("Got %d users who reacted %s on message %s", len(usersReacted), emoji, messageId), utils.SUCCESS)
	return usersReacted
}
