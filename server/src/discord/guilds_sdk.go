package discord

import (
	"encoding/json"
	"io"

	"github.com/ZiarZer/DiscordDel/utils"
)

func GetGuild(guildId string, authorizationToken string) *Guild {
	resp, err := getGuildById(guildId, authorizationToken)
	if err != nil {
		utils.Log(err.Error(), utils.ERROR)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log(err.Error(), utils.ERROR)
		return nil
	}
	if resp.StatusCode != 200 {
		utils.Log(string(body), utils.ERROR)
		return nil
	}
	var guild Guild
	json.Unmarshal(body, &guild)
	return &guild
}

func GetGuildChannels(guildId string, authorizationToken string) []Channel {
	resp, err := getGuildChannels(guildId, authorizationToken)
	if err != nil {
		utils.Log(err.Error(), utils.ERROR)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log(err.Error(), utils.ERROR)
		return nil
	}
	if resp.StatusCode != 200 {
		utils.Log(string(body), utils.ERROR)
	}
	var channels []Channel
	json.Unmarshal(body, &channels)
	return channels
}
