package discord

import (
	"encoding/json"
	"io"

	"github.com/ZiarZer/DiscordDel/utils"
)

func (sdk *DiscordSdk) GetChannel(channelId string, authorizationToken string) *Channel {
	resp, err := getChannelById(channelId, authorizationToken)
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
	var channel Channel
	json.Unmarshal(body, &channel)
	return &channel
}
