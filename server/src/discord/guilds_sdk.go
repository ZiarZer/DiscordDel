package discord

import (
	"encoding/json"
	"io"

	"github.com/ZiarZer/DiscordDel/utils"
)

func (sdk *DiscordSdk) GetGuild(guildId string, authorizationToken string) *Guild {
	resp, err := getGuildById(guildId, authorizationToken)
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
	var guild Guild
	json.Unmarshal(body, &guild)
	return &guild
}

func (sdk *DiscordSdk) GetGuildChannels(guildId string, authorizationToken string) []Channel {
	resp, err := getGuildChannels(guildId, authorizationToken)
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
	}
	var channels []Channel
	json.Unmarshal(body, &channels)
	return channels
}
