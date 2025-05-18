package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (sdk *DiscordSdk) Login(ctx context.Context) *types.User {
	resp, err := sdk.ApiClient.login(ctx)
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
	var loggedUser types.User
	json.Unmarshal(body, &loggedUser)
	sdk.Log(fmt.Sprintf("Successfully authenticated as %s (%s)", loggedUser.Username, loggedUser.Id), utils.SUCCESS)
	return &loggedUser
}

func (sdk *DiscordSdk) GetUserGuilds(ctx context.Context) []types.Guild {
	resp, err := sdk.ApiClient.getUserGuilds(ctx)
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
	var guilds []types.Guild
	json.Unmarshal(body, &guilds)
	sdk.Log(fmt.Sprintf("Successfully got %d guilds for current user", len(guilds)), utils.SUCCESS)
	sdk.Repo.InsertMultipleGuilds(guilds)
	return guilds
}
