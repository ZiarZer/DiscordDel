package discord

import (
	"encoding/json"
	"fmt"
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
	sdk.Log(fmt.Sprintf("Successfully got channel %s (#%s)", channel.Name, channelId), utils.SUCCESS)
	return &channel
}
